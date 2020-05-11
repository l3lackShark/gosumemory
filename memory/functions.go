package memory

import (
	"log"
	"strings"
	"time"

	"github.com/Andoryuuta/kiwi"
	"github.com/l3lackShark/gosumemory/values"
)

//UpdateTime Intervall between value updates
var UpdateTime int
var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")

//SongsFolderPath is full path to osu! Songs. Gets set automatically on Windows (through memory)
var SongsFolderPath string

//readHitErrorArray Gets an array of ints representing UnstableRate. (a little innacurate, shows values with 2 hitobjects delay)
func readHitErrorArray() ([]int32, error) {
	base, err := proc.ReadUint32(uintptr(values.DynamicAddresses.PlayContainer38 + 0x38))
	if err != nil {
		return nil, err
	}
	hitErrorStruct, err := proc.ReadUint32(uintptr(base + 0x4))
	if err != nil {
		return nil, err
	}
	leng, err := proc.ReadUint32(uintptr(base + 0xC))
	if err != nil {
		return nil, err
	}
	var buf32 []int32
	for i := 0x8; i < int(leng*0x4); i += 0x4 {
		value, err := proc.ReadInt32(uintptr(hitErrorStruct + uint32(i)))
		if err != nil {
			return nil, err
		}
		buf32 = append(buf32, value)
	}
	return buf32, nil
}

//Init the whole thing and get osu! memory values to start working with it.
func Init() {
	for {
		var err error
		proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
		if procerr != nil {
			for procerr != nil {
				proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
				log.Println("It seems that we lost the process, retrying!")
				time.Sleep(1 * time.Second)
			}
			values.MenuData.IsReady = false
			err := InitBase()
			for err != nil {
				err = InitBase()
				time.Sleep(1 * time.Second)
			}
		}
		if values.MenuData.IsReady == false {
			err := InitBase()
			for err != nil {
				err = InitBase()
				log.Println("Failure mid getting offsets, retrying!")
				time.Sleep(1 * time.Second)

			}
		}

		values.MenuData.OsuStatus, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
		if err != nil {
			log.Println("Could not get osuStatus Value!")
		}

		var tempBeatmapID uint32 = 0
		switch values.MenuData.OsuStatus {
		case 2:
			values.DynamicAddresses.PlayContainer38, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayContainer-0x4), 0x0, 0x38) //TODO: Should only be read once per map change
			if err != nil {
				log.Println(err)
			}
			xor1, err := proc.ReadUint32Ptr(uintptr(values.DynamicAddresses.PlayContainer38+0x1C), 0xC)
			xor2, err := proc.ReadUint32Ptr(uintptr(values.DynamicAddresses.PlayContainer38+0x1C), 0x8)
			if err != nil {
				log.Println(err, "xor")
			}
			accOffset, err := proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayContainer-0x4), 0x0, 0x48)
			values.GameplayData.AppliedMods = int32(xor1 ^ xor2)
			values.GameplayData.Combo, err = proc.ReadInt32(uintptr(values.DynamicAddresses.PlayContainer38 + 0x90))
			values.GameplayData.MaxCombo, err = proc.ReadInt32(uintptr(values.DynamicAddresses.PlayContainer38 + 0x68))
			values.GameplayData.GameMode, err = proc.ReadInt32(uintptr(values.DynamicAddresses.PlayContainer38 + 0x64))
			values.GameplayData.Score, err = proc.ReadInt32(uintptr(values.DynamicAddresses.PlayContainer38 + 0x74))
			values.GameplayData.Hit100c, err = proc.ReadInt16(uintptr(values.DynamicAddresses.PlayContainer38 + 0x84))
			values.GameplayData.Hit300c, err = proc.ReadInt16(uintptr(values.DynamicAddresses.PlayContainer38 + 0x86))
			values.GameplayData.Hit50c, err = proc.ReadInt16(uintptr(values.DynamicAddresses.PlayContainer38 + 0x88))
			values.GameplayData.HitMiss, err = proc.ReadInt16(uintptr(values.DynamicAddresses.PlayContainer38 + 0x8E))
			values.GameplayData.Accuracy, err = proc.ReadFloat64(uintptr(accOffset + 0x14))
			values.MenuData.PlayTime, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayTime+0x5), 0x0)
			values.GameplayData.HitErrorArray, err = readHitErrorArray()
			if err != nil {
				log.Println("GameplayData failure", err)
			}

		default: //This data available at all times
			values.MenuData.BeatmapAddr, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Base-0xC), 0x0)
			if err != nil {
				log.Println(err)
			}
			values.MenuData.BeatmapID, err = proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr + 0xC4))
			if err != nil {
				log.Println(err)
			}
			if tempBeatmapID != values.MenuData.BeatmapID { //On map change
				values.MenuData.BeatmapSetID, err = proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr + 0xC8))
				beatmapStrOffset, err := proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr) + 0x7C)
				values.MenuData.BeatmapString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapStrOffset) + 0x8)
				beatmapBGStringOffset, err := proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr) + 0x68)
				values.MenuData.BGPath, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapBGStringOffset) + 0x8)
				beatmapOsuFileStrOffset, err := proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr) + 0x8C)
				values.MenuData.BeatmapOsuFileString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapOsuFileStrOffset) + 0x8)
				beatmapFolderStrOffset, err := proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr) + 0x74)
				values.MenuData.BeatmapFolderString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapFolderStrOffset) + 0x8)
				values.MenuData.BeatmapAR, err = proc.ReadFloat32(uintptr(values.MenuData.BeatmapAddr + 0x2C))
				values.MenuData.BeatmapCS, err = proc.ReadFloat32(uintptr(values.MenuData.BeatmapAddr + 0x30))
				values.MenuData.BeatmapHP, err = proc.ReadFloat32(uintptr(values.MenuData.BeatmapAddr + 0x34))
				values.MenuData.BeatmapOD, err = proc.ReadFloat32(uintptr(values.MenuData.BeatmapAddr + 0x38))
				values.MenuData.PlayTime, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayTime+0x5), 0x0)
				if err != nil {
					log.Println("MenuData failure")
				}

				if strings.HasSuffix(values.MenuData.BeatmapOsuFileString, ".osu") == true && len(values.MenuData.BGPath) > 0 {
					values.MenuData.InnerBGPath = values.MenuData.BeatmapFolderString + "/" + values.MenuData.BGPath

				} else {
					log.Println("skipping bg reloading")
				}

				tempBeatmapID = values.MenuData.BeatmapID
			}

		}

		time.Sleep(time.Duration(UpdateTime) * time.Millisecond)
	}

}
