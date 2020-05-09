// +build linux

package patterns

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Andoryuuta/kiwi"
)

var pid uint64
var osuStatusBase uint32

//OsuStaticAddresses (should be updated every client restart)
var OsuStaticAddresses = StaticAddresses{}

//ResolveOsuStatus Gets osuStatusValue to start working with it.
func ResolveOsuStatus() int32 {

	var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	if procerr != nil {
		log.Println("osu! is not running!")

	}

	result, err := proc.ReadUint32Ptr(uintptr(osuStatusBase-0x4), 0x0)
	if err != nil {
		log.Println("Could not get osuStatus Value!")
		return -5
	}
	return int32(result)

}

//InitBase initializes base static addresses.
func InitBase() error {
	var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	if procerr != nil {
		InitBase()
	}
	maps, err := readMaps(int(proc.PID))
	if err != nil {
		log.Println("Process error!")
		return err

	}
	mem, err := os.Open(fmt.Sprintf("/proc/%d/mem", int(proc.PID)))
	if err != nil {
		log.Println("Coud not open /proc")
		return err

	}
	defer mem.Close()

	OsuStaticAddresses.Status, err = scan(mem, maps, osuSignatures.status)
	if err != nil {
		return err
	}
	osuStatus, err := proc.ReadUint32Ptr(uintptr(OsuStaticAddresses.Status-0x4), 0x0)
	if err != nil {
		return err
	}
	for osuStatus != 5 {
		fmt.Println("Please go to song select in order to proceed!")
		time.Sleep(500 * time.Millisecond)
		osuStatus, err = proc.ReadUint32Ptr(uintptr(OsuStaticAddresses.Status-0x4), 0x0)
		if err != nil {
			return err
		}

	}
	OsuStaticAddresses.BPM, err = scan(mem, maps, osuSignatures.bpm)
	if err != nil {
		return err
	}
	OsuStaticAddresses.Base, err = scan(mem, maps, osuSignatures.base)
	if err != nil {
		return err
	}
	OsuStaticAddresses.InMenuMods, err = scan(mem, maps, osuSignatures.inMenuMods)
	if err != nil {
		return err
	}
	OsuStaticAddresses.PlayTime, err = scan(mem, maps, osuSignatures.playTime)
	if err != nil {
		return err
	}
	OsuStaticAddresses.PlayContainer, err = scan(mem, maps, osuSignatures.playContainer)
	if err != nil {
		return err
	}

	return nil
}
