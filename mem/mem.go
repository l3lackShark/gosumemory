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
		Start() int64
		Size() int64
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

func Scan(p Process, pattern string) (uint64, error) {
	maps, err := p.Maps()
	if err != nil {
		return 0, err
	}

	var largestMap uint64
	for _, reg := range maps {
		if size := reg.Size(); size > largestMap && size < 1073741824 {
			largestMap = size
		}
	}
	if largestMap == 0 {
		return 0, ErrPatternNotFound
	}

	pat, err := parsePattern(pattern)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, int(largestMap))
	for _, reg := range maps {
		addr := reg.Start()
		size := reg.Size()
		if size >= 1073741824 {
			continue
		}
		n, err := p.ReadAt(buf[0:size], int64(addr))
		if err != nil || n != int(size) {
			continue
		}
		needle := pat.Bytes[0]
		mask := pat.Mask[0]

		var i uint64
	outer:
		for i = 0; i+4 < size; i++ {
			haystack := binary.LittleEndian.Uint32(buf[i : i+4])
			if needle^haystack&mask == 0 {
				for j := range pat.Bytes {
					needle := pat.Bytes[j]
					mask := pat.Mask[j]
					haystack := binary.LittleEndian.Uint32(
						buf[i+uint64(4*j) : i+4+uint64(4*j)])
					if needle^haystack&mask != 0 {
						continue outer
					}
				}
				return reg.Start() + i, nil
			}
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
