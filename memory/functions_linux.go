// +build linux

package memory

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Andoryuuta/kiwi"
	"github.com/l3lackShark/gosumemory/values"
)

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

	osuStaticAddresses.Status, err = scan(mem, maps, osuSignatures.status)
	if err != nil {
		return err
	}
	osuStatus, err := proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
	if err != nil {
		return err
	}
	for osuStatus != 5 {
		fmt.Println("Please go to song select in order to proceed!")
		time.Sleep(500 * time.Millisecond)
		osuStatus, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
		if err != nil {
			return err
		}

	}
	osuStaticAddresses.BPM, err = scan(mem, maps, osuSignatures.bpm)
	if err != nil {
		return err
	}
	osuStaticAddresses.Base, err = scan(mem, maps, osuSignatures.base)
	if err != nil {
		return err
	}
	osuStaticAddresses.InMenuMods, err = scan(mem, maps, osuSignatures.inMenuMods)
	if err != nil {
		return err
	}
	osuStaticAddresses.PlayTime, err = scan(mem, maps, osuSignatures.playTime)
	if err != nil {
		return err
	}
	osuStaticAddresses.PlayContainer, err = scan(mem, maps, osuSignatures.playContainer)
	if err != nil {
		return err
	}
	values.MenuData.IsReady = true
	return nil
}
