package memory

import (
	"errors"
	"fmt"
	"log"
	"math"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/mem"
)

func modsResolver(xor uint32) string {
	return Mods(xor).String()
}

//UpdateTime Intervall between value updates
var UpdateTime int

//UnderWine?
var UnderWine bool

var isTournamentClient bool
var tourneyPID int

//var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
var leaderStart int32
var hasLeaderboard = false

//SongsFolderPath is full path to osu! Songs. Gets set automatically on Windows (through memory)
var SongsFolderPath string

var process, procerr = mem.FindProcess(osuProcessRegex)

var tempBeatmapString string = ""
var tempRetries int32

//Init the whole thing and get osu! memory values to start working with it.
func Init() {
	if UnderWine == true || runtime.GOOS != "windows" {
		leaderStart = 0xC
	} else {
		leaderStart = 0x8
	}
	var err error
	for {
		process, procerr = mem.FindProcess(osuProcessRegex)
		if procerr != nil {
			DynamicAddresses.IsReady = false
			for procerr != nil {
				process, procerr = mem.FindProcess(osuProcessRegex)
				log.Println("It seems that we lost the process, retrying!")
				time.Sleep(1 * time.Second)
			}
			isTournamentClient, err = initBase()
			for err != nil {
				isTournamentClient, err = initBase()
				time.Sleep(1 * time.Second)
			}
		}
		if DynamicAddresses.IsReady == false {
			isTournamentClient, err = initBase()
			for err != nil {
				isTournamentClient, err = initBase()
				if err != nil {
					log.Println("Failure mid getting offsets, retrying")
				}
				time.Sleep(1 * time.Second)

			}
		} else {
			if isTournamentClient {
				for i := 0; i < len(process); i++ {
					if i != tourneyManagerID {
						readProcessMemory(i)
					}
				}
				fmt.Println("We are in tournament mode!")

			} else {
				readProcessMemory(0)
			}
		}
		time.Sleep(time.Duration(UpdateTime) * time.Millisecond)
	}

}

func readProcessMemory(proc int) error {

	err := mem.Read(process[proc],
		&patterns[proc].PreSongSelectAddresses,
		&menuData[proc].PreSongSelectData)
	if err != nil {
		DynamicAddresses.IsReady = false
		log.Println("It appears that we lost the precess, retrying", err)
		return err
	}
	MenuData.OsuStatus = menuData[proc].Status

	mem.Read(process[proc], &patterns[proc], &alwaysData[proc])

	MenuData.ChatChecker = alwaysData[proc].ChatStatus
	MenuData.Bm.Time.PlayTime = alwaysData[proc].PlayTime
	MenuData.SkinFolder = alwaysData[proc].SkinFolder
	switch menuData[proc].Status {
	case 2:
		if MenuData.Bm.Time.PlayTime < 150 || menuData[proc].Path == "" { //To catch up with the F2-->Enter
			err := bmUpdateData(proc)
			if err != nil {
				pp.Println(err)
			}
		}
		if gameplayData[proc].Retries > tempRetries {
			tempRetries = gameplayData[proc].Retries
			GameplayData = GameplayValues{}
			gameplayData[proc] = gameplayD{}

		}
		getGamplayData(proc)
	case 1:
		err = bmUpdateData(proc)
		if err != nil {
			pp.Println(err)
		}
	case 7:
	default:
		tempRetries = -1
		GameplayData = GameplayValues{}
		gameplayData[proc] = gameplayD{}
		hasLeaderboard = false
		err = bmUpdateData(proc)
		if err != nil {
			pp.Println(err)
		}
	}
	return nil
}

func bmUpdateData(proc int) error {
	mem.Read(process[proc], &patterns[proc], &menuData[proc])

	bmString := menuData[proc].Path
	if strings.HasSuffix(bmString, ".osu") && tempBeatmapString != bmString { //On map change
		for i := 0; i < 50; i++ {
			if menuData[proc].BackgroundFilename != "" {
				break
			}
			time.Sleep(25 * time.Millisecond)
			mem.Read(process[proc], &patterns[proc], &menuData[proc])
		}
		tempBeatmapString = bmString
		MenuData.Bm.BeatmapID = menuData[proc].MapID
		MenuData.Bm.BeatmapSetID = menuData[proc].SetID
		MenuData.Bm.Path = path{
			AudioPath:            menuData[proc].AudioFilename,
			BGPath:               menuData[proc].BackgroundFilename,
			BeatmapOsuFileString: menuData[proc].Path,
			BeatmapFolderString:  menuData[proc].Folder,
			FullMP3Path:          filepath.Join(SongsFolderPath, menuData[proc].Folder, menuData[proc].AudioFilename),
			FullDotOsu:           filepath.Join(SongsFolderPath, menuData[proc].Folder, bmString),
			InnerBGPath:          filepath.Join(menuData[proc].Folder, menuData[proc].BackgroundFilename),
		}
		MenuData.Bm.Stats.MemoryAR = menuData[proc].AR
		MenuData.Bm.Stats.MemoryCS = menuData[proc].CS
		MenuData.Bm.Stats.MemoryHP = menuData[proc].HP
		MenuData.Bm.Stats.MemoryOD = menuData[proc].OD
		MenuData.Bm.Metadata.Artist = menuData[proc].Artist
		MenuData.Bm.Metadata.Title = menuData[proc].Title
		MenuData.Bm.Metadata.Mapper = menuData[proc].Creator
		MenuData.Bm.Metadata.Version = menuData[proc].Difficulty
		MenuData.GameMode = menuData[proc].MenuGameMode
		MenuData.Bm.RandkedStatus = menuData[proc].RankedStatus
		MenuData.Bm.BeatmapMD5 = menuData[proc].MD5
	}
	if alwaysData[proc].MenuMods == 0 {
		MenuData.Mods.PpMods = "NM"
		MenuData.Mods.AppliedMods = int32(alwaysData[proc].MenuMods)
	} else {
		MenuData.Mods.AppliedMods = int32(alwaysData[proc].MenuMods)
		MenuData.Mods.PpMods = Mods(alwaysData[proc].MenuMods).String()
	}

	return nil
}
func getGamplayData(proc int) {
	mem.Read(process[proc], &patterns[proc], &gameplayData[proc])
	//GameplayData.BitwiseKeypress = gameplayData[proc].BitwiseKeypress
	GameplayData.Combo.Current = gameplayData[proc].Combo
	GameplayData.Combo.Max = gameplayData[proc].MaxCombo
	fmt.Println(gameplayData[proc].Score)
	GameplayData.GameMode = gameplayData[proc].Mode
	GameplayData.Score = gameplayData[proc].Score
	GameplayData.Hits.H100 = gameplayData[proc].Hit100
	GameplayData.Hits.HKatu = gameplayData[proc].HitKatu
	GameplayData.Hits.H200M = gameplayData[proc].Hit200M
	GameplayData.Hits.H300 = gameplayData[proc].Hit300
	GameplayData.Hits.HGeki = gameplayData[proc].HitGeki
	GameplayData.Hits.H50 = gameplayData[proc].Hit50
	GameplayData.Hits.H0 = gameplayData[proc].HitMiss
	if GameplayData.Combo.Temp > GameplayData.Combo.Max {
		GameplayData.Combo.Temp = 0
	}
	if GameplayData.Combo.Current < GameplayData.Combo.Temp && GameplayData.Hits.H0Temp == GameplayData.Hits.H0 {
		GameplayData.Hits.HSB++
	}
	GameplayData.Hits.H0Temp = GameplayData.Hits.H0
	GameplayData.Combo.Temp = GameplayData.Combo.Current
	MenuData.Mods.AppliedMods = int32(gameplayData[proc].ModsXor1 ^ gameplayData[proc].ModsXor1)
	GameplayData.Accuracy = gameplayData[proc].Accuracy
	GameplayData.Hp.Normal = gameplayData[proc].PlayerHP
	GameplayData.Hp.Smooth = gameplayData[proc].PlayerHPSmooth
	GameplayData.Name = gameplayData[proc].PlayerName
	fmt.Println(gameplayData[proc].PlayerName)
	MenuData.Mods.AppliedMods = int32(gameplayData[proc].ModsXor1 ^ gameplayData[proc].ModsXor2)
	if MenuData.Mods.AppliedMods == 0 {
		MenuData.Mods.PpMods = "NM"
	} else {
		MenuData.Mods.PpMods = Mods(gameplayData[proc].ModsXor1 ^ gameplayData[proc].ModsXor2).String()
	}
	if GameplayData.Combo.Max > 0 {
		GameplayData.Hits.HitErrorArray = gameplayData[proc].HitErrors
		GameplayData.Hits.UnstableRate, _ = calculateUR(GameplayData.Hits.HitErrorArray)
	}
	getLeaderboard(proc)
}

func getLeaderboard(proc int) {
	var board leaderboard
	if gameplayData[proc].LeaderBoard == 0 {
		board.DoesLeaderBoardExists = false
		GameplayData.Leaderboard = board
		return
	}
	board.DoesLeaderBoardExists = true
	ourPlayerStruct, _ := mem.ReadUint32(process[proc], int64(gameplayData[proc].LeaderBoard)+0x10, 0)
	board.OurPlayer = readLeaderPlayerStruct(proc, int64(ourPlayerStruct))
	board.OurPlayer.Mods = MenuData.Mods.PpMods //ourplayer mods is sometimes delayed so better default to PlayContainer Here
	playersArray, _ := mem.ReadUint32(process[proc], int64(gameplayData[proc].LeaderBoard)+0x4)
	amOfSlots, _ := mem.ReadInt32(process[proc], int64(playersArray+0xC))
	if amOfSlots < 1 || amOfSlots > 64 {
		return
	}
	items, _ := mem.ReadInt32(process[proc], int64(playersArray+0x4))
	board.Slots = make([]leaderPlayer, amOfSlots)
	for i, j := 0x8, 0; j < int(amOfSlots); i, j = i+0x4, j+1 {
		slot, _ := mem.ReadUint32(process[proc], int64(items), int64(i))
		board.Slots[j] = readLeaderPlayerStruct(proc, int64(slot))
	}
	GameplayData.Leaderboard = board
}

func readLeaderPlayerStruct(proc int, base int64) leaderPlayer {
	addresses := struct{ Base int64 }{base}
	var player struct {
		Name      string `mem:"[Base + 0x8]"`
		Score     int32  `mem:"Base + 0x30"`
		Combo     int16  `mem:"[Base + 0x20] + 0x94"`
		MaxCombo  int16  `mem:"[Base + 0x20] + 0x68"`
		ModsXor1  uint32 `mem:"[[Base + 0x20] + 0x1C] + 0x8"`
		ModsXor2  uint32 `mem:"[[Base + 0x20] + 0x1C] + 0xC"`
		H300      int16  `mem:"[Base + 0x20] + 0x8A"`
		H100      int16  `mem:"[Base + 0x20] + 0x88"`
		H50       int16  `mem:"[Base + 0x20] + 0x8C"`
		H0        int16  `mem:"[Base + 0x20] + 0x92"`
		Team      int32  `mem:"Base + 0x40"`
		Position  int32  `mem:"Base + 0x2C"`
		IsPassing int8   `mem:"Base + 0x4B"`
	}
	mem.Read(process[proc], &addresses, &player)
	mods := modsResolver(player.ModsXor1 ^ player.ModsXor2)
	if mods == "" {
		mods = "NM"
	}
	return leaderPlayer{
		Name:      player.Name,
		Score:     player.Score,
		Combo:     player.Combo,
		MaxCombo:  player.MaxCombo,
		Mods:      mods,
		H300:      player.H300,
		H100:      player.H100,
		H50:       player.H50,
		H0:        player.H0,
		Team:      player.Team,
		Position:  player.Position,
		IsPassing: player.IsPassing,
	}
}

func calculateUR(HitErrorArray []int32) (float64, error) {
	if len(HitErrorArray) < 1 {
		return 0, errors.New("Empty hit error array")
	}
	var totalAll float32 //double
	for _, hit := range HitErrorArray {
		totalAll += float32(hit)
	}
	var average float32 = totalAll / float32(len(HitErrorArray))
	var variance float64 = 0
	for _, hit := range HitErrorArray {
		variance += math.Pow(float64(hit)-float64(average), 2)
	}
	variance = variance / float64(len(HitErrorArray))
	return math.Sqrt(variance) * 10, nil

}
