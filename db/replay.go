package db

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/memory"

	"github.com/ulikunitz/xz/lzma"
)

type osr struct {
	Gamemode         uint8
	OsuVer           int32
	MD5              string
	PlayerName       string
	BmChecksum       string
	Hit300s          uint16
	Hit100s          uint16
	Hit50s           uint16
	HitGekis         uint16
	HitKatus         uint16
	HitMisses        uint16
	Score            int32
	MaxCombo         uint16
	IsPerfect        bool
	Mods             int32
	Lifebar          string
	DateTime         int64
	LengthReplayData int32
	ReplayData       []uint8
	ScoreID          int64
}

var tempBeatmapFailTime int32

const ticksUnix = 621355968000000000 //C# DateTime

//WriteOSR does the write replay magic
func InitOSRWriting() {
	if _, err := os.Stat("FailedReplays"); os.IsNotExist(err) {
		fmt.Println("FailedReplays Directory does not exist. Making one..")
		err := os.Mkdir("FailedReplays", 0644)
		if err != nil {
			pp.Println("Wasn't able to create the 'FailedReplays' direcory. Failed replays functionality will be unavailable'")
			return
		}
	}
	writeOSR()
	return
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

func compressToLZMA(input string) []byte {
	text := input
	var buf bytes.Buffer
	// compress text
	w, err := lzma.WriterConfig{DictCap: 16 * 1024 * 1024}.NewWriter(&buf)
	if err != nil {
		log.Fatalf("xz.NewWriter error %s", err)
	}
	if _, err := io.WriteString(w, text); err != nil {
		log.Fatalf("WriteString error %s", err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("w.Close error %s", err)
	}
	return buf.Bytes()
}

func convertMemoryDataToOSRStruct() osr {
	osrStruct := memory.GameplayData.Replay
	var lzma []string
	for i, replayTick := range osrStruct.Replays {
		if i > 0 {
			replayTick.Time = replayTick.Time - osrStruct.Replays[i-1].Time
		}

		lzma = append(lzma, fmt.Sprintf("%d|%f|%f|%d", replayTick.Time, replayTick.X, replayTick.Y, replayTick.WasButtonPressed)) //0|256|-500|0, f.e.
	}
	lzma = append(lzma, "-12345|0|0|0,") //every replay has this at the end
	decompressedLZMAStr := strings.Join(lzma, ",")
	compressed := compressToLZMA(decompressedLZMAStr)

	var replay = osr{
		Gamemode:         uint8(memory.GameplayData.GameMode),
		OsuVer:           20190828, //doesn't really matter
		MD5:              memory.MenuData.Bm.BeatmapMD5,
		PlayerName:       memory.GameplayData.Name + "(Failed)",
		BmChecksum:       "", //not needed for a functioning replay
		Hit300s:          uint16(memory.GameplayData.Hits.H300),
		Hit100s:          uint16(memory.GameplayData.Hits.H100),
		Hit50s:           uint16(memory.GameplayData.Hits.H50),
		HitGekis:         uint16(memory.GameplayData.Hits.HGeki),
		HitKatus:         uint16(memory.GameplayData.Hits.HKatu),
		HitMisses:        uint16(memory.GameplayData.Hits.H0),
		Score:            memory.GameplayData.Score,
		MaxCombo:         uint16(memory.GameplayData.Combo.Max),
		IsPerfect:        false,
		Mods:             memory.MenuData.Mods.AppliedMods,
		Lifebar:          "",                                     //not needed for a functioning replay
		DateTime:         time.Now().Unix()*10000000 + ticksUnix, //(C# DateTime)
		LengthReplayData: int32(len(compressed)),
		ReplayData:       []uint8(compressed),
		ScoreID:          0,
	}

	return replay
}

func gamemodeToStr(num int32) string {
	switch num {
	case 0:
		return "osuSTD"
	case 1:
		return "osuTaiko"
	case 2:
		return "osuCatch"
	case 3:
		return "osuMania"
	}
	return "Unsupported num"
}
