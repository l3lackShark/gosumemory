package memory

import (
	"errors"
	"log"
	"math"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cast"

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
var isTournamentMode bool
var tourneyProcs []mem.Process
var tourneyErr error

//var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
var leaderStart int32
var hasLeaderboard = false

//SongsFolderPath is full path to osu! Songs. Gets set automatically on Windows (through memory)
var SongsFolderPath string

var allProcs []mem.Process
var process mem.Process
var procerr error
var tempRetries int32

//Init the whole thing and get osu! memory values to start working with it.
func Init() {
	if UnderWine == true || runtime.GOOS != "windows" {
		leaderStart = 0xC
	} else {
		leaderStart = 0x8
	}

	allProcs, procerr = mem.FindProcess(osuProcessRegex)
	for {
		start := time.Now()
		if procerr != nil {
			DynamicAddresses.IsReady = false
			for procerr != nil {
				allProcs, procerr = mem.FindProcess(osuProcessRegex)
				log.Println("It seems that we lost the process, retrying!")
				time.Sleep(1 * time.Second)
			}
			err := initBase()
			for err != nil {
				err = initBase()
				time.Sleep(1 * time.Second)
			}
		}
		if DynamicAddresses.IsReady == false {
			err := initBase()
			for err != nil {
				err = initBase()
				if err != nil {
					log.Println("Failure mid getting offsets, retrying", err)
				}
				time.Sleep(1 * time.Second)

			}
		} else {
			err := mem.Read(process,
				&patterns.PreSongSelectAddresses,
				&menuData.PreSongSelectData)
			if err != nil {
				DynamicAddresses.IsReady = false
				log.Println("It appears that we lost the precess, retrying", err)
				continue
			}
			MenuData.OsuStatus = menuData.Status

			mem.Read(process, &patterns, &alwaysData)
			MenuData.ChatChecker = alwaysData.ChatStatus
			MenuData.Bm.Time.PlayTime = alwaysData.PlayTime
			MenuData.SkinFolder = alwaysData.SkinFolder
			SettingsData.ShowInterface = cast.ToBool(int(alwaysData.ShowInterface))
			switch menuData.Status {
			case 2:
				if MenuData.Bm.Time.PlayTime < 150 || menuData.Path == "" { //To catch up with the F2-->Enter
					err := bmUpdateData()
					if err != nil {
						pp.Println(err)
					}
				}
				if gameplayData.Retries > tempRetries {
					tempRetries = gameplayData.Retries
					GameplayData = GameplayValues{}
					gameplayData = gameplayD{}

				}
				getGamplayData()
			case 1:
				err = bmUpdateData()
				if err != nil {
					pp.Println(err)
				}
			case 7:
			default:
				tempRetries = -1
				GameplayData = GameplayValues{}
				gameplayData = gameplayD{}
				hasLeaderboard = false
				err = bmUpdateData()
				if err != nil {
					pp.Println(err)
				}

			}
		}
		if isTournamentMode {
			err := mem.Read(allProcs[0],
				&patterns,
				&tourneyManagerData)
			if err != nil {
				DynamicAddresses.IsReady = false
				log.Println("It appears that we lost the precess, retrying", err)
				continue
			}
			TourneyData.Manager.BO = tourneyManagerData.BO
			TourneyData.Manager.IPCState = tourneyManagerData.IPCState
			TourneyData.Manager.ScoreVisible = cast.ToBool(int(tourneyManagerData.ScoreVisible))
			TourneyData.Manager.StarsVisible = cast.ToBool(int(tourneyManagerData.StarsVisible))
			TourneyData.Manager.StarsLeft = tourneyManagerData.LeftStars
			TourneyData.Manager.StarsRight = tourneyManagerData.RightStars
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
					//fmt.Println(fmt.Sprintf("Client #%d is in play mode!", i))
					getTourneyGameplayData(proc, i)
				}

			}

		}
		elapsed := time.Since(start)
		log.Printf("Cycle took %s", elapsed)
		time.Sleep(time.Duration(UpdateTime-int(elapsed.Milliseconds())) * time.Millisecond)

	}

}

var tempBeatmapString string = ""

func bmUpdateData() error {
	mem.Read(process, &patterns, &menuData)

	bmString := menuData.Path
	if strings.HasSuffix(bmString, ".osu") && tempBeatmapString != bmString { //On map change
		for i := 0; i < 50; i++ {
			if menuData.BackgroundFilename != "" {
				break
			}
			time.Sleep(25 * time.Millisecond)
			mem.Read(process, &patterns, &menuData)
		}
		tempBeatmapString = bmString
		MenuData.Bm.BeatmapID = menuData.MapID
		MenuData.Bm.BeatmapSetID = menuData.SetID
		MenuData.Bm.Path = path{
			AudioPath:            menuData.AudioFilename,
			BGPath:               menuData.BackgroundFilename,
			BeatmapOsuFileString: menuData.Path,
			BeatmapFolderString:  menuData.Folder,
			FullMP3Path:          filepath.Join(SongsFolderPath, menuData.Folder, menuData.AudioFilename),
			FullDotOsu:           filepath.Join(SongsFolderPath, menuData.Folder, bmString),
			InnerBGPath:          filepath.Join(menuData.Folder, menuData.BackgroundFilename),
		}
		MenuData.Bm.Stats.MemoryAR = menuData.AR
		MenuData.Bm.Stats.MemoryCS = menuData.CS
		MenuData.Bm.Stats.MemoryHP = menuData.HP
		MenuData.Bm.Stats.MemoryOD = menuData.OD
		MenuData.Bm.Metadata.Artist = menuData.Artist
		MenuData.Bm.Metadata.Title = menuData.Title
		MenuData.Bm.Metadata.Mapper = menuData.Creator
		MenuData.Bm.Metadata.Version = menuData.Difficulty
		MenuData.GameMode = menuData.MenuGameMode
		MenuData.Bm.RandkedStatus = menuData.RankedStatus
		MenuData.Bm.BeatmapMD5 = menuData.MD5
	}
	if alwaysData.MenuMods == 0 {
		MenuData.Mods.PpMods = "NM"
		MenuData.Mods.AppliedMods = int32(alwaysData.MenuMods)
	} else {
		MenuData.Mods.AppliedMods = int32(alwaysData.MenuMods)
		MenuData.Mods.PpMods = Mods(alwaysData.MenuMods).String()
	}

	return nil
}
func getGamplayData() {
	mem.Read(process, &patterns, &gameplayData)
	//GameplayData.BitwiseKeypress = gameplayData.BitwiseKeypress
	GameplayData.Combo.Current = gameplayData.Combo
	GameplayData.Combo.Max = gameplayData.MaxCombo
	GameplayData.GameMode = gameplayData.Mode
	GameplayData.Score = gameplayData.Score
	GameplayData.Hits.H100 = gameplayData.Hit100
	GameplayData.Hits.HKatu = gameplayData.HitKatu
	GameplayData.Hits.H200M = gameplayData.Hit200M
	GameplayData.Hits.H300 = gameplayData.Hit300
	GameplayData.Hits.HGeki = gameplayData.HitGeki
	GameplayData.Hits.H50 = gameplayData.Hit50
	GameplayData.Hits.H0 = gameplayData.HitMiss
	if GameplayData.Combo.Temp > GameplayData.Combo.Max {
		GameplayData.Combo.Temp = 0
	}
	if GameplayData.Combo.Current < GameplayData.Combo.Temp && GameplayData.Hits.H0Temp == GameplayData.Hits.H0 {
		GameplayData.Hits.HSB++
	}
	GameplayData.Hits.H0Temp = GameplayData.Hits.H0
	GameplayData.Combo.Temp = GameplayData.Combo.Current
	MenuData.Mods.AppliedMods = int32(gameplayData.ModsXor1 ^ gameplayData.ModsXor1)
	GameplayData.Accuracy = gameplayData.Accuracy
	GameplayData.Hp.Normal = gameplayData.PlayerHP
	GameplayData.Hp.Smooth = gameplayData.PlayerHPSmooth
	GameplayData.Name = gameplayData.PlayerName
	MenuData.Mods.AppliedMods = int32(gameplayData.ModsXor1 ^ gameplayData.ModsXor2)
	if MenuData.Mods.AppliedMods == 0 {
		MenuData.Mods.PpMods = "NM"
	} else {
		MenuData.Mods.PpMods = Mods(gameplayData.ModsXor1 ^ gameplayData.ModsXor2).String()
	}
	if GameplayData.Combo.Max > 0 {
		GameplayData.Hits.HitErrorArray = gameplayData.HitErrors
		baseUR, _ := calculateUR(GameplayData.Hits.HitErrorArray)
		if strings.Contains(MenuData.Mods.PpMods, "DT") || strings.Contains(MenuData.Mods.PpMods, "NC") {
			GameplayData.Hits.UnstableRate = baseUR / 1.5
		} else if strings.Contains(MenuData.Mods.PpMods, "HT") {
			GameplayData.Hits.UnstableRate = baseUR * 1.33
		} else {
			GameplayData.Hits.UnstableRate = baseUR
		}
	}
	getLeaderboard()
}

func getLeaderboard() {
	var board leaderboard
	if gameplayData.LeaderBoard == 0 {
		board.DoesLeaderBoardExists = false
		GameplayData.Leaderboard = board
		return
	}
	board.DoesLeaderBoardExists = true
	ourPlayerStruct, _ := mem.ReadUint32(process, int64(gameplayData.LeaderBoard)+0x10, 0)
	board.OurPlayer, board.IsLeaderBoardVisible = readLeaderPlayerStruct(int64(ourPlayerStruct))
	board.OurPlayer.Mods = MenuData.Mods.PpMods //ourplayer mods is sometimes delayed so better default to PlayContainer Here
	playersArray, _ := mem.ReadUint32(process, int64(gameplayData.LeaderBoard)+0x4)
	amOfSlots, _ := mem.ReadInt32(process, int64(playersArray+0xC))
	if amOfSlots < 1 || amOfSlots > 64 {
		return
	}
	items, _ := mem.ReadInt32(process, int64(playersArray+0x4))
	board.Slots = make([]leaderPlayer, amOfSlots)
	for i, j := 0x8, 0; j < int(amOfSlots); i, j = i+0x4, j+1 {
		slot, _ := mem.ReadUint32(process, int64(items), int64(i))
		board.Slots[j], _ = readLeaderPlayerStruct(int64(slot))
	}
	GameplayData.Leaderboard = board
}

func readLeaderPlayerStruct(base int64) (leaderPlayer, bool) {
	addresses := struct{ Base int64 }{base}
	var player struct {
		Name                 string `mem:"[Base + 0x8]"`
		Score                int32  `mem:"Base + 0x30"`
		Combo                int16  `mem:"[Base + 0x20] + 0x94"`
		MaxCombo             int16  `mem:"[Base + 0x20] + 0x68"`
		ModsXor1             uint32 `mem:"[[Base + 0x20] + 0x1C] + 0x8"`
		ModsXor2             uint32 `mem:"[[Base + 0x20] + 0x1C] + 0xC"`
		H300                 int16  `mem:"[Base + 0x20] + 0x8A"`
		H100                 int16  `mem:"[Base + 0x20] + 0x88"`
		H50                  int16  `mem:"[Base + 0x20] + 0x8C"`
		H0                   int16  `mem:"[Base + 0x20] + 0x92"`
		Team                 int32  `mem:"Base + 0x40"`
		Position             int32  `mem:"Base + 0x2C"`
		IsPassing            int8   `mem:"Base + 0x4B"`
		IsLeaderboardVisible int8   `mem:"[Base + 0x24] + 0x20"`
	}
	mem.Read(process, &addresses, &player)
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
	}, cast.ToBool(int(player.IsLeaderboardVisible))
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
