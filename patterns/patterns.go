package patterns

import (
	"fmt"
	"log"
	"os"

	"github.com/Andoryuuta/kiwi"
)

//Patterns is Base osu signatures stuct
type Patterns struct {
	status        string
	bpm           string
	base          string
	inMenuMods    string
	playTime      string
	playContainer string
}

//ResolveOsuStatus Gets osuStatusValue to start working with it.
func ResolveOsuStatus() int32 {
	var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	if procerr != nil {
		log.Fatalln("osu! is not running!")
	}
	OsuSignatures := Patterns{
		status: "48 83 F8 04 73 1E",
	}
	maps, err := readMaps(int(proc.PID))
	if err != nil {
		log.Fatalln("Please provide process/Process error!")
	}
	mem, err := os.Open(fmt.Sprintf("/proc/%d/mem", int(proc.PID)))
	if err != nil {
		fmt.Println("Coud not open /proc (missing sudo?")

	}
	defer mem.Close()

	osuStatusValue, err := scan(mem, maps, OsuSignatures.status)
	if err != nil {
		fmt.Println("Could not get signature!")

	}
	result, err := proc.ReadInt32(uintptr(osuStatusValue) - 0x4)
	if err != nil {
		log.Println("Could not get osuStatus Value!")
	}
	return result

}
