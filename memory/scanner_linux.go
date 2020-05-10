// +build linux

package memory

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type mmap struct {
	Start    uintptr
	End      uintptr
	Perms    string
	Offset   uintptr
	DevMajor int
	DevMinor int
	Inode    int
	Path     string
}

func (m *mmap) Size() uintptr {
	return m.End - m.Start
}

func scan(f *os.File, maps []mmap, pattern string) (uint32, error) {
	var largestMap uintptr
	for _, amap := range maps {
		if amap.Perms[0] != 'r' || amap.Perms[2] != 'x' {
			continue
		}
		if amap.Size() > largestMap {
			largestMap = amap.Size()
		}
	}
	pat, err := parsePattern(pattern)
	if err != nil {
		return 0, err
	}
	buf := make([]byte, largestMap)
	for _, amap := range maps {
		if amap.Perms[0] != 'r' || amap.Perms[2] != 'x' {
			continue
		}
		size := amap.Size()
		_, err := f.Seek(int64(amap.Start), 0)
		if err != nil {
			return 0, err
		}
		_, err = io.ReadFull(f, buf[0:size])
		if err != nil {
			continue
		}
		needle := pat.Bytes[0]
		mask := pat.Mask[0]
		var j uintptr
	outer:
		for j = 0; (j + 4) < size; j++ {
			haystack := binary.LittleEndian.Uint32(buf[j : j+4])
			if needle^haystack&mask == 0 {
				for k := range pat.Bytes {
					needle := pat.Bytes[k]
					mask := pat.Mask[k]
					haystack := binary.LittleEndian.Uint32(
						buf[j+uintptr(4*k) : j+4+uintptr(4*k)])
					if needle^haystack&mask != 0 {
						continue outer
					}
				}
				return uint32(amap.Start) + uint32(j), nil
			}
		}
	}
	return 0, errPatternNotFound
}

func readMaps(pid int) ([]mmap, error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/maps", pid))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var maps []mmap
	s := bufio.NewScanner(f)
	for s.Scan() {
		var amap mmap
		_, err := fmt.Sscanf(
			s.Text(), "%x-%x %s %x %x:%x %d %s",
			&amap.Start, &amap.End, &amap.Perms, &amap.Offset,
			&amap.DevMajor, &amap.DevMinor, &amap.Inode, &amap.Path,
		)
		if err != nil && err != io.EOF {
			continue
		}
		maps = append(maps, amap)
	}
	return maps, nil
}

type pattern struct {
	Bytes []uint32
	Mask  []uint32
}

func parsePattern(s string) (*pattern, error) {
	var bytes, mask []byte
	for _, bytestr := range strings.Split(s, " ") {
		if bytestr == "??" {
			bytes = append(bytes, 0x00)
			mask = append(mask, 0x00)
			continue
		}
		b, err := strconv.ParseUint(bytestr, 16, 8)
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, byte(b))
		mask = append(mask, 0xFF)
	}
	var p pattern
	for i := 0; i < len(bytes); i += 4 {
		var byt, mas uint32
		for i, b := range bytes[i : i+4] {
			byt |= uint32(b) << (i * 8)
		}
		for i, m := range mask[i : i+4] {
			mas |= uint32(m) << (i * 8)
		}
		p.Bytes = append(p.Bytes, byt)
		p.Mask = append(p.Mask, mas)
	}
	return &p, nil
}

var (
	errNoPIDMatched = errors.New("No PID matched the criteria")
)

var (
	errPatternNotFound = errors.New("Pattern not found")
)
