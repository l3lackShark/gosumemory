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
			isReady = false
			err := InitBase()
			for err != nil {
				err = InitBase()
				log.Println("It seems that we lost the process, retrying!")
				time.Sleep(1 * time.Second)
			}
		}
		if isReady == false {
			InitBase()
		}

		values.OsuData.OsuStatus, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
		if err != nil {
			log.Println("Could not get osuStatus Value!")
		}
		switch values.OsuData.OsuStatus {
		case 2:
			values.OsuData.PlayContainer38, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayContainer-0x4), 0x0, 0x38)
			if err != nil {
				log.Println(err)
			}
		default:
			values.OsuData.BeatmapAddr, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Base-0xC), 0x0)
			if err != nil {
				log.Println(err)
			}
			values.OsuData.BeatMapID, err = proc.ReadUint32(uintptr(values.OsuData.BeatmapAddr + 0xC4))
			if err != nil {
				log.Println(err)
			}

		}
		time.Sleep(100 * time.Millisecond)
	}

}
