package memory

import (
	"regexp"

	"github.com/l3lackShark/gosumemory/mem"
)

func GetGameInstances() (*[]mem.Process, error) {
	osuProcessRegex := regexp.MustCompile(`.*osu!\.exe.*`)
	instances, err := mem.FindProcess(osuProcessRegex, "osu!lazer", "osu!framework")
	return &instances, err
}
