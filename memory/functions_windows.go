//+build windows

package memory

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Andoryuuta/kiwi"
	"github.com/spf13/cast"
)

//InitBase initializes base static addresses. (In hopes to deprecate C#)
func InitBase() error {
	var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	if procerr != nil {
		InitBase()
	}

	cmd, err := exec.Command("deps/OsuStatusAddr.exe").Output()
	if err != nil {
		return err
	}
	outStr := string(cmd)
	outStr = strings.Replace(outStr, "\n", "", -1)
	outStr = strings.Replace(outStr, "\r", "", -1)
	osuStaticAddresses.Status = cast.ToUint32(outStr)
	fmt.Printf("OsuStatusAddr: %x\n", osuStaticAddresses.Status)
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
	cmd, err = exec.Command("deps/OsuBPMAddr.exe").Output()
	if err != nil {
		return err
	}
	outStr = string(cmd)
	outStr = strings.Replace(outStr, "\n", "", -1)
	outStr = strings.Replace(outStr, "\r", "", -1)
	osuStaticAddresses.BPM = cast.ToUint32(outStr)
	fmt.Printf("OsuBPMAddr: %x\n", osuStaticAddresses.BPM)

	cmd, err = exec.Command("deps/OsuBaseAddr.exe").Output()
	if err != nil {
		return err
	}
	outStr = string(cmd)
	outStr = strings.Replace(outStr, "\n", "", -1)
	outStr = strings.Replace(outStr, "\r", "", -1)
	osuStaticAddresses.Base = cast.ToUint32(outStr)
	fmt.Printf("OsuBaseAddr: %x\n", osuStaticAddresses.Base)

	cmd, err = exec.Command("deps/InMenuAppliedModsAddr.exe").Output()
	if err != nil {
		return err
	}
	outStr = string(cmd)
	outStr = strings.Replace(outStr, "\n", "", -1)
	outStr = strings.Replace(outStr, "\r", "", -1)
	osuStaticAddresses.InMenuMods = cast.ToUint32(outStr)
	fmt.Printf("OsuInMenuModsAddr: %x\n", osuStaticAddresses.InMenuMods)

	cmd, err = exec.Command("deps/OsuPlayTimeAddr.exe").Output()
	if err != nil {
		return err
	}
	outStr = string(cmd)
	outStr = strings.Replace(outStr, "\n", "", -1)
	outStr = strings.Replace(outStr, "\r", "", -1)
	osuStaticAddresses.PlayTime = cast.ToUint32(outStr)
	fmt.Printf("OsuPlayTimeAddr: %x\n", osuStaticAddresses.PlayTime)

	cmd, err = exec.Command("deps/OsuplayContainer.exe").Output()
	if err != nil {
		return err
	}
	outStr = string(cmd)
	outStr = strings.Replace(outStr, "\n", "", -1)
	outStr = strings.Replace(outStr, "\r", "", -1)
	osuStaticAddresses.PlayContainer = cast.ToUint32(outStr)
	fmt.Printf("OsuPlayContainerAddr: %x\n", osuStaticAddresses.PlayContainer)

	DynamicAddresses.IsReady = true
	proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	return nil
}
