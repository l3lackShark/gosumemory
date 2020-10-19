package memory

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
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

func resolveSongsFolder() (string, error) {
	var err error
	osuExecutablePath, err := process.ExecutablePath()
	if err != nil {
		return "", err
	}
	if !strings.Contains(osuExecutablePath, `:\`) {
		log.Println("Automatic executable path finder has failed. Please try again or manually specify it. (see --help) GOT: ", osuExecutablePath)
		time.Sleep(5 * time.Second)
		return "", errors.New("osu! executable was not found")
	}
	rootFolder := strings.TrimSuffix(osuExecutablePath, "osu!.exe")
	songsFolder := filepath.Join(rootFolder, "Songs")
	if menuData.PreSongSelectData.SongsFolder == "Songs" {
		return songsFolder, nil
	}
	return menuData.PreSongSelectData.SongsFolder, nil
}

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
	if runtime.GOOS == "windows" && SongsFolderPath == "auto" {
		SongsFolderPath, err = resolveSongsFolder()
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("[MEMORY] Songs folder:", SongsFolderPath)

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

	DynamicAddresses.IsReady = true

	return nil
}
