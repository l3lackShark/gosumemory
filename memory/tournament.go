package memory

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/mem"
)

var tourneyManagerID = 0

//initTournement should be called on tournament manager
func initTournement() error {

	//read tournament.cfg to check how many clients we are expecting
	osuExecutablePath, err := process[0].ExecutablePath()
	if err != nil {
		return err
	}
	if !strings.Contains(osuExecutablePath, `:\`) {
		log.Println("Automatic executable path finder has failed. The program will now ext. GOT: ", osuExecutablePath)
		return errors.New("osu! executable was not found")
	}
	cfgFile, err := os.Open(filepath.Join(filepath.Dir(osuExecutablePath), "tournament.cfg"))
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	var totalClients int
	scanner := bufio.NewScanner(cfgFile)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "TeamSize") {
			teamSize, err := strconv.Atoi(scanner.Text()[len(scanner.Text())-1:])
			if err != nil {
				return err
			}
			totalClients = teamSize * 2
			fmt.Println("Total expected amount of tournament clients:", totalClients)
			break
		}
	}

	if totalClients == 0 {
		return errors.New("total clients is 0")
	}
	fmt.Println("[TOURNAMENT] Awaiting all clients to load...")
	for len(process) != totalClients+1 { //+1 is Tournament Manager
		process, err = mem.FindProcess(osuProcessRegex)
		if err != nil {
			return err
		}
		fmt.Println("[TOURNAMENT] Loaded", len(process), "clients..", "wating for", totalClients-len(process), "more...")
		time.Sleep(500 * time.Millisecond)
	}

	menuData = make([]menuD, len(process))
	patterns = make([]staticAddresses, len(process))
	gameplayData = make([]gameplayD, len(process))
	alwaysData = make([]allTimesD, len(process))

	pp.Println(process)
	for i := range process {
		err = mem.ResolvePatterns(process[i], &patterns[i].PreSongSelectAddresses)
		if err != nil {
			return err
		}
		err = mem.Read(process[i],
			&patterns[i].PreSongSelectAddresses,
			&menuData[i].PreSongSelectData)
		if err != nil {
			return err
		}
		if menuData[i].PreSongSelectData.Status == 22 {
			tourneyManagerID = i
		}
		fmt.Println(process[i].Pid())
		if i != tourneyManagerID {
			err = mem.ResolvePatterns(process[i], &patterns[i])
			if err != nil {
				return err
			}
		}

	}

	return nil
}
