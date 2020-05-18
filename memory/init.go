package memory

import (
	"fmt"
	"regexp"
	"time"

	"github.com/k0kubun/pp"

	"github.com/l3lackShark/gosumemory/mem"

	"github.com/Andoryuuta/kiwi"
	"github.com/spf13/cast"
)

//InitBase initializes base static addresses. (In hopes to deprecate C#)
func InitBase() error {
	var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	if procerr != nil {
		return procerr
	}
	//Migration to the new win32 wrapper by tdeo:
	re := regexp.MustCompile(`.*osu!\.exe.*`)
	newproc, newprocerr := mem.FindProcess(re)
	if newprocerr != nil {
		pp.Println("There was an error in the attempt to find a process!.. ", newprocerr)
		return newprocerr
	}
	pp.Println(newproc)
	osuStatusAddr, err := mem.Scan(newproc, "48 83 F8 04 73 1E")
	if err != nil {
		return err
	}
	osuStaticAddresses.Status = cast.ToUint32(osuStatusAddr)
	fmt.Printf("OsuStatusAddr: 0x%x\n", osuStaticAddresses.Status)
	osuStatus, err := proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
	if err != nil {
		return err
	}
	for osuStatus == 0 {
		fmt.Println("Please go to song select in order to proceed!")
		time.Sleep(500 * time.Millisecond)
		osuStatus, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
		if err != nil {
			return err
		}

	}

	var patterns NewPatterns
	err = mem.ResolvePatterns(newproc, &patterns)
	if err != nil {
		return err
	}
	pp.Println(patterns)
	osuStaticAddresses.BPM = cast.ToUint32(patterns.BPM)
	osuStaticAddresses.Base = cast.ToUint32(patterns.Base)
	osuStaticAddresses.InMenuMods = cast.ToUint32(patterns.InMenuMods)
	osuStaticAddresses.PlayTime = cast.ToUint32(patterns.PlayTime)
	osuStaticAddresses.PlayContainer = cast.ToUint32(patterns.PlayContainer)
	osuStaticAddresses.LeaderBoard = cast.ToUint32(patterns.LeaderBoard + 0x1)
	DynamicAddresses.IsReady = true
	return nil
}
