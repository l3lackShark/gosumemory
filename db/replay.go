package db

import (
	"bufio"
	"encoding/binary"
	"os"
	"reflect"

	"github.com/k0kubun/pp"
)

type osr struct {
	Gamemode   uint8
	OsuVer     int32
	MD5        string
	PlayerName string
	BmChecksum string
	Hit300s    uint16
	Hit100s    uint16
	Hit50s     uint16
	HitGekis   uint16
	HitKatus   uint16
	HitMisses  uint16
	Score      int32
	MaxCombo   uint16
	IsPerfect  bool //bool
	Mods       int32
	Lifebar    string
	DateTime   int64
	ReplayData int32
	ScoreID    int64
}

//WriteOSR does the write replay magic
func WriteOSR() error {
	var OsrStruct = osr{
		Gamemode:   0,
		OsuVer:     20200715,
		MD5:        "a185ae7fa76162b434a973bfd5426e0d",
		PlayerName: "BlackShark",
		BmChecksum: "c02eb46a8b577e95ee5d43c7c74b9d2e",
		Hit300s:    1041,
		Hit100s:    163,
		Hit50s:     22,
		HitGekis:   140,
		HitKatus:   77,
		HitMisses:  16,
		Score:      6395654,
		MaxCombo:   454,
		IsPerfect:  false,
		Mods:       0,
		Lifebar:    "",
		DateTime:   636231068270000000,
		ReplayData: 0,
		ScoreID:    0,
	}

	file, err := os.Create("test.osr")
	defer file.Close()
	if err != nil {
		pp.Println("Could not create osr file")
		return nil
	}
	replayWriter := bufio.NewWriter(file)

	v := reflect.ValueOf(OsrStruct)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		switch v.Field(i).Kind() {
		case reflect.String:
			replayWriter.WriteByte(0x0B)
			replayWriter.WriteByte(byte(len(v.Field(i).String())))
			replayWriter.WriteString(v.Field(i).String())
		case reflect.Uint8:
			writeUint8(replayWriter, uint8(v.Field(i).Uint()))
		case reflect.Uint16:
			writeUint16(replayWriter, uint16(v.Field(i).Uint()))
		case reflect.Int32:
			writeInt32(replayWriter, int32(v.Field(i).Int()))
		case reflect.Bool:
			writeBool(replayWriter, v.Field(i).Bool())
		case reflect.Int64:
			writeInt64(replayWriter, v.Field(i).Int())
		}

	}
	replayWriter.Flush()
	return nil
}

func writeUint8(replayFile *bufio.Writer, number uint8) {
	binary.Write(replayFile, binary.LittleEndian, number)
}
func writeInt32(replayFile *bufio.Writer, number int32) {
	binary.Write(replayFile, binary.LittleEndian, number)
}
func writeUint16(replayFile *bufio.Writer, number uint16) {
	binary.Write(replayFile, binary.LittleEndian, number)
}
func writeInt64(replayFile *bufio.Writer, number int64) {
	binary.Write(replayFile, binary.LittleEndian, number)
}
func writeBool(replayFile *bufio.Writer, number bool) {
	binary.Write(replayFile, binary.LittleEndian, number)
}
