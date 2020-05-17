package memory

import (
	"log"
	"strings"
	"time"

	"github.com/Andoryuuta/kiwi"
)

//UpdateTime Intervall between value updates
var UpdateTime int
var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")

//SongsFolderPath is full path to osu! Songs. Gets set automatically on Windows (through memory)
var SongsFolderPath string

//readHitErrorArray Gets an array of ints representing UnstableRate. (a little innacurate, shows values with 2 hitobjects delay)
func readHitErrorArray() ([]int32, error) {
	base, err := proc.ReadUint32(uintptr(DynamicAddresses.PlayContainer38 + 0x38))
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
			DynamicAddresses.IsReady = false
			for procerr != nil {
				proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
				log.Println("It seems that we lost the process, retrying!")
				time.Sleep(1 * time.Second)
			}
			DynamicAddresses.IsReady = false
			err := InitBase()
			for err != nil {
				err = InitBase()
				time.Sleep(1 * time.Second)
			}
		}
		if DynamicAddresses.IsReady == false {
			err := InitBase()
			for err != nil {
				err = InitBase()
				log.Println("Failure mid getting offsets, retrying!")
				time.Sleep(1 * time.Second)

			}
		}

		MenuData.OsuStatus, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
		if err != nil {
			log.Println("Could not get osuStatus Value, retrying")
			InitBase()
		}

		var tempBeatmapID uint32 = 0
		switch MenuData.OsuStatus {
		case 2:
			DynamicAddresses.PlayContainer38, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayContainer-0x4), 0x0, 0x38) //TODO: Should only be read once per map change
			if err != nil {
				log.Println(err)
			}
			xor1, err := proc.ReadUint32Ptr(uintptr(DynamicAddresses.PlayContainer38+0x1C), 0xC)
			xor2, err := proc.ReadUint32Ptr(uintptr(DynamicAddresses.PlayContainer38+0x1C), 0x8)
			if err != nil {
				log.Println(err, "xor")
			}
			accOffset, err := proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayContainer-0x4), 0x0, 0x48)
			GameplayData.Mods.AppliedMods = int32(xor1 ^ xor2)
			GameplayData.Combo.Current, err = proc.ReadInt32(uintptr(DynamicAddresses.PlayContainer38 + 0x90))
			GameplayData.Combo.Max, err = proc.ReadInt32(uintptr(DynamicAddresses.PlayContainer38 + 0x68))
			GameplayData.GameMode, err = proc.ReadInt32(uintptr(DynamicAddresses.PlayContainer38 + 0x64))
			GameplayData.Score, err = proc.ReadInt32(uintptr(DynamicAddresses.PlayContainer38 + 0x74))
			GameplayData.Hits.H100, err = proc.ReadInt16(uintptr(DynamicAddresses.PlayContainer38 + 0x84))
			GameplayData.Hits.H300, err = proc.ReadInt16(uintptr(DynamicAddresses.PlayContainer38 + 0x86))
			GameplayData.Hits.H50, err = proc.ReadInt16(uintptr(DynamicAddresses.PlayContainer38 + 0x88))
			GameplayData.Hits.H0, err = proc.ReadInt16(uintptr(DynamicAddresses.PlayContainer38 + 0x8E))
			GameplayData.Accuracy, err = proc.ReadFloat64(uintptr(accOffset + 0x14))
			MenuData.Bm.Time.PlayTime, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayTime+0x5), 0x0)
			GameplayData.Hits.HitErrorArray, err = readHitErrorArray()
			if err != nil {
				log.Println("GameplayData failure", err)
			}

		default: //This data available at all times
			DynamicAddresses.BeatmapAddr, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Base-0xC), 0x0)
			if err != nil {
				log.Println(err)
			}
			MenuData.Bm.BeatmapID, err = proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr + 0xC4))
			if err != nil {
				log.Println(err)
			}
			if tempBeatmapID != MenuData.Bm.BeatmapID { //On map change
				time.Sleep(time.Duration(UpdateTime) * time.Millisecond)
				MenuData.Bm.BeatmapSetID, err = proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr + 0xC8))
				beatmapStrOffset, err := proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr) + 0x7C)
				MenuData.Bm.BeatmapString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapStrOffset) + 0x8)
				beatmapBGStringOffset, err := proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr) + 0x68)
				MenuData.Bm.Path.BGPath, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapBGStringOffset) + 0x8)
				beatmapOsuFileStrOffset, err := proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr) + 0x8C)
				MenuData.Bm.Path.BeatmapOsuFileString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapOsuFileStrOffset) + 0x8)
				beatmapFolderStrOffset, err := proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr) + 0x74)
				MenuData.Bm.Path.BeatmapFolderString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapFolderStrOffset) + 0x8)
				MenuData.Bm.Stats.BeatmapAR, err = proc.ReadFloat32(uintptr(DynamicAddresses.BeatmapAddr + 0x2C))
				MenuData.Bm.Stats.BeatmapCS, err = proc.ReadFloat32(uintptr(DynamicAddresses.BeatmapAddr + 0x30))
				MenuData.Bm.Stats.BeatmapHP, err = proc.ReadFloat32(uintptr(DynamicAddresses.BeatmapAddr + 0x34))
				MenuData.Bm.Stats.BeatmapOD, err = proc.ReadFloat32(uintptr(DynamicAddresses.BeatmapAddr + 0x38))
				MenuData.Bm.Time.PlayTime, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayTime+0x5), 0x0)
				if err != nil {
					log.Println("MenuData failure")
				}

				if strings.HasSuffix(MenuData.Bm.Path.BeatmapOsuFileString, ".osu") == true && len(MenuData.Bm.Path.BGPath) > 0 {
					MenuData.Bm.Path.InnerBGPath = MenuData.Bm.Path.BeatmapFolderString + "/" + MenuData.Bm.Path.BGPath

				} else {
					log.Println("skipping bg reloading")
				}

				tempBeatmapID = MenuData.Bm.BeatmapID
			}

		}

		time.Sleep(time.Duration(UpdateTime) * time.Millisecond)
	}

}
