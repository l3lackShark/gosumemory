package memory

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/l3lackShark/gosumemory/mem"
)

var osuProcessRegex = regexp.MustCompile(`.*osu!\.exe.*`)
var patterns staticAddresses

var menuData menuD
var gameplayData gameplayD
var alwaysData allTimesD

func initBase() error {
	process, err := mem.FindProcess(osuProcessRegex)
	if err != nil {
		return err
	}

	err = mem.ResolvePatterns(process, &patterns.PreSongSelectAddresses)
	if err != nil {
		return err
	}

	err = mem.Read(process,
		&patterns.PreSongSelectAddresses,
		&menuData.PreSongSelectData)
	if err != nil {
		return err
	}
	fmt.Println("[MEMORY] Got osu!status addr...")
	if menuData.Status == 0 {
		log.Println("Please go to song select to proceed!")
		for menuData.Status == 0 {
			time.Sleep(100 * time.Millisecond)
			err := mem.Read(process,
				&patterns.PreSongSelectAddresses,
				&menuData.PreSongSelectData)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("[MEMORY] Resolving patterns...")
	err = mem.ResolvePatterns(process, &patterns)
	if err != nil {
		return err
	}
	fmt.Println("[MEMORY] Got all patterns...")
	if runtime.GOOS == "windows" && SongsFolderPath == "auto" {
		SongsFolderPath, err = mem.ReadString(process, int64(menuData.PreSongSelectData.SongsFolderPathAddr+0x18), 0)
		if err != nil || strings.Contains(SongsFolderPath, `:\`) == false {
			log.Println("Automatic Songs folder finder has failed. Please try again or manually specify it. (see --help) GOT: ", SongsFolderPath, err)
			time.Sleep(5 * time.Second)
			log.Fatalln(err)
		}
	}
	fmt.Printf("[MEMORY] Songs Folder Path: %s\n", SongsFolderPath)

	DynamicAddresses.IsReady = true

	return nil
}
