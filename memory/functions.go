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

	"github.com/spf13/cast"

	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/mem"
)

func modsResolver(xor uint32) string {
	return Mods(xor).String()
}

// UpdateTime Intervall between value updates
var UpdateTime int

// UnderWine?
var UnderWine bool

// MemCycle test
var MemCycle bool
var isTournamentMode bool
var tourneyProcs []mem.Process
var tourneyErr error

// Were we in the Result Screen?
var dirtyResults bool = false

// var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
var leaderStart int32

// SongsFolderPath is full path to osu! Songs. Gets set automatically on Windows (through memory)
var SongsFolderPath string

var allProcs []mem.Process
var process mem.Process
var procerr error
var tempRetries int32

// Hit events
var hitEventBaseAddress struct {
	Playing int32
	Replay  int32
}

var hitEventPreOffset [4]int32 = [4]int32{-1, -1, -1, -1}
var hitEventAddress [5]int32 = [5]int32{0, 0, 0, 0, 0}
var hitEventIndexIncrement int32 = 0
var prePlayTime int32 = 2147483647
var preStatus uint32 = 100
var hitEventFrame []hitEvent
var currentHitEventFrameIndex int
var prevHitEventFrameSize int = 0

func getHitEventOffset(depth int32, index int32) int32 {
	switch depth {
	case 0:
		return 0x34
	case 1:
		return 0x4
	case 2:
		return 0x8 + index*0x4
	case 3:
		return 0
	default:
		return -1 // this should not happen
	}
}

type addrWithMask struct {
	string
	int
}

func resolveHitEventAddress() error {
	var _hitEventMemPatterns struct {
		Playing []addrWithMask
		Replay  []addrWithMask
	} = struct {
		Playing []addrWithMask
		Replay  []addrWithMask
	}{ // hit event addr may have different pattern
		[]addrWithMask{
			{"83 7E 60 00 74 2C A1 ?? ?? ?? ?? 8B 50 1C 8B 4A 04", 7},
			{"5D C3 A1 ?? ?? ?? ?? 8B 50 1C 8B 4A 04", 3},
		},
		[]addrWithMask{
			{"D9 5D C0 EB 4E A1 ?? ?? ?? ?? 8B 48 34 4E", 6},
			{"74 4D A1 ?? ?? ?? ?? 8B 58 34 8D 46 FF", 3},
			{"A1 ?? ?? ?? ?? 8B 40 34 8B 70 0C", 1},
			{"75 0E 33 D2 89 15 ?? ?? ?? ?? 89 15", 6},
		},
	}

	if hitEventBaseAddress.Playing == 0 {
		var errMsg string = ""
		var resolved = false
		for i, t := range _hitEventMemPatterns.Playing {
			base, err := mem.Scan(process, t.string)
			if err != nil || base == 0 {
				errMsg += err.Error() + "; "
				continue
			}
			addr, err := mem.ReadInt32(process, base, int64(t.int))
			if err != nil || addr == 0 {
				errMsg += err.Error() + "; "
				continue
			} else {
				hitEventBaseAddress.Playing = addr
				resolved = true
				log.Printf("Resolved hit event(playing) at 0x%x by pattern #%d.\n", addr, i)
				break
			}
		}
		if !resolved {
			return fmt.Errorf("err resolving hit event(playing) address: %s", errMsg)
		}
	}
	if hitEventBaseAddress.Replay == 0 {
		var errMsg string = ""
		var resolved = false
		for i, t := range _hitEventMemPatterns.Replay {
			base, err := mem.Scan(process, t.string)
			if err != nil || base == 0 {
				errMsg += err.Error() + "; "
				continue
			}
			addr, err := mem.ReadInt32(process, base, int64(t.int))
			if err != nil || addr == 0 {
				errMsg += err.Error() + "; "
				continue
			} else {
				hitEventBaseAddress.Replay = addr
				resolved = true
				log.Printf("Resolved hit event(replay) at 0x%x by pattern #%d.\n", addr, i)
				break
			}
		}
		if !resolved {
			return fmt.Errorf("err resolving hit event(replay) address: %s", errMsg)
		}
	}
	return nil
}

// Init the whole thing and get osu! memory values to start working with it.
func Init() {
	if UnderWine == true || runtime.GOOS != "windows" { //Arrays start at 0xC in Linux for some reason, has to be wine specific
		leaderStart = 0xC
	} else {
		leaderStart = 0x8
	}

	allProcs, procerr = mem.FindProcess(osuProcessRegex, "osu!lazer", "osu!framework")
	for {
		start := time.Now()
		if procerr != nil {
			DynamicAddresses.IsReady = false
			for procerr != nil {
				allProcs, procerr = mem.FindProcess(osuProcessRegex, "osu!lazer", "osu!framework")
				log.Println("It seems that we lost the process, retrying! ERROR:", procerr)
				time.Sleep(1 * time.Second)
			}
			err := initBase()
			for err != nil {
				log.Println("Failure mid getting offsets, retrying! ERROR:", err)
				err = initBase()
				time.Sleep(1 * time.Second)
			}
		}
		if DynamicAddresses.IsReady == false {
			err := initBase()
			for err != nil {
				log.Println("Failure mid getting offsets, retrying! ERROR:", err)
				err = initBase()
				time.Sleep(1 * time.Second)
			}
		} else {
			err := mem.Read(process,
				&patterns.PreSongSelectAddresses,
				&menuData.PreSongSelectData)
			if err != nil {
				DynamicAddresses.IsReady = false
				log.Println("It appears that we lost the process, retrying! ERROR:", err)
				continue
			}
			MenuData.OsuStatus = menuData.Status

			mem.Read(process, &patterns, &alwaysData)
			MenuData.ChatChecker = alwaysData.ChatStatus
			MenuData.Bm.Time.PlayTime = alwaysData.PlayTime
			SettingsData.Folders.Skin = alwaysData.SkinFolder

			SettingsData.ShowInterface = cast.ToBool(int(alwaysData.ShowInterface))

			resolveHitEventAddress()

			if preStatus == 2 && menuData.Status != 2 {
				hitEventIndexIncrement = 0
				preStatus = menuData.Status
				hitEventFrame = []hitEvent{}
				currentHitEventFrameIndex = 0
				prevHitEventFrameSize = 0
				GameplayData.HitEvent = []hitEvent{}
			}

			switch menuData.Status {
			case 0:
				err = bmUpdateData()
				if err != nil {
					pp.Println(err)
				}
				mem.Read(process, &patterns, &mainMenuData)
				MenuData.MainMenuValues.BassDensity = calculateBassDensity(mainMenuData.AudioVelocityBase, &process)
			case 2:
				if MenuData.Bm.Time.PlayTime < 150 || menuData.Path == "" { //To catch up with the F2-->Enter
					err := bmUpdateData()
					if err != nil {
						pp.Println(err)
					}
				}
				if gameplayData.Retries > tempRetries || dirtyResults {
					tempRetries = gameplayData.Retries
					GameplayData = GameplayValues{}
					gameplayData = gameplayD{}
					dirtyResults = false
				}
				getGamplayData()
			case 1:
				err = bmUpdateData()
				if err != nil {
					pp.Println(err)
				}
			case 7:
				err = bmUpdateData()
				if err != nil {
					pp.Println(err)
				}
				mem.Read(process, &patterns, &resultsScreenData)

				if resultsScreenData.ModsXor1 == 0 { //not initialized yet
					for i := 0; i < 10; i++ {
						mem.Read(process, &patterns, &resultsScreenData)
						if resultsScreenData.ModsXor1 != 0 {
							break
						}
						time.Sleep(50 * time.Millisecond)
					}
				}
				ResultsScreenData.H300 = resultsScreenData.Hit300
				ResultsScreenData.H100 = resultsScreenData.Hit100
				ResultsScreenData.H50 = resultsScreenData.Hit50
				ResultsScreenData.H0 = resultsScreenData.HitMiss
				ResultsScreenData.MaxCombo = resultsScreenData.MaxCombo
				ResultsScreenData.Name = resultsScreenData.PlayerName
				ResultsScreenData.Score = resultsScreenData.Score
				ResultsScreenData.HGeki = resultsScreenData.HitGeki
				ResultsScreenData.HKatu = resultsScreenData.HitKatu

				ResultsScreenData.Mods.AppliedMods = resultsScreenData.ModsXor1 ^ resultsScreenData.ModsXor2
				if ResultsScreenData.Mods.AppliedMods == 0 {
					ResultsScreenData.Mods.PpMods = "NM"
				} else {
					ResultsScreenData.Mods.PpMods = Mods(resultsScreenData.ModsXor1 ^ resultsScreenData.ModsXor2).String()
				}
				if !dirtyResults {
					dirtyResults = true
				}
			default:
				tempRetries = -1
				GameplayData = GameplayValues{}
				gameplayData = gameplayD{}
				err = bmUpdateData()
				if err != nil {
					pp.Println(err)
				}

			}
		}
		if isTournamentMode {

			if err := getTourneyIPC(); err != nil {
				DynamicAddresses.IsReady = false
				log.Println("It appears that we lost the precess, retrying", err)
				continue
			}
		}
		if menuData.Status != 7 {
			ResultsScreenData = ResultsScreenValues{}
		}
		elapsed := time.Since(start)
		if MemCycle {
			log.Printf("Cycle took %s", elapsed)
		}
		time.Sleep(time.Duration(UpdateTime-int(elapsed.Milliseconds())) * time.Millisecond)

	}

}

var tempBeatmapString string
var tempGameMode int32 = 5

func bmUpdateData() error {
	mem.Read(process, &patterns, &menuData)

	bmString := menuData.Path
	if (strings.HasSuffix(bmString, ".osu") && tempBeatmapString != bmString) || (strings.HasSuffix(bmString, ".osu") && tempGameMode != menuData.MenuGameMode) { //On map/mode change
		for i := 0; i < 50; i++ {
			if menuData.BackgroundFilename != "" {
				break
			}
			time.Sleep(25 * time.Millisecond)
			mem.Read(process, &patterns, &menuData)
		}
		tempGameMode = menuData.MenuGameMode
		tempBeatmapString = bmString
		MenuData.Bm.BeatmapID = menuData.MapID
		MenuData.Bm.BeatmapSetID = menuData.SetID
		MenuData.Bm.Stats.MemoryAR = menuData.AR
		MenuData.Bm.Stats.MemoryCS = menuData.CS
		MenuData.Bm.Stats.MemoryHP = menuData.HP
		MenuData.Bm.Stats.MemoryOD = menuData.OD
		MenuData.Bm.Stats.TotalHitObjects = menuData.ObjectCount
		MenuData.Bm.Metadata.Artist = menuData.Artist
		MenuData.Bm.Metadata.ArtistOriginal = menuData.ArtistOriginal
		MenuData.Bm.Metadata.Title = menuData.Title
		MenuData.Bm.Metadata.TitleOriginal = menuData.TitleOriginal
		MenuData.Bm.Metadata.Mapper = menuData.Creator
		MenuData.Bm.Metadata.Version = menuData.Difficulty
		MenuData.GameMode = menuData.MenuGameMode
		MenuData.Bm.RandkedStatus = menuData.RankedStatus
		MenuData.Bm.BeatmapMD5 = menuData.MD5
		MenuData.Bm.Path = path{
			AudioPath:            menuData.AudioFilename,
			BGPath:               menuData.BackgroundFilename,
			BeatmapOsuFileString: menuData.Path,
			BeatmapFolderString:  menuData.Folder,
			FullMP3Path:          filepath.Join(SongsFolderPath, menuData.Folder, menuData.AudioFilename),
			FullDotOsu:           filepath.Join(SongsFolderPath, menuData.Folder, bmString),
			InnerBGPath:          filepath.Join(menuData.Folder, menuData.BackgroundFilename),
		}
	}
	if menuData.Status != 7 && menuData.Status != 14 {
		if alwaysData.MenuMods == 0 {
			MenuData.Mods.PpMods = "NM"
			MenuData.Mods.AppliedMods = int32(alwaysData.MenuMods)
		} else {
			MenuData.Mods.AppliedMods = int32(alwaysData.MenuMods)
			MenuData.Mods.PpMods = Mods(alwaysData.MenuMods).String()
		}
	}

	return nil
}
func getGamplayData() {
	err := mem.Read(process, &patterns, &gameplayData)
	if err != nil && !strings.Contains(err.Error(), "LeaderBoard") && !strings.Contains(err.Error(), "KeyOverlay") { //those could be disabled
		return //struct not initialized yet
	}
	//GameplayData.BitwiseKeypress = gameplayData.BitwiseKeypress
	GameplayData.Combo.Current = gameplayData.Combo
	GameplayData.Combo.Max = gameplayData.MaxCombo
	GameplayData.GameMode = gameplayData.Mode
	GameplayData.Hits.H100 = gameplayData.Hit100
	GameplayData.Hits.HKatu = gameplayData.HitKatu
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
	GameplayData.Accuracy = cast.ToFloat64(fmt.Sprintf("%.2f", gameplayData.Accuracy))
	GameplayData.Hp.Normal = gameplayData.PlayerHP
	GameplayData.Hp.Smooth = gameplayData.PlayerHPSmooth
	GameplayData.Name = gameplayData.PlayerName
	MenuData.Mods.AppliedMods = gameplayData.ModsXor1 ^ gameplayData.ModsXor2
	if MenuData.Mods.AppliedMods == 0 {
		MenuData.Mods.PpMods = "NM"
	} else {
		MenuData.Mods.PpMods = Mods(gameplayData.ModsXor1 ^ gameplayData.ModsXor2).String()
	}
	if strings.Contains(MenuData.Mods.PpMods, "V2") {
		GameplayData.Score = gameplayData.ScoreV2
	} else {
		GameplayData.Score = gameplayData.Score
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
	getKeyOveraly()

	if preStatus != 2 || prePlayTime > alwaysData.PlayTime {
		hitEventIndexIncrement = 0
	}

	preStatus = menuData.Status
	prePlayTime = alwaysData.PlayTime

	if menuData.Status == 2 {
		if hitEventBaseAddress.Playing != 0 {
			var currPlayTime = alwaysData.PlayTime

			hitEventAddress[0] = hitEventBaseAddress.Playing
			getHitEvent()

			var currHitEventFrameSize = len(hitEventFrame)
			if prevHitEventFrameSize != currHitEventFrameSize {
				if currHitEventFrameSize == 1 {
					GameplayData.HitEvent = hitEventFrame
				} else {
					GameplayData.HitEvent = hitEventFrame[prevHitEventFrameSize:currHitEventFrameSize]
				}
			} else {
				for i := currentHitEventFrameIndex; i < len(hitEventFrame)-1; i++ {
					var currHitEventFrame = hitEventFrame[i]
					if currPlayTime > currHitEventFrame.TimeStamp && currPlayTime <= hitEventFrame[i+1].TimeStamp {
						GameplayData.HitEvent = hitEventFrame[currentHitEventFrameIndex : i+1]
						currentHitEventFrameIndex = i + 1
						break
					}
				}
			}
			prevHitEventFrameSize = currHitEventFrameSize
		}
	}
}

func getHitEvent() {
	hitEventPreOffset = [4]int32{-1, -1, -1, -1}
	for {
		hitEvent, err := readCurrentHitEvent(hitEventIndexIncrement)

		if err != nil || hitEvent.X == -1 {
			break
		} else {
			hitEventFrame = append(hitEventFrame, hitEvent)
			hitEventIndexIncrement++
		}
	}
}

func readCurrentHitEvent(index int32) (hitEvent, error) {
	var changed = false
	var hitEvent hitEvent

	for depth := 0; depth < 4; depth++ {
		offset := getHitEventOffset(int32(depth), index)
		if offset != hitEventPreOffset[depth] || changed {
			hitEventPreOffset[depth] = offset
			changed = true
			addr, err := mem.ReadInt32(process, cast.ToInt64(hitEventAddress[depth]))
			hitEventAddress[depth+1] = addr
			if err != nil || hitEventAddress[depth+1] == 0 {
				hitEvent.X = -1
				return hitEvent, err
			}
			hitEventAddress[depth+1] += offset
		}
	}

	addr4x64 := cast.ToInt64(hitEventAddress[4])

	x, err := mem.ReadFloat32(process, addr4x64, 4)
	if err != nil {
		return hitEvent, err
	}
	y, err := mem.ReadFloat32(process, addr4x64, 8)
	if err != nil {
		return hitEvent, err
	}
	keyBit, err := mem.ReadUint8(process, addr4x64, 12)
	if err != nil {
		return hitEvent, err
	}
	time, err := mem.ReadInt32(process, addr4x64, 16)
	if err != nil {
		return hitEvent, err
	}

	// print memory
	// var u8Arr []string = []string{}
	// for i := 0; i < 100; i++ {
	// 	u8, _ := mem.ReadUint8(process, addr4x64, int64(i))
	// 	u8Arr = append(u8Arr, fmt.Sprintf("%2x", u8))
	// }
	// fmt.Println(strings.Join(u8Arr, " "))

	hitEvent.X = x
	hitEvent.Y = y
	hitEvent.TimeStamp = time
	hitEvent.K1 = keyBit&0b0101 == 0b0101
	hitEvent.K2 = keyBit&0b1010 == 0b1010
	hitEvent.M1 = (keyBit&0b1100 == 0b0000) && ((keyBit&0b0001 == 0b0001) || (keyBit&0b0001 == 0b0011))
	hitEvent.M2 = (keyBit&0b1100 == 0b0000) && ((keyBit&0b0010 == 0b0010) || (keyBit&0b0010 == 0b0011))
	return hitEvent, nil
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
	for i, j := leaderStart, 0; j < int(amOfSlots); i, j = i+0x4, j+1 {
		slot, _ := mem.ReadUint32(process, int64(items), int64(i))
		board.Slots[j], _ = readLeaderPlayerStruct(int64(slot))
	}
	GameplayData.Leaderboard = board
}

type ManiaStars struct {
	NoMod float64
	DT    float64
	HT    float64
}

func ReadManiaStars() (ManiaStars, error) {
	addresses := struct{ Base int64 }{int64(menuData.StarRatingStruct)} //Beatmap + 0x88
	var entries struct {
		Data uint32 `mem:"[Base + 0x14] + 0x8"`
	}
	err := mem.Read(process, &addresses, &entries)
	if err != nil || entries.Data == 0 {
		return ManiaStars{}, errors.New("[MEMORY] Could not find star rating for this map (internal) This probably means that difficulty calculation is in progress")
	}
	starRating := struct{ Base int64 }{int64(entries.Data)}
	var stars struct {
		NoMod float64 `mem:"Base + 0x18"`
		DT    float64 `mem:"Base + 0x30"`
		HT    float64 `mem:"Base + 0x48"`
	}
	err = mem.Read(process, &starRating, &stars)
	if err != nil {
		return ManiaStars{}, errors.New("[MEMORY] Empty star rating (internal)")
	}
	return ManiaStars{stars.NoMod, stars.DT, stars.HT}, nil
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
		return 0, errors.New("empty hit error array")
	}
	var totalAll float32 //double
	for _, hit := range HitErrorArray {
		totalAll += float32(hit)
	}
	var average = totalAll / float32(len(HitErrorArray))
	var variance float64 = 0
	for _, hit := range HitErrorArray {
		variance += math.Pow(float64(hit)-float64(average), 2)
	}
	variance = variance / float64(len(HitErrorArray))
	return math.Sqrt(variance) * 10, nil

}

var currentAudioVelocity float64

func calculateBassDensity(base uint32, proc *mem.Process) float64 {
	var bass float32
	for i, j := leaderStart, 0; j < 40; i, j = i+0x4, j+1 {
		value, err := mem.ReadFloat32(*proc, int64(base), int64(i))
		if err != nil {
			return 0.5
		}
		bass += 2 * value * (40 - float32(j)) / 40
	}
	if math.IsNaN(currentAudioVelocity) || math.IsNaN(float64(bass)) {
		currentAudioVelocity = 0
		return 0.5
	}
	currentAudioVelocity = math.Max(currentAudioVelocity, math.Min(float64(bass)*1.5, 6))
	currentAudioVelocity *= 0.95
	return (1 + currentAudioVelocity) * 0.5

}

func getKeyOveraly() {
	addresses := struct{ Base int64 }{int64(gameplayData.KeyOverlayArrayAddr)}
	var entries struct {
		//Pressed usually works with <20 update rate. It's recommended to create a buffer and predict presses by count to save CPU overhead
		K1Count int32 `mem:"[Base + 0x8] + 0x14"`
		K2Count int32 `mem:"[Base + 0xC] + 0x14"`
		M1Count int32 `mem:"[Base + 0x10] + 0x14"`
		M2Count int32 `mem:"[Base + 0x14] + 0x14"`
	}
	err := mem.Read(process, &addresses, &entries)
	if err != nil {
		return
	}

	var out keyCount

	out.K1 = entries.K1Count
	out.K2 = entries.K2Count
	out.M1 = entries.M1Count
	out.M2 = entries.M2Count
	GameplayData.KeyCount = out //needs complete rewrite in 1.4.0
}
