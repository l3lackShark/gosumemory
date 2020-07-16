package mem

import (
	"encoding/binary"
	"io"
	"math"
	"unicode/utf16"
)

func readFullAt(r io.ReaderAt, buf []byte, off int64) (n int, err error) {
	for n < len(buf) && err == nil {
		var nn int
		nn, err = r.ReadAt(buf[n:], off+int64(n))
		n += nn
	}
	if n >= len(buf) {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

func bytesToInt(buf []byte) uint64 {
	var num uint64
	for i := range buf {
		num |= uint64(buf[i]) << (8 * i)
	}
	return num
}

func readUintRaw(r io.ReaderAt, addr int64, size int) (uint64, error) {
	var buf [8]byte
	if _, err := readFullAt(r, buf[:size], addr); err != nil {
		return 0, err
	}
	return bytesToInt(buf[:size]), nil
}

func readUint(r io.ReaderAt, addr int64, size int, offsets ...int64) (uint64, error) {
	last := len(offsets) - 1
	for _, offset := range offsets[:last] {
		newaddr, err := readUintRaw(r, addr+offset, 4)
		if err != nil {
			return 0, err
		}
		addr = int64(newaddr)
	}
	return readUintRaw(r, addr+offsets[last], size)
}

func readUintArray(r io.ReaderAt, addr int64, size int,
	offsets ...int64) ([]uint64, error) {
	base, err := ReadInt32(r, addr, offsets...)
	if err != nil {
		return nil, err
	}

	length, err := ReadInt32(r, int64(base), 12)
	if err != nil {
		return nil, err
	}

	data, err := ReadInt32(r, int64(base), 4)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, int(length)*size)
	_, err = readFullAt(r, buf, int64(data))
	if err != nil {
		return nil, err
	}

	buf64 := make([]uint64, length)
	for i := 0; i < int(length); i++ {
		buf64[i] = bytesToInt(buf[i*size : i*size+size])
	}

	return buf64, nil
}

func ReadString(r io.ReaderAt, addr int64, offsets ...int64) (string, error) {
	base, err := ReadInt32(r, addr, offsets...)
	if err != nil {
		return "", err
	}

	length, err := ReadInt32(r, int64(base), 4)
	if err != nil {
		return "", err
	}

	buf := make([]byte, length*2)

	_, err = readFullAt(r, buf, int64(base)+8)
	if err != nil {
		return "", err
	}

	buf16 := make([]uint16, length)
	for i := int32(0); i < length; i += 1 {
		buf16[i] = binary.LittleEndian.Uint16(buf[i*2 : i*2+2])
	}

	return string(utf16.Decode(buf16)), nil
}

func ReadInt8(r io.ReaderAt, addr int64, offsets ...int64) (int8, error) {
	num, err := readUint(r, addr, 1, offsets...)
	return int8(int64(num)), err
}

func ReadInt8Array(r io.ReaderAt, addr int64, offsets ...int64) ([]int8, error) {
	array, err := readUintArray(r, addr, 1, offsets...)
	array64 := make([]int8, len(array))
	for i, v := range array {
		array64[i] = int8(int64(v))
	}
	return array64, err
}

func ReadUint8(r io.ReaderAt, addr int64, offsets ...int64) (uint8, error) {
	num, err := readUint(r, addr, 1, offsets...)
	return uint8(num), err
}

func ReadUint8Array(r io.ReaderAt, addr int64, offsets ...int64) ([]uint8, error) {
	array, err := readUintArray(r, addr, 1, offsets...)
	array8 := make([]uint8, len(array))
	for i, v := range array {
		array8[i] = uint8(v)
	}
	return array8, err
}

func ReadInt16(r io.ReaderAt, addr int64, offsets ...int64) (int16, error) {
	num, err := readUint(r, addr, 2, offsets...)
	return int16(int64(num)), err
}

func ReadInt16Array(r io.ReaderAt, addr int64, offsets ...int64) ([]int16, error) {
	array, err := readUintArray(r, addr, 2, offsets...)
	array16 := make([]int16, len(array))
	for i, v := range array {
		array16[i] = int16(int64(v))
	}
	return array16, err
}

func ReadUint16(r io.ReaderAt, addr int64, offsets ...int64) (uint16, error) {
	num, err := readUint(r, addr, 2, offsets...)
	return uint16(num), err
}

func ReadUint16Array(r io.ReaderAt, addr int64, offsets ...int64) ([]uint16, error) {
	array, err := readUintArray(r, addr, 2, offsets...)
	array16 := make([]uint16, len(array))
	for i, v := range array {
		array16[i] = uint16(v)
	}
	return array16, err
}

func ReadInt32(r io.ReaderAt, addr int64, offsets ...int64) (int32, error) {
	num, err := readUint(r, addr, 4, offsets...)
	return int32(int64(num)), err
}

func ReadInt32Array(r io.ReaderAt, addr int64, offsets ...int64) ([]int32, error) {
	array, err := readUintArray(r, addr, 4, offsets...)
	array32 := make([]int32, len(array))
	for i, v := range array {
		array32[i] = int32(int64(v))
	}
	return array32, err
}

func ReadUint32(r io.ReaderAt, addr int64, offsets ...int64) (uint32, error) {
	num, err := readUint(r, addr, 4, offsets...)
	return uint32(num), err
}

func ReadUint32Array(r io.ReaderAt, addr int64, offsets ...int64) ([]uint32, error) {
	array, err := readUintArray(r, addr, 4, offsets...)
	array32 := make([]uint32, len(array))
	for i, v := range array {
		array32[i] = uint32(v)
	}
	return array32, err
}

func ReadInt64(r io.ReaderAt, addr int64, offsets ...int64) (int64, error) {
	num, err := readUint(r, addr, 8, offsets...)
	return int64(num), err
}

func ReadInt64Array(r io.ReaderAt, addr int64, offsets ...int64) ([]int64, error) {
	array, err := readUintArray(r, addr, 8, offsets...)
	array64 := make([]int64, len(array))
	for i, v := range array {
		array64[i] = int64(v)
	}
	return array64, err
}

func ReadUint64(r io.ReaderAt, addr int64, offsets ...int64) (uint64, error) {
	num, err := readUint(r, addr, 8, offsets...)
	return uint64(num), err
}

func ReadUint64Array(r io.ReaderAt, addr int64, offsets ...int64) ([]uint64, error) {
	array, err := readUintArray(r, addr, 8, offsets...)
	array64 := make([]uint64, len(array))
	for i, v := range array {
		array64[i] = v
	}
	return array64, err
}

func ReadFloat32(r io.ReaderAt, addr int64, offsets ...int64) (float32, error) {
	num, err := readUint(r, addr, 4, offsets...)
	return math.Float32frombits(uint32(num)), err
}

func ReadFloat32Array(r io.ReaderAt, addr int64, offsets ...int64) ([]float32, error) {
	array, err := readUintArray(r, addr, 4, offsets...)
	array32 := make([]float32, len(array))
	for i, v := range array {
		array32[i] = math.Float32frombits(uint32(v))
	}
	return array32, err
}

func ReadFloat64(r io.ReaderAt, addr int64, offsets ...int64) (float64, error) {
	num, err := readUint(r, addr, 8, offsets...)
	return math.Float64frombits(num), err
}

func ReadFloat64Array(r io.ReaderAt, addr int64, offsets ...int64) ([]float64, error) {
	array, err := readUintArray(r, addr, 8, offsets...)
	array64 := make([]float64, len(array))
	for i, v := range array {
		array64[i] = math.Float64frombits(v)
	}
	return array64, err
}
