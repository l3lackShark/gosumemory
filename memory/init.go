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

	"github.com/l3lackShark/gosumemory/config"
	"github.com/l3lackShark/gosumemory/injctr"
	"github.com/l3lackShark/gosumemory/mem"
	"github.com/spf13/cast"
)

var osuProcessRegex = regexp.MustCompile(`.*osu!\.exe.*`)
var patterns staticAddresses
var tourneyPatterns []staticAddresses
var tourneySpecificPatterns []tourneyStaticAddresses
var tourneyMenuData []menuD
var tourneyManagerData tourneyD
var tourneyGameplayData []gameplayD
var tourneyAlwaysData []allTimesD

var menuData menuD
var songsFolderData songsFolderD
var gameplayData gameplayD
var alwaysData allTimesD
var mainMenuData mainMenuD
var resultsScreenData resultsScreenD

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
	if songsFolderData.SongsFolder == "Songs" || songsFolderData.SongsFolder == "CompatibilityContext" { //dirty hack to fix old stable offset
		return songsFolder, nil
	}
	return songsFolderData.SongsFolder, nil
}

func initBase() error {
	var err error
	isTournamentMode = false
	allProcs, err = mem.FindProcess(osuProcessRegex, "osu!lazer", "osu!framework")
	if err != nil {
		return err
	}
	process = allProcs[0]

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
		err = mem.Read(process,
			&patterns.PreSongSelectAddresses,
			&songsFolderData)
		if err != nil {
			return err
		}
		SongsFolderPath, err = resolveSongsFolder()
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("[MEMORY] Songs folder:", SongsFolderPath)
	pepath, err := process.ExecutablePath()
	if err != nil {
		panic(err)
	}
	SettingsData.Folders.Game = filepath.Dir(pepath)

	if menuData.PreSongSelectData.Status == 22 || len(allProcs) > 1 {
		fmt.Println("[MEMORY] Operating in tournament mode!")
		tourneyProcs, tourneyErr = resolveTourneyClients(allProcs)
		if tourneyErr != nil {
			return err
		}
		isTournamentMode = true
		tourneyPatterns = make([]staticAddresses, len(tourneyProcs))
		tourneySpecificPatterns = make([]tourneyStaticAddresses, len(tourneyProcs))
		TourneyData.IPCClients = make([]ipcClient, len(tourneyProcs))
		tourneyMenuData = make([]menuD, len(tourneyProcs))
		tourneyGameplayData = make([]gameplayD, len(tourneyProcs))
		tourneyAlwaysData = make([]allTimesD, len(tourneyProcs))
		for i, proc := range tourneyProcs {
			err = mem.ResolvePatterns(proc, &tourneyPatterns[i].PreSongSelectAddresses)
			if err != nil {
				return err
			}
			err = mem.Read(proc,
				&tourneyPatterns[i].PreSongSelectAddresses,
				&tourneyMenuData[i].PreSongSelectData)
			if err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("[MEMORY] Got osu!status addr for client #%d...", i))
			fmt.Println(fmt.Sprintf("[MEMORY] Resolving patterns for client #%d...", i))
			err = mem.ResolvePatterns(proc, &tourneyPatterns[i])
			if err != nil {
				return err
			}
			err = mem.ResolvePatterns(proc, &tourneySpecificPatterns[i])
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

	SettingsData.Folders.Songs = SongsFolderPath

	fmt.Println("[MEMORY] Got all patterns...")
	fmt.Println("WARNING: Mania pp calcualtion is experimental and only works if you choose mania gamemode in the SongSelect!")
	fmt.Println(fmt.Sprintf("Initialization complete, you can now visit http://%s or add it as a browser source in OBS", config.Config["serverip"]))
	DynamicAddresses.IsReady = true
	if cast.ToBool(config.Config["enabled"]) {
		err = injctr.Injct(process.Pid())
		if err != nil {
			log.Printf("Failed to inject into osu's process, in game overlay will be unavailable. %e\n", err)
		}
	} else {
		fmt.Println("[MEMORY] In-Game overlay is disabled, but could be enabled in config.ini!")
	}

	return nil
}
