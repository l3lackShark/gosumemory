package memory

import (
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/k0kubun/pp"

	"github.com/l3lackShark/kiwi"
)

func modsResolver(xor uint32) string {
	return Mods(xor).String()
}

//UpdateTime Intervall between value updates
var UpdateTime int
var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
var leaderStart int32
var hasLeaderboard = false

//SongsFolderPath is full path to osu! Songs. Gets set automatically on Windows (through memory)
var SongsFolderPath string

func oncePerBeatmapChange() error {
	var err error
	DynamicAddresses.LeaderBoardStruct, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.LeaderBoard), 0x4, 0x74, 0x24, 0x4, 0x4)
	if err != nil {
		//pp.Println("Could not get leaderboard stuff! ", err, osuStaticAddresses.LeaderBoard)
		return err
	}
	GameplayData.Leaderboard.OurPlayer.Addr, err = proc.ReadUint32Ptr(uintptr(DynamicAddresses.LeaderBoardStruct+uint32(leaderStart)), 0x24, 0x10)
	if err != nil {
		pp.Println("Could not get current player! ", err)
		return err
	}

	nameAddr, err := proc.ReadUint32(uintptr(GameplayData.Leaderboard.OurPlayer.Addr + 0x8))
	GameplayData.Leaderboard.OurPlayer.Name, err = proc.ReadNullTerminatedUTF16String(uintptr(nameAddr + 0x8))
	if err != nil {
		pp.Println("Could not get current player name! ", err)
		return err
	}

	return nil
}

func leaderPlayerCountResolver() error {
	DynamicAddresses.LeaderSlotAddr = nil
	DynamicAddresses.LeaderBaseSlotAddr = nil
	for i := leaderStart; i < 0xE4; i += 0x4 {
		slot, err := proc.ReadUint32Ptr(uintptr(DynamicAddresses.LeaderBoardStruct + uint32(i)))
		if err != nil || slot == 0x0 {
			return err
		}
		DynamicAddresses.LeaderBaseSlotAddr = append(DynamicAddresses.LeaderBaseSlotAddr, slot)
		slotaddr, err := proc.ReadUint32(uintptr(slot) + 0x20)
		if err != nil {
			return err
		}
		if slotaddr == 0x0 { //osu has 64 slots in leaderboard array for some reason, those that are unused point to 0
			GameplayData.Leaderboard.OurPlayer.AmountOfSlots = int32((i - leaderStart + 0x4) / 4)
			return nil
		}
		DynamicAddresses.LeaderSlotAddr = append(DynamicAddresses.LeaderSlotAddr, slotaddr)
	}

	return nil
}

var comboResult []int16
var maxComboResult []int32
var scoreResult []int32
var h300Result []int16
var h100Result []int16
var h50Result []int16
var h0Result []int16
var nameResult []string

func leaderSlotsData() error {
	nameResult = nil
	comboResult = nil
	maxComboResult = nil
	scoreResult = nil
	h300Result = nil
	h100Result = nil
	h50Result = nil
	h0Result = nil
	if len(DynamicAddresses.LeaderSlotAddr) >= 1 {

		for i := 0; i < len(DynamicAddresses.LeaderSlotAddr); i++ {

			nameoffset, err := proc.ReadUint32(uintptr(DynamicAddresses.LeaderBaseSlotAddr[i] + 0x8))
			if err != nil || nameoffset == 0x0 {
				return err
			}
			name, err := proc.ReadNullTerminatedUTF16String(uintptr(nameoffset) + 0x8)
			if err != nil {
				return err
			}
			combo, err := proc.ReadInt16(uintptr(DynamicAddresses.LeaderSlotAddr[i]) + 0x90)
			if err != nil {
				return err
			}
			maxcombo, err := proc.ReadInt32(uintptr(DynamicAddresses.LeaderSlotAddr[i]) + 0x68)
			if err != nil {
				return err
			}
			score, err := proc.ReadInt32(uintptr(DynamicAddresses.LeaderSlotAddr[i]) + 0x74)
			if err != nil {
				return err
			}
			hit300, err := proc.ReadInt16(uintptr(DynamicAddresses.LeaderSlotAddr[i]) + 0x86)
			if err != nil {
				return err
			}
			hit100, err := proc.ReadInt16(uintptr(DynamicAddresses.LeaderSlotAddr[i]) + 0x84)
			if err != nil {
				return err
			}
			hit50, err := proc.ReadInt16(uintptr(DynamicAddresses.LeaderSlotAddr[i]) + 0x88)
			if err != nil {
				return err
			}
			hit0, err := proc.ReadInt16(uintptr(DynamicAddresses.LeaderSlotAddr[i]) + 0x8E)
			if err != nil {
				return err

			}
			nameResult = append(nameResult, name)
			comboResult = append(comboResult, combo)
			maxComboResult = append(maxComboResult, maxcombo)
			scoreResult = append(scoreResult, score)
			h300Result = append(h300Result, hit300)
			h100Result = append(h100Result, hit100)
			h50Result = append(h50Result, hit50)
			h0Result = append(h0Result, hit0)

		}
		GameplayData.Leaderboard.Slots.Combo = comboResult
		GameplayData.Leaderboard.Slots.MaxCombo = maxComboResult
		GameplayData.Leaderboard.Slots.Score = scoreResult
		GameplayData.Leaderboard.Slots.H300 = h300Result
		GameplayData.Leaderboard.Slots.H100 = h100Result
		GameplayData.Leaderboard.Slots.H50 = h50Result
		GameplayData.Leaderboard.Slots.H0 = h0Result
		GameplayData.Leaderboard.Slots.Name = nameResult
	}

	return nil
}

//readHitErrorArray Gets an array of ints representing UnstableRate. (a little innacurate, shows values with 2 hitobjects delay)
func readHitErrorArray() ([]int32, error) {
	base, err := proc.ReadUint32(uintptr(DynamicAddresses.PlayContainer38 + 0x38))
	if err != nil {
		return nil, err
	}
	hitErrorStruct, err := proc.ReadUint32(uintptr(base + 0x4))
	if err != nil {
		return nil, err
	}
	leng, err := proc.ReadUint32(uintptr(base + 0xC))
	if err != nil {
		return nil, err
	}
	var buf32 []int32
	for i := 0x8; i < int(leng*0x4); i += 0x4 {
		value, err := proc.ReadInt32(uintptr(hitErrorStruct + uint32(i)))
		if err != nil {
			return nil, err
		}
		buf32 = append(buf32, value)
	}
	return buf32, nil
}

//Init the whole thing and get osu! memory values to start working with it.
func Init() {
	if runtime.GOOS == "windows" {
		leaderStart = 0x8
	} else {
		leaderStart = 0xC
	}
	var tempBeatmapID uint32 = 0
	for {
		var err error
		proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
		if procerr != nil {
			DynamicAddresses.IsReady = false
			for procerr != nil {
				proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
				log.Println("It seems that we lost the process, retrying!")
				time.Sleep(1 * time.Second)
			}
			DynamicAddresses.IsReady = false
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
				log.Println("Failure mid getting offsets, retrying!")
				time.Sleep(1 * time.Second)

			}
		}

		MenuData.OsuStatus, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Status-0x4), 0x0)
		if err != nil {
			log.Println("Could not get osuStatus Value, retrying")
			initBase()
		}

		switch MenuData.OsuStatus {
		case 2:
			DynamicAddresses.PlayContainer38, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayContainer-0x4), 0x0, 0x38) //TODO: Should only be read once per map change
			if err != nil {
				//log.Println(err)
			}
			xor1, err := proc.ReadUint32Ptr(uintptr(DynamicAddresses.PlayContainer38+0x1C), 0xC)
			xor2, err := proc.ReadUint32Ptr(uintptr(DynamicAddresses.PlayContainer38+0x1C), 0x8)

			accOffset, err := proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayContainer-0x4), 0x0, 0x48)
			MenuData.Mods.AppliedMods = int32(xor1 ^ xor2)
			GameplayData.Combo.Current, err = proc.ReadInt32(uintptr(DynamicAddresses.PlayContainer38 + 0x90))
			GameplayData.Combo.Max, err = proc.ReadInt32(uintptr(DynamicAddresses.PlayContainer38 + 0x68))
			GameplayData.GameMode, err = proc.ReadInt32(uintptr(DynamicAddresses.PlayContainer38 + 0x64))
			GameplayData.Score, err = proc.ReadInt32(uintptr(DynamicAddresses.PlayContainer38 + 0x74))
			GameplayData.Hits.H100, err = proc.ReadInt16(uintptr(DynamicAddresses.PlayContainer38 + 0x84))
			GameplayData.Hits.H300, err = proc.ReadInt16(uintptr(DynamicAddresses.PlayContainer38 + 0x86))
			GameplayData.Hits.H50, err = proc.ReadInt16(uintptr(DynamicAddresses.PlayContainer38 + 0x88))
			GameplayData.Hits.H0, err = proc.ReadInt16(uintptr(DynamicAddresses.PlayContainer38 + 0x8E))
			GameplayData.Accuracy, err = proc.ReadFloat64(uintptr(accOffset + 0x14))
			MenuData.Bm.Time.PlayTime, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayTime+0x5), 0x0)
			GameplayData.Hits.HitErrorArray, err = readHitErrorArray()
			if err != nil {
				//log.Println(err)
			}

			if MenuData.Bm.Time.PlayTime <= 15000 { //hardcoded for now as current pointer chain is unstable and tends to change within first 15 seconds
				err := oncePerBeatmapChange()
				if err != nil {
					hasLeaderboard = false
				} else {
					hasLeaderboard = true
				}
			}
			leaderPlayerCountResolver() //should probably run this on another thread
			if hasLeaderboard == true {
				err = leaderSlotsData()
				if err != nil {
					log.Println("Leaderboard data error: ", err)
				}
				GameplayData.Leaderboard.OurPlayer.Position, err = proc.ReadInt32(uintptr(GameplayData.Leaderboard.OurPlayer.Addr + 0x2C))
			}
			MenuData.Mods.PpMods = Mods(MenuData.Mods.AppliedMods).String()

		default: //This data is available at all times
			hasLeaderboard = false
			DynamicAddresses.BeatmapAddr, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.Base-0xC), 0x0)
			if err != nil {
				log.Println("Dynamic beatmap addr error: ", err)
			}
			MenuData.Bm.BeatmapID, err = proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr + 0xC4))
			if err != nil {
				//log.Println("Dynamic beatmap id error: ", err) //Gets triggered on F2
			}

			if tempBeatmapID != MenuData.Bm.BeatmapID { //On map change
				time.Sleep(time.Duration(UpdateTime) * time.Millisecond)
				MenuData.Bm.BeatmapSetID, err = proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr + 0xC8))
				beatmapStrOffset, err := proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr) + 0x7C)
				MenuData.Bm.BeatmapString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapStrOffset) + 0x8)
				beatmapBGStringOffset, err := proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr) + 0x68)
				MenuData.Bm.Path.BGPath, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapBGStringOffset) + 0x8)
				beatmapOsuFileStrOffset, err := proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr) + 0x8C)
				MenuData.Bm.Path.BeatmapOsuFileString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapOsuFileStrOffset) + 0x8)
				beatmapFolderStrOffset, err := proc.ReadUint32(uintptr(DynamicAddresses.BeatmapAddr) + 0x74)
				MenuData.Bm.Path.BeatmapFolderString, err = proc.ReadNullTerminatedUTF16String(uintptr(beatmapFolderStrOffset) + 0x8)
				// MenuData.Bm.Stats.BeatmapAR, err = proc.ReadFloat32(uintptr(DynamicAddresses.BeatmapAddr + 0x2C))
				// MenuData.Bm.Stats.BeatmapCS, err = proc.ReadFloat32(uintptr(DynamicAddresses.BeatmapAddr + 0x30))
				// MenuData.Bm.Stats.BeatmapHP, err = proc.ReadFloat32(uintptr(DynamicAddresses.BeatmapAddr + 0x34))
				// MenuData.Bm.Stats.BeatmapOD, err = proc.ReadFloat32(uintptr(DynamicAddresses.BeatmapAddr + 0x38))
				MenuData.Bm.Time.PlayTime, err = proc.ReadUint32Ptr(uintptr(osuStaticAddresses.PlayTime+0x5), 0x0)
				if err != nil {
					log.Println("MenuData failure")
				}

				if strings.HasSuffix(MenuData.Bm.Path.BeatmapOsuFileString, ".osu") == true && len(MenuData.Bm.Path.BGPath) > 0 {
					MenuData.Bm.Path.InnerBGPath = MenuData.Bm.Path.BeatmapFolderString + "/" + MenuData.Bm.Path.BGPath

				} else {
					log.Println("skipping bg reloading")
				}

				tempBeatmapID = MenuData.Bm.BeatmapID

			}
			menuMods, err := proc.ReadUint32Ptr(uintptr(osuStaticAddresses.InMenuMods+0x9), 0x0)
			if err != nil {
				pp.Println(err)
			} else {
				if menuMods == 0 {
					MenuData.Mods.PpMods = "NM"
					MenuData.Mods.AppliedMods = int32(menuMods)
				} else {
					MenuData.Mods.AppliedMods = int32(menuMods)
					MenuData.Mods.PpMods = Mods(menuMods).String()
				}

			}

		}

		time.Sleep(time.Duration(UpdateTime) * time.Millisecond)
	}

}
