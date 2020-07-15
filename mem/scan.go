package mem

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"
	"strings"
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

func maskbyte(needle, mask uint32) (search byte, offset int) {
	for offset := 0; offset < 4; offset++ {
		maskByte := byte(mask >> (offset * 8) & 0xFF)
		needleByte := byte(needle >> (offset * 8) & 0xFF)
		if maskByte != 0x00 && needleByte != 0x00 {
			return needleByte, offset
		}
	}
	panic("empty mask (bad pattern)")
}

func search(buf []byte, pat pattern) (int, bool) {
	progress := 0
	maskbyte, byteoffset := maskbyte(pat.Bytes[0], pat.Mask[0])
	for {
		i := bytes.IndexByte(buf, maskbyte)
		if i == -1 {
			return 0, false
		}
		begin := i - byteoffset
		end := begin + len(pat.Bytes)*4
		if begin < 0 || end > len(buf) {
			progress += i + 1
			buf = buf[i+1:]
			continue
		}
		success := true
		for j := range pat.Bytes {
			needle, mask := pat.Bytes[j], pat.Mask[j]
			slice := buf[begin+(j*4) : begin+(j*4)+4]
			haystack := binary.LittleEndian.Uint32(slice)
			if needle^haystack&mask != 0 {
				success = false
				break
			}
		}
		if success {
			return begin + progress, true
		}
		progress += i + 1
		buf = buf[i+1:]
	}
}

func find(p Process, pat pattern, reg Map) (int64, error) {
	const bufsize = 65536
	var buf [bufsize]byte
	for i := int64(0); i < reg.Size(); {
		ntoread := reg.Size() - i
		if ntoread >= bufsize {
			ntoread = bufsize - 1
		}
		n, err := p.ReadAt(buf[:ntoread], int64(reg.Start()+i))
		if err != nil {
			return 0, err
		}
		if at, ok := search(buf[:n], pat); ok {
			return int64(reg.Start() + i + int64(at)), nil
		}
		diff := n - len(pat.Bytes)*8
		if diff <= 1 {
			diff = 1
		}
		i += int64(diff)
	}
	return 0, ErrPatternNotFound
}

func Scan(p Process, pattern string) (int64, error) {
	maps, err := p.Maps()
	if err != nil {
		return 0, err
	}

	pat, err := parsePattern(pattern)
	if err != nil {
		return 0, err
	}

	for _, reg := range maps {
		if i, err := find(p, pat, reg); err == nil {
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
