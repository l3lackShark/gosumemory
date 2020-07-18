package mem

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/pkg/errors"
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

func Read(r io.ReaderAt, addresses interface{}, p interface{}) error {
	addrpval := reflect.ValueOf(addresses)
	addrval := addrpval.Elem()

	if addrval.Kind() != reflect.Struct {
		panic("addresses must be a pointer to a struct")
	}

	pval := reflect.ValueOf(p)
	val := pval.Elem()
	valt := val.Type()

	if val.Kind() != reflect.Struct {
		panic("p must be a pointer to a struct")
	}

	var anyerr error
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldt := valt.Field(i)
		tag, ok := fieldt.Tag.Lookup("mem")
		if !ok {
			continue
		}
		evalFunc := func(addr int64) (int64, error) {
			v, err := ReadInt32(r, addr, 0x0)
			return int64(v), err
		}
		var varFunc func(name string) (int64, error)
		varFunc = func(name string) (int64, error) {
			field := addrval.FieldByName(name)
			if field.IsValid() {
				addr := field.Interface().(int64)
				return addr, nil
			}
			method := addrval.Addr().MethodByName(name)
			if method.IsValid() {
				ret := method.Call([]reflect.Value{})
				expr, err := parseMem(
					ret[0].Interface().(string), varFunc)
				if err != nil {
					return 0, err
				}
				return expr.eval(evalFunc)
			}
			return 0, fmt.Errorf("undefined variable %s", name)
		}
		expr, err := parseMem(tag, varFunc)
		if err != nil {
			return errors.Wrapf(err,
				"failed to parse mem tag for %s.%s",
				valt.Name(), fieldt.Name)
		}
		addr, err := expr.eval(evalFunc)
		if err != nil {
			return errors.Wrapf(err, "failed to read %s.%s",
				valt.Name(), fieldt.Name)
		}
		anyerr = readPrimitive(r, field.Addr().Interface(), addr, 0)
	}

	return anyerr
}

func readPrimitive(r io.ReaderAt, p interface{},
	addr int64, offsets ...int64) error {
	var err error
	switch p := p.(type) {
	case *int8:
		*p, err = ReadInt8(r, addr, offsets...)
	case *int16:
		*p, err = ReadInt16(r, addr, offsets...)
	case *int32:
		*p, err = ReadInt32(r, addr, offsets...)
	case *int64:
		*p, err = ReadInt64(r, addr, offsets...)
	case *uint8:
		*p, err = ReadUint8(r, addr, offsets...)
	case *uint16:
		*p, err = ReadUint16(r, addr, offsets...)
	case *uint32:
		*p, err = ReadUint32(r, addr, offsets...)
	case *uint64:
		*p, err = ReadUint64(r, addr, offsets...)
	case *float32:
		*p, err = ReadFloat32(r, addr, offsets...)
	case *float64:
		*p, err = ReadFloat64(r, addr, offsets...)
	case *[]int8:
		*p, err = ReadInt8Array(r, addr, offsets...)
	case *[]int16:
		*p, err = ReadInt16Array(r, addr, offsets...)
	case *[]int32:
		*p, err = ReadInt32Array(r, addr, offsets...)
	case *[]int64:
		*p, err = ReadInt64Array(r, addr, offsets...)
	case *[]uint8:
		*p, err = ReadUint8Array(r, addr, offsets...)
	case *[]uint16:
		*p, err = ReadUint16Array(r, addr, offsets...)
	case *[]uint32:
		*p, err = ReadUint32Array(r, addr, offsets...)
	case *[]uint64:
		*p, err = ReadUint64Array(r, addr, offsets...)
	case *[]float32:
		*p, err = ReadFloat32Array(r, addr, offsets...)
	case *[]float64:
		*p, err = ReadFloat64Array(r, addr, offsets...)
	case *string:
		*p, err = ReadString(r, addr, offsets...)
	default:
		err = fmt.Errorf("unknown type %T", p)
	}
	return err

}

type mem struct {
	Child  *mem
	Offset int64
}

func (m *mem) String() string {
	var b strings.Builder
	if m.Child != nil {
		fmt.Fprintf(&b, "[%s]", m.Child)
	}
	if m.Child != nil && m.Offset != 0 {
		b.WriteString(" + ")
	}
	if m.Offset != 0 {
		fmt.Fprintf(&b, "0x%x", m.Offset)
	}
	return b.String()
}

func (m *mem) eval(f func(p int64) (int64, error)) (int64, error) {
	if m.Child != nil {
		childAddr, err := m.Child.eval(f)
		if err != nil {
			return 0, err
		}

		dereferenced, err := f(childAddr)
		if err != nil {
			return 0, err
		}

		return dereferenced + m.Offset, nil
	} else {
		return m.Offset, nil
	}
}

func parseMem(tag string,
	varFunc func(name string) (int64, error)) (*mem, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(tag))
	s.Mode = scanner.ScanIdents | scanner.ScanInts
	return parseMemExpr(&s, varFunc, false)
}

func parseMemExpr(s *scanner.Scanner,
	varFunc func(name string) (int64, error), inBrackets bool) (*mem, error) {
	expr := &mem{}
	switch tok := s.Scan(); tok {
	case '[':
		inner, err := parseMemExpr(s, varFunc, true)
		if err != nil {
			return nil, err
		}
		expr.Child = inner
	case scanner.Ident:
		name := s.TokenText()
		var err error
		expr.Offset, err = varFunc(name)
		if err != nil {
			return nil, err
		}
	case scanner.Int:
		var err error
		expr.Offset, err = strconv.ParseInt(s.TokenText(), 0, 64)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected token %d (%s)",
			tok, s.TokenText())
	}

	switch tok := s.Scan(); tok {
	case '+', '-':
		rest, err := parseMemExpr(s, varFunc, inBrackets)
		if err != nil {
			return nil, err
		}
		switch tok {
		case '+':
			expr.Offset += rest.Offset
		case '-':
			expr.Offset -= rest.Offset
		}
		return expr, nil
	case scanner.EOF, ']':
		if tok == ']' && !inBrackets {
			return nil, fmt.Errorf("unexpected token %d (%s)",
				tok, s.TokenText())
		}
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token %d (%s)",
			tok, s.TokenText())
	}
}
