package mem

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrNoProcess       = errors.New("no process matching the criteria was found")
	ErrPatternNotFound = errors.New("no memory matched the pattern")
)

type (
	Process interface {
		io.Closer
		io.ReaderAt
		Pid() int
		Maps() ([]Map, error)
	}

	Map interface {
		Start() uint64
		Size() uint64
	}
)

type pattern struct {
	Bytes []uint32
	Mask  []uint32
}

func parsePattern(s string) (pattern, error) {
	var bytes, mask []byte
	for _, bytestr := range strings.Split(s, " ") {
		if bytestr == "??" {
			bytes = append(bytes, 0x00)
			mask = append(mask, 0x00)
			continue
		}
		b, err := strconv.ParseUint(bytestr, 16, 8)
		if err != nil {
			return pattern{}, err
		}
		bytes = append(bytes, byte(b))
		mask = append(mask, 0xFF)
	}

	var p pattern
	for i := 0; i < len(bytes); i += 4 {
		var byt uint32
		for i, b := range bytes[i : i+4] {
			byt |= uint32(b) << (i * 8)
		}
		p.Bytes = append(p.Bytes, byt)
		var mas uint32
		for i, m := range mask[i : i+4] {
			mas |= uint32(m) << (i * 8)
		}
		p.Mask = append(p.Mask, mas)
	}

	return p, nil
}

func search(buf []byte, needle uint32, mask uint32) (int, bool) {
	for i := 0; i+4 < len(buf); i++ {
		haystack := binary.LittleEndian.Uint32(buf[i : i+4])
		if needle^haystack&mask == 0 {
			return i, true
		}
	}
	return 0, false
}

func find(p Process, pat pattern, reg Map) (uint64, bool) {
	const bufsize = 32768
outer:
	for i := uint64(0); i < reg.Size(); i += bufsize - 3 {
		var buf [bufsize]byte
		n, err := p.ReadAt(buf[:], int64(reg.Start()+i))
		if err != nil {
			continue
		}
		at, ok := search(buf[:n], pat.Bytes[0], pat.Mask[0])
		if ok && at+(4*len(pat.Bytes)) < len(buf) {
			for j := range pat.Bytes {
				needle := pat.Bytes[j]
				mask := pat.Mask[j]
				haystack := binary.LittleEndian.Uint32(buf[at+(4*j) : at+4+(4*j)])
				if needle^haystack&mask != 0 {
					continue outer
				}
			}
			return uint64(reg.Start() + i + uint64(at)), true
		}
	}
	return 0, false
}

func Scan(p Process, pattern string) (uint64, error) {
	maps, err := p.Maps()
	if err != nil {
		return 0, err
	}

	pat, err := parsePattern(pattern)
	if err != nil {
		return 0, err
	}

	for _, reg := range maps {
		i, ok := find(p, pat, reg)
		if ok {
			return i, nil
		}
	}

	return 0, ErrPatternNotFound
}

func ResolvePatterns(p Process, offsets interface{}) error {
	pval := reflect.ValueOf(offsets)
	val := reflect.Indirect(pval)
	valt := val.Type()
	if pval.Kind() != reflect.Ptr || val.Kind() != reflect.Struct {
		panic("offsets must be a pointer to a struct")
	}
	var anyErr error
	for i := 0; i < val.NumField(); i++ {
		field := valt.Field(i)
		sig, ok := field.Tag.Lookup("sig")
		if !ok {
			continue
		}
		offset, err := Scan(p, sig)
		if err != nil {
			anyErr = err
			continue
		}
		val.Field(i).Set(reflect.ValueOf(offset))
	}

	return anyErr
}
