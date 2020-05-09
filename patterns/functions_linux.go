// +build linux

package patterns

import (
	"fmt"
	"log"
	"os"

	"github.com/Andoryuuta/kiwi"
)

var pid uint64

//ResolveOsuStatus Gets osuStatusValue to start working with it.
func ResolveOsuStatus() int32 {

	var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")

	if procerr != nil {
		log.Println("osu! is not running!")
		return -1
	}

	maps, err := readMaps(int(proc.PID))
	if err != nil {
		log.Println("Please provide process/Process error!")
		return -2

	}
	mem, err := os.Open(fmt.Sprintf("/proc/%d/mem", int(proc.PID)))
	if err != nil {
		log.Println("Coud not open /proc (missing sudo?")
		return -3

	}
	defer mem.Close()

	osuStatusBase, err := scan(mem, maps, osuSignatures.status)
	if err != nil {
		log.Println("Could not get signature!")
		return -4

	}
	result, err := proc.ReadUint32Ptr(uintptr(osuStatusBase-0x4), 0x0)
	if err != nil {
		log.Println("Could not get osuStatus Value!")
		return -5
	}
	return int32(result)

}
