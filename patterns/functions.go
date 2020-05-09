package patterns

import (
	"log"
	"time"

	"github.com/Andoryuuta/kiwi"
	"github.com/l3lackShark/gosumemory/values"
)

var isReady bool = false

//Init the whole thing and get osuStatusValue to start working with it.
func Init() {
	for {
		var err error
		var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
		if procerr != nil {
			isReady = false
			err := InitBase()
			for err != nil {
				err = InitBase()
			}
		}
		if isReady == false {
			InitBase()
		} else {
			isReady = true
		}

		values.OsuData.OsuStatus, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
		if err != nil {
			log.Println("Could not get osuStatus Value!")
		}
		time.Sleep(500 * time.Millisecond)
	}

}
