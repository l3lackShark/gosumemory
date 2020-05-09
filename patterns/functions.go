package patterns

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Andoryuuta/kiwi"
	"github.com/l3lackShark/gosumemory/values"
)

var isReady bool = false

//Init the whole thing and get osu! memory values to start working with it.
func Init() {
	for {
		var err error
		var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
		if procerr != nil {
			for procerr != nil {
				proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
				log.Println("It seems that we lost the process, retrying!")
				time.Sleep(1 * time.Second)
			}
			isReady = false
			err := InitBase()
			for err != nil {
				err = InitBase()
				log.Println("It seems that we lost the process, retrying!(2)")
				time.Sleep(1 * time.Second)
			}
		}
		if isReady == false {
			InitBase()
		}

		values.MenuData.OsuStatus, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
		if err != nil {
			log.Println("Could not get osuStatus Value!")
		}
		var tempBeatmapID uint32 = 0
		switch values.MenuData.OsuStatus {
		case 2:
			// values.MenuData.PlayContainer38, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayContainer-0x4), 0x0, 0x38)
			// if err != nil {
			// 	log.Println(err)
			// }
		default:
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
				if err != nil {
					log.Println(err)
				}
				beatmapStrOffset, err := proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr) + 0x7C)
				values.MenuData.BeatmapString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapStrOffset) + 0x8)
				if err != nil {
					log.Println(err)
				}
				beatmapBGStringOffset, err := proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr) + 0x68)
				values.MenuData.BGPath, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapBGStringOffset) + 0x8)
				if err != nil {
					log.Println(err)
				}
				beatmapOsuFileStrOffset, err := proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr) + 0x8C)
				values.MenuData.BeatmapOsuFileString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapOsuFileStrOffset) + 0x8)
				if err != nil {
					log.Println(err)
				}
				beatmapFolderStrOffset, err := proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr) + 0x74)
				values.MenuData.BeatmapFolderString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapFolderStrOffset) + 0x8)
				if err != nil {
					log.Println(err)
				}
				values.MenuData.BeatmapAR, err = proc.ReadFloat32(uintptr(values.MenuData.BeatmapAddr + 0x2C))
				if err != nil {
					log.Println(err)
				}
				values.MenuData.BeatmapCS, err = proc.ReadFloat32(uintptr(values.MenuData.BeatmapAddr + 0x30))
				if err != nil {
					log.Println(err)
				}
				values.MenuData.BeatmapHP, err = proc.ReadFloat32(uintptr(values.MenuData.BeatmapAddr + 0x34))
				if err != nil {
					log.Println(err)
				}
				values.MenuData.BeatmapOD, err = proc.ReadFloat32(uintptr(values.MenuData.BeatmapAddr + 0x38))
				if err != nil {
					log.Println(err)
				}
				if strings.HasSuffix(values.MenuData.BeatmapOsuFileString, ".osu") == true && len(values.MenuData.BGPath) > 0 {
					values.MenuData.InnerBGPath = values.MenuData.BeatmapFolderString + "/" + values.MenuData.BGPath

				} else {
					fmt.Println("skipping bg reloading")
				}

				tempBeatmapID = values.MenuData.BeatmapID
			}

		}

		time.Sleep(100 * time.Millisecond)
	}

}
