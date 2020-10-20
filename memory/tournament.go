//+build windows

package memory

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/l3lackShark/gosumemory/mem"
)

func resolveTourneyClients(procs []mem.Process) ([]mem.Process, error) {
	if runtime.GOOS != "windows" {
		panic("Tournament client is not yet implemented on Linux")
	}
	//read tournament.cfg to check how many clients we are expecting
	osuExecutablePath, err := procs[0].ExecutablePath()
	if err != nil {
		return nil, err
	}
	if !strings.Contains(osuExecutablePath, `:\`) {
		log.Println("Automatic executable path finder has failed. The program will now ext. GOT: ", osuExecutablePath)
		return nil, errors.New("osu! executable was not found")
	}
	cfgFile, err := os.Open(filepath.Join(filepath.Dir(osuExecutablePath), "tournament.cfg"))
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	var totalClients int
	scanner := bufio.NewScanner(cfgFile)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "TeamSize") {
			teamSize, err := strconv.Atoi(scanner.Text()[len(scanner.Text())-1:])
			if err != nil {
				return nil, err
			}
			totalClients = teamSize * 2
			fmt.Println("Total expected amount of tournament clients:", totalClients)
			break
		}
	}

	if totalClients == 0 {
		return nil, errors.New("total clients is 0")
	}
	fmt.Println("[TOURNAMENT] Awaiting all clients to load...")
	var tourneyClients []mem.Process
	for len(procs) != totalClients+1 {
		procs, err = mem.FindProcess(osuProcessRegex)
		if err != nil {
			return nil, err
		}
	}
	for i := range procs {
		if i > len(procs)-2 {
			break
		}
		client, err := mem.FindWindow(fmt.Sprintf("Tournament Client %d", i))
		if err != nil {
			return nil, fmt.Errorf("Error getting clients, '%s'", err)
		}
		for _, proc := range procs {
			if int32(proc.Pid()) == mem.GetWindowThreadProcessID(client) {
				tourneyClients = append(tourneyClients, proc)
				break
			}
		}
	}
	return tourneyClients, nil
}

func getTourneyGamplayData(proc mem.Process, iterator int) {
	mem.Read(proc, &tourneyPatterns[iterator], &tourneyGameplayData[iterator])
	TourneyData.Clients[iterator].Gameplay.Combo.Current = tourneyGameplayData[iterator].Combo
	TourneyData.Clients[iterator].Gameplay.Combo.Max = tourneyGameplayData[iterator].MaxCombo
	TourneyData.Clients[iterator].Gameplay.GameMode = tourneyGameplayData[iterator].Mode
	TourneyData.Clients[iterator].Gameplay.Score = tourneyGameplayData[iterator].Score
	TourneyData.Clients[iterator].Gameplay.Hits.H100 = tourneyGameplayData[iterator].Hit100
	TourneyData.Clients[iterator].Gameplay.Hits.HKatu = tourneyGameplayData[iterator].HitKatu
	TourneyData.Clients[iterator].Gameplay.Hits.H200M = tourneyGameplayData[iterator].Hit200M
	TourneyData.Clients[iterator].Gameplay.Hits.H300 = tourneyGameplayData[iterator].Hit300
	TourneyData.Clients[iterator].Gameplay.Hits.HGeki = tourneyGameplayData[iterator].HitGeki
	TourneyData.Clients[iterator].Gameplay.Hits.H50 = tourneyGameplayData[iterator].Hit50
	TourneyData.Clients[iterator].Gameplay.Hits.H0 = tourneyGameplayData[iterator].HitMiss
	if TourneyData.Clients[iterator].Gameplay.Combo.Temp > TourneyData.Clients[iterator].Gameplay.Combo.Max {
		TourneyData.Clients[iterator].Gameplay.Combo.Temp = 0
	}
	if TourneyData.Clients[iterator].Gameplay.Combo.Current < TourneyData.Clients[iterator].Gameplay.Combo.Temp && TourneyData.Clients[iterator].Gameplay.Hits.H0Temp == TourneyData.Clients[iterator].Gameplay.Hits.H0 {
		TourneyData.Clients[iterator].Gameplay.Hits.HSB++
	}
	TourneyData.Clients[iterator].Gameplay.Hits.H0Temp = TourneyData.Clients[iterator].Gameplay.Hits.H0
	TourneyData.Clients[iterator].Gameplay.Combo.Temp = TourneyData.Clients[iterator].Gameplay.Combo.Current
	MenuData.Mods.AppliedMods = int32(tourneyGameplayData[iterator].ModsXor1 ^ tourneyGameplayData[iterator].ModsXor1)
	TourneyData.Clients[iterator].Gameplay.Accuracy = tourneyGameplayData[iterator].Accuracy
	TourneyData.Clients[iterator].Gameplay.Hp.Normal = tourneyGameplayData[iterator].PlayerHP
	TourneyData.Clients[iterator].Gameplay.Hp.Smooth = tourneyGameplayData[iterator].PlayerHPSmooth
	TourneyData.Clients[iterator].Gameplay.Name = tourneyGameplayData[iterator].PlayerName
	MenuData.Mods.AppliedMods = int32(tourneyGameplayData[iterator].ModsXor1 ^ tourneyGameplayData[iterator].ModsXor2)
	if MenuData.Mods.AppliedMods == 0 {
		MenuData.Mods.PpMods = "NM"
	} else {
		MenuData.Mods.PpMods = Mods(tourneyGameplayData[iterator].ModsXor1 ^ tourneyGameplayData[iterator].ModsXor2).String()
	}
	if TourneyData.Clients[iterator].Gameplay.Combo.Max > 0 {
		TourneyData.Clients[iterator].Gameplay.Hits.HitErrorArray = tourneyGameplayData[iterator].HitErrors
		baseUR, _ := calculateUR(TourneyData.Clients[iterator].Gameplay.Hits.HitErrorArray)
		if strings.Contains(MenuData.Mods.PpMods, "DT") || strings.Contains(MenuData.Mods.PpMods, "NC") {
			TourneyData.Clients[iterator].Gameplay.Hits.UnstableRate = baseUR / 1.5
		} else if strings.Contains(MenuData.Mods.PpMods, "HT") {
			TourneyData.Clients[iterator].Gameplay.Hits.UnstableRate = baseUR * 1.33
		} else {
			TourneyData.Clients[iterator].Gameplay.Hits.UnstableRate = baseUR
		}
	}
}
