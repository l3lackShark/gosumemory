package memory

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
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

	windowsUser, err := user.Current()
	if err != nil {
		fmt.Println("[MEMORY] Could not find Windows Username, Please manually specify SongsFolder. (see --help)")
		time.Sleep(5 * time.Second)
		return "", errors.New("Windwos Username was not found")
	}
	userName := strings.Split(windowsUser.Username, "\\")
	configPath := filepath.Join(rootFolder, "osu!."+userName[1]+".cfg")
	config, err := os.Open(configPath)
	if err != nil {
		log.Println("[MEMORY] Could not find your osu!.cfg, Please manually specify SongsFolder. (see --help)")
		time.Sleep(5 * time.Second)
		return "", errors.New("osu!.cfg was not found")
	}
	defer config.Close()
	scanner := bufio.NewScanner(config)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "BeatmapDirectory") {
			split := strings.Split(scanner.Text(), "=")
			split[1] = strings.TrimSpace(split[1])
			if split[1] == "Songs" {
				if _, err := os.Stat(songsFolder); !os.IsNotExist(err) {
					fmt.Printf("[MEMORY] Songs Folder Path: %s\n", songsFolder)
					return songsFolder, nil
				}
			}
			fmt.Println("[MEMORY] It appears that you moved your Songs directory to somehere else. Attempting to find it...")
			if strings.Contains(split[1], `:\`) == false {
				fmt.Println("[MEMORY] This is very confusing... non-default Songs folder but still in the same directory? Encountering...")
				songsFolder = filepath.Join(rootFolder, split[1])
				if _, err := os.Stat(songsFolder); !os.IsNotExist(err) {
					fmt.Printf("[MEMORY] Songs Folder Path: %s\n", songsFolder)
					return songsFolder, nil
				}
				log.Println("[MEMORY] Automatic Songs folder finder has failed. Please try again or manually specify it. (see --help) GOT: ", split[1], "This probably means that you renamed your Songs folder to something else but it's still in the same folder as the main executable. An attempt was made to encounter this, but it failed.")
				time.Sleep(15 * time.Second)
				return "", errors.New("Songs Folder was not found")

			} else {
				if _, err := os.Stat(split[1]); !os.IsNotExist(err) {
					fmt.Printf("[MEMORY] Songs Folder Path: %s\n", split[1])
					return split[1], nil
				}
				log.Println("[MEMORY] Automatic Songs folder finder has failed. Please try again or manually specify it. (see --help) GOT: ", split[1])
				time.Sleep(5 * time.Second)
				return "", errors.New("Songs Folder was not found")

			}
		}
	}
	log.Println("[MEMORY] Songs Folder was not found. Please try again or manually specify it. (see --help)")
	time.Sleep(10 * time.Second)
	return "", errors.New("Songs Folder was not found")
}

func initBase() error {
	process, err := mem.FindProcess(osuProcessRegex)
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" && SongsFolderPath == "auto" {
		SongsFolderPath, err = resolveSongsFolder()
		if err != nil {
			log.Fatalln(err)
		}
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

	DynamicAddresses.IsReady = true

	return nil
}
