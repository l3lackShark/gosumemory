package patterns

import (
	"log"
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
			values.MenuData.BeatmapSetID, err = proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr + 0xC8))
			if err != nil {
				log.Println(err)
			}
			beatmapStrOffset, err := proc.ReadUint32(uintptr(values.MenuData.BeatmapAddr) + 0x7C)
			values.MenuData.BeatmapString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapStrOffset) + 0x8)
			if err != nil {
				log.Println(err)
			}

		}
		time.Sleep(100 * time.Millisecond)
	}

}
