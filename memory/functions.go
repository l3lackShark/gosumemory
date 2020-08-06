package memory

import (
	"log"
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
	GameplayData.Combo.Current = gameplayData.Combo
	GameplayData.Combo.Max = gameplayData.MaxCombo
	GameplayData.GameMode = gameplayData.Mode
	GameplayData.Score = gameplayData.Score
	GameplayData.Hits.H100 = gameplayData.Hit100
	GameplayData.Hits.HKatu = gameplayData.HitKatu
	GameplayData.Hits.H200M = gameplayData.Hit200M
	GameplayData.Hits.H300 = gameplayData.Hit300
	GameplayData.Hits.H350 = gameplayData.Hit350
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
	}

}
