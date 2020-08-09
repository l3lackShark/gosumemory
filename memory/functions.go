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

//var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
var leaderStart int32
var hasLeaderboard = false

//SongsFolderPath is full path to osu! Songs. Gets set automatically on Windows (through memory)
var SongsFolderPath string

var process, procerr = mem.FindProcess(osuProcessRegex)

//Init the whole thing and get osu! memory values to start working with it.
func Init() {
	if UnderWine == true || runtime.GOOS != "windows" {
		leaderStart = 0xC
	} else {
		leaderStart = 0x8
	}

	for {
		process, procerr = mem.FindProcess(osuProcessRegex)
		if procerr != nil {
			DynamicAddresses.IsReady = false
			for procerr != nil {
				process, procerr = mem.FindProcess(osuProcessRegex)
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
					log.Println("Failure mid getting offsets, retrying")
				}
				time.Sleep(1 * time.Second)

			}
		} else {
			err := mem.Read(process,
				&patterns.PreSongSelectAddresses,
				&menuData.PreSongSelectData)
			if err != nil {
				DynamicAddresses.IsReady = false
				log.Println("It appears that we lost the precess, retrying")
				initBase()
			}
			MenuData.OsuStatus = menuData.Status

			mem.Read(process, &patterns, &alwaysData)

			MenuData.ChatChecker = alwaysData.ChatStatus
			MenuData.Bm.Time.PlayTime = alwaysData.PlayTime
			MenuData.SkinFolder = alwaysData.SkinFolder
			switch menuData.Status {

			case 2:
				if MenuData.Bm.Time.PlayTime < 150 { //To catch up with the F2-->Enter
					err := bmUpdateData()
					if err != nil {
						pp.Println(err)
					}
				}
				getGamplayData()
			case 1:
				err = bmUpdateData()
				if err != nil {
					pp.Println(err)
				}
			case 7:
			default:
				GameplayData = GameplayValues{}
				hasLeaderboard = false
				err = bmUpdateData()
				if err != nil {
					pp.Println(err)
				}

			}
			time.Sleep(time.Duration(UpdateTime) * time.Millisecond)
		}
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
	GameplayData.IsFailed = gameplayData.IsFailed
	GameplayData.FailTime = gameplayData.ReplayFailTime
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
	MenuData.Mods.AppliedMods = int32(gameplayData.ModsXor1 ^ gameplayData.ModsXor1)
	GameplayData.Accuracy = gameplayData.Accuracy
	GameplayData.Hp.Normal = gameplayData.PlayerHP
	GameplayData.Hp.Smooth = gameplayData.PlayerHPSmooth
	GameplayData.Name = gameplayData.PlayerName
	MenuData.Mods.AppliedMods = int32(gameplayData.ModsXor1 ^ gameplayData.ModsXor2)
	MenuData.Mods.PpMods = Mods(gameplayData.ModsXor1 ^ gameplayData.ModsXor2).String()
	if GameplayData.Combo.Max > 0 {
		GameplayData.Hits.HitErrorArray = gameplayData.HitErrors
		GameplayData.Hits.UnstableRate, _ = calculateUR(GameplayData.Hits.HitErrorArray)
	}
	getLeaderboard()
	var err error
	GameplayData.Replay, err = readOSREntries()
	if err != nil {
		fmt.Println(err)
	}
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
	board.OurPlayer = readLeaderPlayerStruct(int64(ourPlayerStruct))
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
		board.Slots[j] = readLeaderPlayerStruct(int64(slot))
	}
	GameplayData.Leaderboard = board
}

func readLeaderPlayerStruct(base int64) leaderPlayer {
	name, _ := mem.ReadString(process, base, 0x8, 0)
	score, _ := mem.ReadInt32(process, base+0x30, 0)
	combo, _ := mem.ReadInt16(process, base+0x20, 0, 0x94)
	maxCombo, _ := mem.ReadInt16(process, base+0x20, 0, 0x68)
	modsXor1, _ := mem.ReadUint32(process, base+0x20, 0, 0x1C, 0x8)
	modsXor2, _ := mem.ReadUint32(process, base+0x20, 0, 0x1C, 0xC)
	var mods string
	if modsXor1^modsXor2 != 0 {
		mods = modsResolver(modsXor1 ^ modsXor2)
	} else {
		mods = "NM"
	}
	h300, _ := mem.ReadInt16(process, base+0x20, 0, 0x8A)
	h100, _ := mem.ReadInt16(process, base+0x20, 0, 0x88)
	h50, _ := mem.ReadInt16(process, base+0x20, 0, 0x8C)
	h0, _ := mem.ReadInt16(process, base+0x20, 0, 0x92)
	team, _ := mem.ReadInt32(process, base+0x40, 0)
	position, _ := mem.ReadInt32(process, base+0x2C, 0)
	isPassing, _ := mem.ReadInt8(process, base+0x4B, 0)
	player := leaderPlayer{
		name,
		score,
		combo,
		maxCombo,
		mods,
		h300,
		h100,
		h50,
		h0,
		team,
		position,
		isPassing,
	}
	return player
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

type ReplayArray struct {
	BmHash  string
	Replays []OSREntry
}

type OSREntry struct {
	X                float32
	Y                float32
	WasButtonPressed int8 //Bitwise combination of keys/mouse buttons (0 - no keypress)
	Time             int32
}

func readOSREntries() (ReplayArray, error) {

	items, err := mem.ReadInt32(process, int64(gameplayData.ReplayDataBase)+0xC)
	//	var osr ReplayArray
	if err != nil {
		return ReplayArray{}, err
	}
	if items > 100000 || items < 1 {
		return ReplayArray{}, errors.New("invalid struct or empty array")
	}
	arraysBase, err := mem.ReadInt32(process, int64(gameplayData.ReplayDataBase)+0x4, 0)
	if err != nil {
		return ReplayArray{}, err
	}
	var osr ReplayArray

	osr.Replays = make([]OSREntry, items)
	for i, j := 0x8, 0; j < int(items); i, j = i+0x4, j+1 {
		ourArray, err := mem.ReadUint32(process, int64(arraysBase)+int64(i), 0)
		if err != nil {
			return ReplayArray{}, err
		}
		x, _ := mem.ReadFloat32(process, int64(ourArray)+0x4, 0)
		y, _ := mem.ReadFloat32(process, int64(ourArray)+0x8, 0)
		wasButtonPressed, err := mem.ReadInt8(process, int64(ourArray)+0xC, 0)
		time, _ := mem.ReadInt32(process, int64(ourArray)+0x10, 0)

		osr.Replays[j] = OSREntry{
			x,
			y,
			wasButtonPressed,
			time,
		}
	}
	osr.BmHash = MenuData.Bm.BeatmapMD5
	return osr, nil
}
