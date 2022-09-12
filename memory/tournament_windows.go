//go:build windows
// +build windows

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
	err := mem.Read(proc, &tourneyPatterns[iterator], &tourneyGameplayData[iterator])
	if err != nil && !strings.Contains(err.Error(), "LeaderBoard") && !strings.Contains(err.Error(), "KeyOverlay") { //TODO: fix this mem-side
		return //struct not initialized yet
	}
	TourneyData.IPCClients[iterator].Gameplay.Combo.Current = tourneyGameplayData[iterator].Combo
	TourneyData.IPCClients[iterator].Gameplay.Combo.Max = tourneyGameplayData[iterator].MaxCombo
	TourneyData.IPCClients[iterator].Gameplay.GameMode = tourneyGameplayData[iterator].Mode
	TourneyData.IPCClients[iterator].Gameplay.Score = tourneyGameplayData[iterator].ScoreV2
	TourneyData.IPCClients[iterator].Gameplay.Hits.H100 = tourneyGameplayData[iterator].Hit100
	TourneyData.IPCClients[iterator].Gameplay.Hits.HKatu = tourneyGameplayData[iterator].HitKatu
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
	addresses := struct {
		Base int64
	}{base}
	var data struct {
		SpectatingID int32 `mem:"Base + 0x14"`
		Score        int32 `mem:"Base + 0x18"`
	}

	mem.Read(process, &addresses, &data)

	return data.SpectatingID, data.Score
}

func readSpectatingUser(user int64, proc *mem.Process) (ipcSpec, error) {
	userAddr := struct {
		UserInfo int64
	}{user}
	var userData struct {
		Accuracy    float64 `mem:"[[UserInfo - 0x5]] + 0x4"`
		RankedScore int64   `mem:"[[UserInfo - 0x5]] + 0xC"`
		PlayCount   int32   `mem:"[[UserInfo - 0x5]] + 0x7C"`
		GlobalRank  int32   `mem:"[[UserInfo - 0x5]] + 0x84"`
		PP          int32   `mem:"[[UserInfo - 0x5]] + 0x9C"`
		Name        string  `mem:"[[[UserInfo - 0x5]] + 0x30]"`
		Country     string  `mem:"[[[UserInfo - 0x5]] + 0x2C]"`
		UserID      int32   `mem:"[[UserInfo - 0x5]] + 0x70"`
	}
	err := mem.Read(*proc, &userAddr, &userData)
	if err != nil {
		return ipcSpec{}, errors.New("[TOURNAMENT] Could not read userData")
	}

	return ipcSpec{
		Accuracy:    userData.Accuracy,
		GlobalPP:    userData.PP,
		GlobalRank:  userData.GlobalRank,
		Name:        userData.Name,
		Country:     userData.Country,
		ID:          userData.UserID,
		PlayCount:   userData.PlayCount,
		RankedScore: userData.RankedScore,
	}, nil
}

func getTourneyIPC() error {
	err := mem.Read(process,
		&patterns,
		&tourneyManagerData)
	if err != nil {
		return err
	}
	chatClass := patterns.ChatArea - 0x44
	TourneyData.Manager.Chat, err = readChatData(&chatClass)
	if err != nil {
		log.Println(err)
		DynamicAddresses.IsReady = false
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

	for i, j := leaderStart, 0; j < len(tourneyProcs); i, j = i+0x4, j+1 {
		slot, _ := mem.ReadUint32(process, int64(tourneyManagerData.IPCBaseAddr), int64(i))

		TourneyData.IPCClients[j].ID, TourneyData.IPCClients[j].Gameplay.Score = readTourneyIPCStruct(int64(slot))

	}

	for i, proc := range tourneyProcs {
		err := mem.Read(proc,
			&tourneyPatterns[i].PreSongSelectAddresses,
			&tourneyMenuData[i].PreSongSelectData)
		if err != nil {
			DynamicAddresses.IsReady = false
			log.Println("It appears that we lost the process, retrying", err)
			continue
		}
		if i == 0 {
			if tourneyMenuData[0].PreSongSelectData.Status == 0 {
				mem.Read(proc, &tourneyPatterns[0], &mainMenuData)
				MenuData.MainMenuValues.BassDensity = calculateBassDensity(mainMenuData.AudioVelocityBase, &proc)
			}
		}
		switch tourneyMenuData[i].PreSongSelectData.Status {
		case 2:
			getTourneyGameplayData(proc, i)
		}
		if TourneyData.IPCClients[i].ID > 0 {
			totalClients := len(TourneyData.IPCClients)
			if i < totalClients/2 {
				TourneyData.IPCClients[i].Team = "left"
			} else {
				TourneyData.IPCClients[i].Team = "right"
			}
			TourneyData.IPCClients[i].Spectating, _ = readSpectatingUser(tourneySpecificPatterns[i].UserInfo, &proc)
		} else {
			TourneyData.IPCClients[i] = ipcClient{}
		}
	}
	return nil
}

func readChatData(base *int64) (result []tourneyMessage, err error) {
	addresses := struct{ Base int64 }{*base}
	var data struct {
		Tabs uint32 `mem:"[Base + 0x1C] + 0x4"`
	}
	err = mem.Read(process, &addresses, &data)
	if err != nil {
		return nil, errors.New("[TOURNEY CHAT] Failed reading the main struct")
	}
	length, err := mem.ReadInt32(process, int64(data.Tabs), 4)
	for i, j := leaderStart, 0; j < int(length); i, j = i+0x4, j+1 {
		slot, _ := mem.ReadUint32(process, int64(data.Tabs), int64(i))
		if slot == 0 {
			continue
		}
		addrs := struct{ Base int64 }{int64(slot)}
		var chatData struct {
			ChatTag      string `mem:"[[Base + 0xC] + 0x4]"`
			MessagesAddr uint32 `mem:"[[Base + 0xC] + 0x10] + 0x4"`
		}
		err := mem.Read(process, &addrs, &chatData)
		if err != nil || chatData.ChatTag != "#multiplayer" {
			continue
		}
		msgLength, err := mem.ReadInt32(process, int64(chatData.MessagesAddr), 4)
		var messages []tourneyMessage
		for n, k := leaderStart, 0; k < int(msgLength); n, k = n+0x4, k+1 {
			msgSlot, err := mem.ReadUint32(process, int64(chatData.MessagesAddr), int64(n))
			if err != nil {
				return nil, errors.New("[TOURNEY CHAT] Internal error")
			}
			msgAddrs := struct{ Base int64 }{int64(msgSlot)}
			var chatContent struct {
				TimeName string `mem:"[Base + 0x8]"`
				Content  string `mem:"[Base + 0x4]"`
			}
			err = mem.Read(process, &msgAddrs, &chatContent)
			if chatContent.Content == "" || strings.HasPrefix(chatContent.Content, "!mp") {
				continue
			}
			spl := strings.SplitAfterN(chatContent.TimeName, " ", 2)
			if len(spl) < 2 {
				return nil, errors.New("[TOURNEY CHAT] Internal error, could not split")
			}
			var msg tourneyMessage
			msg.Time = strings.TrimSpace(spl[0])
			msg.Name = strings.TrimSuffix(spl[1], ":")
			msg.MessageBody = chatContent.Content
			for _, client := range TourneyData.IPCClients {
				if client.Spectating.Name == msg.Name {
					msg.Team = client.Team
				}
			}
			if msg.Team == "" {
				if msg.Name == "BanchoBot" {
					msg.Team = "bot"
				} else {
					msg.Team = "unknown"
				}
			}

			messages = append(messages, msg)

		}
		if len(messages) > 0 {
			return messages, nil
		}
		return nil, nil
	}
	return nil, nil
}
