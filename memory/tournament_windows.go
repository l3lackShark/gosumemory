//+build windows

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

	"github.com/l3lackShark/gosumemory/mem"
	"github.com/spf13/cast"
)

func resolveTourneyClients(procs []mem.Process) ([]mem.Process, error) {
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
		counter := 0
		for err != nil {
			if counter >= 30 {
				fmt.Println("Time's up! exiting tournament mode, failed after 30 attempts")
				return nil, errors.New("Tournament client timeout")
			}
			fmt.Println(fmt.Sprintf("[TOURNAMENT] %s, waiting for it...", err))
			time.Sleep(1 * time.Second)
			counter++
			client, err = mem.FindWindow(fmt.Sprintf("Tournament Client %d", i))
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

func getTourneyGameplayData(proc mem.Process, iterator int) {
	mem.Read(proc, &tourneyPatterns[iterator], &tourneyGameplayData[iterator])
	TourneyData.IPCClients[iterator].Gameplay.Combo.Current = tourneyGameplayData[iterator].Combo
	TourneyData.IPCClients[iterator].Gameplay.Combo.Max = tourneyGameplayData[iterator].MaxCombo
	TourneyData.IPCClients[iterator].Gameplay.GameMode = tourneyGameplayData[iterator].Mode
	TourneyData.IPCClients[iterator].Gameplay.Score = tourneyGameplayData[iterator].Score
	TourneyData.IPCClients[iterator].Gameplay.Hits.H100 = tourneyGameplayData[iterator].Hit100
	TourneyData.IPCClients[iterator].Gameplay.Hits.HKatu = tourneyGameplayData[iterator].HitKatu
	TourneyData.IPCClients[iterator].Gameplay.Hits.H200M = tourneyGameplayData[iterator].Hit200M
	TourneyData.IPCClients[iterator].Gameplay.Hits.H300 = tourneyGameplayData[iterator].Hit300
	TourneyData.IPCClients[iterator].Gameplay.Hits.HGeki = tourneyGameplayData[iterator].HitGeki
	TourneyData.IPCClients[iterator].Gameplay.Hits.H50 = tourneyGameplayData[iterator].Hit50
	TourneyData.IPCClients[iterator].Gameplay.Hits.H0 = tourneyGameplayData[iterator].HitMiss
	if TourneyData.IPCClients[iterator].Gameplay.Combo.Temp > TourneyData.IPCClients[iterator].Gameplay.Combo.Max {
		TourneyData.IPCClients[iterator].Gameplay.Combo.Temp = 0
	}
	if TourneyData.IPCClients[iterator].Gameplay.Combo.Current < TourneyData.IPCClients[iterator].Gameplay.Combo.Temp && TourneyData.IPCClients[iterator].Gameplay.Hits.H0Temp == TourneyData.IPCClients[iterator].Gameplay.Hits.H0 {
		TourneyData.IPCClients[iterator].Gameplay.Hits.HSB++
	}
	TourneyData.IPCClients[iterator].Gameplay.Hits.H0Temp = TourneyData.IPCClients[iterator].Gameplay.Hits.H0
	TourneyData.IPCClients[iterator].Gameplay.Combo.Temp = TourneyData.IPCClients[iterator].Gameplay.Combo.Current
	TourneyData.IPCClients[iterator].Gameplay.Accuracy = tourneyGameplayData[iterator].Accuracy
	TourneyData.IPCClients[iterator].Gameplay.Hp.Normal = tourneyGameplayData[iterator].PlayerHP
	TourneyData.IPCClients[iterator].Gameplay.Hp.Smooth = tourneyGameplayData[iterator].PlayerHPSmooth
	TourneyData.IPCClients[iterator].Gameplay.Name = tourneyGameplayData[iterator].PlayerName
	TourneyData.IPCClients[iterator].Gameplay.Mods.AppliedMods = int32(tourneyGameplayData[iterator].ModsXor1 ^ tourneyGameplayData[iterator].ModsXor2)
	if TourneyData.IPCClients[iterator].Gameplay.Mods.AppliedMods == 0 {
		TourneyData.IPCClients[iterator].Gameplay.Mods.PpMods = "NM"
	} else {
		TourneyData.IPCClients[iterator].Gameplay.Mods.PpMods = Mods(tourneyGameplayData[iterator].ModsXor1 ^ tourneyGameplayData[iterator].ModsXor2).String()
	}
	if TourneyData.IPCClients[iterator].Gameplay.Combo.Max > 0 {
		TourneyData.IPCClients[iterator].Gameplay.Hits.HitErrorArray = tourneyGameplayData[iterator].HitErrors
		baseUR, _ := calculateUR(TourneyData.IPCClients[iterator].Gameplay.Hits.HitErrorArray)
		if strings.Contains(TourneyData.IPCClients[iterator].Gameplay.Mods.PpMods, "DT") || strings.Contains(TourneyData.IPCClients[iterator].Gameplay.Mods.PpMods, "NC") {
			TourneyData.IPCClients[iterator].Gameplay.Hits.UnstableRate = baseUR / 1.5
		} else if strings.Contains(TourneyData.IPCClients[iterator].Gameplay.Mods.PpMods, "HT") {
			TourneyData.IPCClients[iterator].Gameplay.Hits.UnstableRate = baseUR * 1.33
		} else {
			TourneyData.IPCClients[iterator].Gameplay.Hits.UnstableRate = baseUR
		}
	}
}

func readTourneyIPCStruct(base int64) (int32, int32) {
	addresses := struct{ Base int64 }{base}
	var data struct {
		SpectatingID int32 `mem:"Base + 0x14"`
		Score        int32 `mem:"Base + 0x18"`
	}
	mem.Read(process, &addresses, &data)
	return data.SpectatingID, data.Score
}

func getTourneyIPC() error {
	err := mem.Read(process,
		&patterns,
		&tourneyManagerData)
	if err != nil {
		return err
	}
	TourneyData.Manager.BO = tourneyManagerData.BO
	TourneyData.Manager.IPCState = tourneyManagerData.IPCState
	TourneyData.Manager.Bools.ScoreVisible = cast.ToBool(int(tourneyManagerData.ScoreVisible))
	TourneyData.Manager.Bools.StarsVisible = cast.ToBool(int(tourneyManagerData.StarsVisible))
	TourneyData.Manager.Stars.Left = tourneyManagerData.LeftStars
	TourneyData.Manager.Stars.Right = tourneyManagerData.RightStars
	TourneyData.Manager.Name.Left = tourneyManagerData.TeamOneName
	TourneyData.Manager.Name.Right = tourneyManagerData.TeamTwoName
	TourneyData.Manager.Gameplay.Score.Left = tourneyManagerData.TeamOneScore
	TourneyData.Manager.Gameplay.Score.Right = tourneyManagerData.TeamTwoScore
	if TourneyData.Manager.IPCState != 3 && TourneyData.Manager.IPCState != 4 { //Playing, Ranking
		TourneyData.Manager.Gameplay = tmGameplay{}
		for i := range tourneyGameplayData {
			TourneyData.IPCClients[i].Gameplay = tourneyGameplay{}
		}
	}

	for i, proc := range tourneyProcs {
		err := mem.Read(proc,
			&tourneyPatterns[i].PreSongSelectAddresses,
			&tourneyMenuData[i].PreSongSelectData)
		if err != nil {
			DynamicAddresses.IsReady = false
			log.Println("It appears that we lost the precess, retrying", err)
			continue
		}
		if tourneyMenuData[i].PreSongSelectData.Status == 2 {
			getTourneyGameplayData(proc, i)
		}

	}

	for i, j := leaderStart, 0; j < int(tourneyManagerData.TotalAmOfClients); i, j = i+0x4, j+1 {
		slot, _ := mem.ReadUint32(process, int64(tourneyManagerData.IPCBaseAddr), int64(i))
		TourneyData.IPCClients[j].SpectatingID, TourneyData.IPCClients[j].Gameplay.Score = readTourneyIPCStruct(int64(slot))
	}
	return nil
}
