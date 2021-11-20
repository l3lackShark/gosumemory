package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/helpers"
	"github.com/l3lackShark/gosumemory/mem"
	"github.com/l3lackShark/gosumemory/memory"
)

const (
	osuStatusMainMenu   = 0
	osuStatusPlaying    = 2
	osuStatusSongSelect = 5
)

type agent struct {
	config
	clients []struct {
		instance        mem.Process
		patterns        memory.StaticAddresses
		menuData        memory.MenuD
		songsFolderData memory.SongsFolderD
	}
}

type config struct {
	updateTime  int
	songsFolder string
}

func main() {
	//some init config stuff will go here
	a := agent{
		config: config{
			updateTime: 100,
		},
	}
	a.runMainLoop()
	// a.patterns.hello
}

func (a *agent) runMainLoop() {
StartOver:
	fmt.Println("Waiting for osu! to launch...")
	instances, err := memory.GetGameInstances()
	for err != nil {
		helpers.Sleep(a.config.updateTime)
		instances, err = memory.GetGameInstances()
	}
	instanceCount := len(*instances)

	if instanceCount > 1 {
		fmt.Println("Operating in tourney mode") //TODO: Add tourney support
	}

	a.clients = make([]struct {
		instance        mem.Process
		patterns        memory.StaticAddresses
		menuData        memory.MenuD
		songsFolderData memory.SongsFolderD
	}, instanceCount)

	for i := 0; i < instanceCount; i++ {
		a.clients[i].instance = (*instances)[i]
		err = mem.ResolvePatterns(a.clients[i].instance, &a.clients[i].patterns.PreSongSelectAddresses)
		if err != nil {
			log.Println(err)
			goto StartOver
		}
		err = mem.Read(a.clients[i].instance,
			&a.clients[i].patterns.PreSongSelectAddresses,
			&a.clients[i].menuData.PreSongSelectData)
		if err != nil {
			log.Println(err)
			goto StartOver
		}
		fmt.Println("[MEMORY] Resolving patterns...")
		err = mem.ResolvePatterns(a.clients[i].instance, &a.clients[i].patterns)
		if err != nil {
			log.Println(err)
			goto StartOver
		}
	}
	pp.Println("Got the game\n", a.clients[0].menuData.PreSongSelectData.Status)

	//TODO: make this flexible with config
	a.config.songsFolder, err = a.resolveSongsFolder()
	if err != nil {
		log.Println(err)
		goto StartOver
	}
	fmt.Println("Songs Folder:", a.config.songsFolder)

	//run the main loop, goto startover if we loose the instance
	for {
		for i := 0; i < instanceCount; i++ {
			err := mem.Read(a.clients[i].instance,
				&a.clients[i].patterns.PreSongSelectAddresses,
				&a.clients[i].menuData.PreSongSelectData)
			if err != nil {
				log.Println("It appears that we lost the precess, retrying! ERROR:", err)
				goto StartOver
			}
			switch a.clients[i].menuData.PreSongSelectData.Status {
			case osuStatusPlaying:
				fmt.Println("I'm playing :)")

			}
		}
		helpers.Sleep(a.updateTime)
	}
}

func (a *agent) resolveSongsFolder() (string, error) {
	if runtime.GOOS == "windows" {
		var err error

		err = mem.Read(a.clients[0].instance,
			&a.clients[0].patterns.PreSongSelectAddresses,
			&a.clients[0].songsFolderData)
		if err != nil {
			return "", err
		}

		osuExecutablePath, err := a.clients[0].instance.ExecutablePath()
		if err != nil {
			return "", err
		}
		if !strings.Contains(osuExecutablePath, `:\`) {
			log.Println("Automatic executable path finder has failed. Please try again or manually specify it. (see --help) GOT: ", osuExecutablePath)
			return "", fmt.Errorf("osu! executable was not found")
		}
		rootFolder := strings.TrimSuffix(osuExecutablePath, "osu!.exe")
		songsFolder := filepath.Join(rootFolder, "Songs")
		if a.clients[0].songsFolderData.SongsFolder == "Songs" {
			return songsFolder, nil
		}
		return a.clients[0].songsFolderData.SongsFolder, nil
	}
	return "", fmt.Errorf("unsupported OS") //TODO: add Linux
}
