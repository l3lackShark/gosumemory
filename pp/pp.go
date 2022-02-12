package pp

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/Wieku/gosu-pp/beatmap"
	"github.com/Wieku/gosu-pp/beatmap/difficulty"
	"github.com/Wieku/gosu-pp/performance/osu"
	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/spf13/cast"
)

//#cgo LDFLAGS: -lm
//#cgo CPPFLAGS: -DOPPAI_STATIC_HEADER
//#include <stdlib.h>
//#include "oppai.c"
import "C"

var ez C.ezpp_t

type PP struct {
	Total         C.float
	FC            C.float
	Strain        []float64
	StarRating    C.float
	AimStars      C.float
	SpeedStars    C.float
	AimPP         C.float
	SpeedPP       C.float
	Accuracy      C.float
	N300          C.int
	N100          C.int
	N50           C.int
	NMiss         C.int
	AR            C.float
	CS            C.float
	OD            C.float
	HP            C.float
	Artist        string
	ArtistUnicode string
	Title         string
	TitleUnicode  string
	Version       string
	Creator       string
	NCircles      C.int
	NSliders      C.int
	NSpinners     C.int
	ODMS          C.float
	Mode          C.int
	Combo         C.int
	MaxCombo      C.int
	Mods          C.int
	ScoreVersion  C.int
}

var strainArray []float64
var tempBeatmapFile string
var tempGameMode int32 = 4
var currMaxCombo C.int

func readData(data *PP, ez C.ezpp_t, needStrain bool, path string) error {

	if strings.HasSuffix(path, ".osu") {
		cpath := C.CString(path)

		defer C.free(unsafe.Pointer(cpath))
		if rc := C.ezpp(ez, cpath); rc < 0 {
			memory.MenuData.PP.PpStrains = []float64{0}
			return errors.New(C.GoString(C.errstr(rc)))
		}
		C.ezpp_set_base_ar(ez, C.float(memory.MenuData.Bm.Stats.MemoryAR))
		C.ezpp_set_base_od(ez, C.float(memory.MenuData.Bm.Stats.MemoryOD))
		C.ezpp_set_base_cs(ez, C.float(memory.MenuData.Bm.Stats.MemoryCS))
		C.ezpp_set_base_hp(ez, C.float(memory.MenuData.Bm.Stats.MemoryHP))
		C.ezpp_set_accuracy_percent(ez, C.float(memory.GameplayData.Accuracy))
		C.ezpp_set_mods(ez, C.int(memory.MenuData.Mods.AppliedMods))
		*data = PP{
			Artist:     C.GoString(C.ezpp_artist(ez)),
			Title:      C.GoString(C.ezpp_title(ez)),
			Version:    C.GoString(C.ezpp_version(ez)),
			Creator:    C.GoString(C.ezpp_creator(ez)),
			AR:         C.ezpp_ar(ez),
			CS:         C.ezpp_cs(ez),
			OD:         C.ezpp_od(ez),
			HP:         C.ezpp_hp(ez),
			StarRating: C.ezpp_stars(ez),
		}
		memory.MenuData.Bm.Stats.BeatmapSR = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.StarRating)))
		memory.MenuData.Bm.Stats.BeatmapAR = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.AR)))
		memory.MenuData.Bm.Stats.BeatmapCS = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.CS)))
		memory.MenuData.Bm.Stats.BeatmapOD = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.OD)))
		memory.MenuData.Bm.Stats.BeatmapHP = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.HP)))

		if needStrain == true {
			C.ezpp_set_end_time(ez, 0)
			C.ezpp_set_combo(ez, 0)
			C.ezpp_set_nmiss(ez, 0)
			memory.MenuData.Bm.Stats.BeatmapMaxCombo = int32(C.ezpp_max_combo(ez))
			memory.MenuData.Bm.Stats.FullSR = cast.ToFloat32(fmt.Sprintf("%.2f", float32(C.ezpp_stars(ez))))
			var bpmChanges []int
			var bpmMultiplier float64 = 1
			if strings.Contains(memory.MenuData.Mods.PpMods, "DT") || strings.Contains(memory.MenuData.Mods.PpMods, "NC") {
				bpmMultiplier = 1.5
			} else if strings.Contains(memory.MenuData.Mods.PpMods, "HT") {
				bpmMultiplier = 0.75
			}
			for i := 0; i < int(C.ezpp_ntiming_points(ez)); i++ {
				msPerBeat := float64(C.ezpp_timing_ms_per_beat(ez, C.int(i)))
				timingChanges := int(C.ezpp_timing_change(ez, C.int(i)))
				if timingChanges == 1 {
					bpmFormula := int(math.Round(1 / msPerBeat * 1000 * 60 * bpmMultiplier))
					if bpmFormula > 0 {
						bpmChanges = append(bpmChanges, bpmFormula)
					}
				}
			}
			memory.MenuData.Bm.Stats.BeatmapBPM.Minimal, memory.MenuData.Bm.Stats.BeatmapBPM.Maximal = minMax(bpmChanges)
			strainArray = nil
			seek := 0
			var window []float64
			var total []float64
			// for seek < int(C.ezpp_time_at(ez, C.ezpp_nobjects(ez)-1)) { //len-1
			for int32(seek) < memory.MenuData.Bm.Time.Mp3Time {
				for obj := 0; obj <= int(C.ezpp_nobjects(ez)-1); obj++ {
					if tempBeatmapFile != memory.MenuData.Bm.Path.BeatmapOsuFileString {
						return nil //Interrupt calcualtion if user has changed the map.
					}
					if int(C.ezpp_time_at(ez, C.int(obj))) >= seek && int(C.ezpp_time_at(ez, C.int(obj))) <= seek+3000 {
						window = append(window, float64(C.ezpp_strain_at(ez, C.int(obj), 0))+float64(C.ezpp_strain_at(ez, C.int(obj), 1)))
					}
				}
				sum := 0.0
				for _, num := range window {
					sum += num
				}
				total = append(total, sum/math.Max(float64(len(window)), 1))
				window = nil
				seek += 500
			}
			strainArray = total
			memory.MenuData.Bm.Time.FirstObj = int32(C.ezpp_time_at(ez, 0))
			memory.MenuData.Bm.Time.FullTime = int32(C.ezpp_time_at(ez, C.ezpp_nobjects(ez)-1))
		} else {
			C.ezpp_set_end_time(ez, C.float(memory.MenuData.Bm.Time.PlayTime))
			currMaxCombo = C.ezpp_max_combo(ez) //for RestSS
			C.ezpp_set_combo(ez, C.int(memory.GameplayData.Combo.Max))
			C.ezpp_set_nmiss(ez, C.int(memory.GameplayData.Hits.H0))
		}

		*data = PP{
			Total:  C.ezpp_pp(ez),
			Strain: strainArray,

			AimStars:   C.ezpp_aim_stars(ez),
			SpeedStars: C.ezpp_speed_stars(ez),
			AimPP:      C.ezpp_aim_pp(ez),
			SpeedPP:    C.ezpp_speed_pp(ez),
			Accuracy:   C.ezpp_accuracy_percent(ez),
			N300:       C.ezpp_n300(ez),
			N100:       C.ezpp_n100(ez),
			N50:        C.ezpp_n50(ez),
			NMiss:      C.ezpp_nmiss(ez),
			//ArtistUnicode: C.GoString(C.ezpp_artist_unicode(ez)),
			//	TitleUnicode:  C.GoString(C.ezpp_title_unicode(ez)),
			NCircles:     C.ezpp_ncircles(ez),
			NSliders:     C.ezpp_nsliders(ez),
			NSpinners:    C.ezpp_nspinners(ez),
			ODMS:         C.ezpp_odms(ez),
			Mode:         C.ezpp_mode(ez),
			Combo:        C.ezpp_combo(ez),
			MaxCombo:     C.ezpp_max_combo(ez),
			Mods:         C.ezpp_mods(ez),
			ScoreVersion: C.ezpp_score_version(ez),
		}
		memory.MenuData.PP.PpStrains = data.Strain
	}
	return nil
}

var maniaSR float64
var maniaHitObjects float64
var tempMods string

func GetData() {

	ez = C.ezpp_new()
	C.ezpp_set_autocalc(ez, 1)
	//defer C.ezpp_free(ez)

	for {

		if memory.DynamicAddresses.IsReady == true {
			switch memory.MenuData.GameMode {
			case 0, 1:
				var data PP
				if tempBeatmapFile != memory.MenuData.Bm.Path.BeatmapOsuFileString || memory.MenuData.Mods.PpMods != tempMods || memory.MenuData.GameMode != tempGameMode { //On map/mods change
					tempGameMode = memory.MenuData.GameMode // looks very ugly but will rewrite everything in 1.4.0
					tempBadJudgments = 0
					path := memory.MenuData.Bm.Path.FullDotOsu
					tempBeatmapFile = memory.MenuData.Bm.Path.BeatmapOsuFileString
					tempMods = memory.MenuData.Mods.PpMods
					tempGameMode = memory.MenuData.GameMode
					mp3Time, err := calculateMP3Time()
					if err == nil {
						memory.MenuData.Bm.Time.Mp3Time = mp3Time
					}
					//Get Strains
					readData(&data, ez, true, path)

					//pp.Println(memory.MenuData.Bm.Metadata)
				}

				switch memory.MenuData.OsuStatus {
				case 2:
					path := memory.MenuData.Bm.Path.FullDotOsu
					readData(&data, ez, false, path)
					if memory.GameplayData.Combo.Max > 1 && float64(data.Total) > 0 {
						//pre-Wieku rewrite crutch
						if memory.GameplayData.GameMode == 0 {

							res, err := wiekuCalcCrutch(path, memory.GameplayData.Combo.Max, memory.GameplayData.Hits.H300, memory.GameplayData.Hits.H100, memory.GameplayData.Hits.H50, memory.GameplayData.Hits.H0)
							if err != nil {
								pp.Println(err)
								memory.GameplayData.PP.Pp = cast.ToInt32(float64(data.Total))
							}
							memory.GameplayData.PP.Pp = cast.ToInt32(res)

						} else {
							memory.GameplayData.PP.Pp = cast.ToInt32(float64(data.Total))
						}
					}
				case 7:
					//idle
				case 5:
					memory.GameplayData.PP.Pp = 0
				}

			case 3:

				if tempBeatmapFile != memory.MenuData.Bm.Path.BeatmapOsuFileString || memory.MenuData.Mods.PpMods != tempMods || memory.MenuData.GameMode != tempGameMode { //On map/mods/mode change
					tempGameMode = memory.MenuData.GameMode // looks very ugly but will rewrite everything in 1.4.0
					tempBeatmapFile = memory.MenuData.Bm.Path.BeatmapOsuFileString
					tempMods = memory.MenuData.Mods.PpMods

					maniaSR = 0.0
					memory.MenuData.Bm.Time.FullTime = 0        //Not implemented for mania yet
					memory.MenuData.Bm.Stats.BeatmapAR = 0      //Not implemented for mania yet
					memory.MenuData.Bm.Stats.BeatmapCS = 0      //Not implemented for mania yet
					memory.MenuData.Bm.Stats.BeatmapOD = 0      //Not implemented for mania yet
					memory.MenuData.Bm.Stats.BeatmapHP = 0      //Not implemented for mania yet
					memory.MenuData.PP.PpStrains = []float64{0} //Not implemented for mania yet

					maniaStars, err := memory.ReadManiaStars()
					if err != nil {
						pp.Println(err)
					}
					if maniaStars.NoMod == 0 { //diff calc in progress
						for i := 0; i < 50; i++ {
							maniaStars, _ = memory.ReadManiaStars()
							if maniaStars.NoMod > 0 {
								break
							}
							time.Sleep(100 * time.Millisecond)
						}
					}

					maniaHitObjects = float64(memory.MenuData.Bm.Stats.TotalHitObjects)

					if strings.Contains(memory.MenuData.Mods.PpMods, "DT") {
						maniaSR = maniaStars.DT
					} else if strings.Contains(memory.MenuData.Mods.PpMods, "HT") {
						maniaSR = maniaStars.HT
					} else {
						maniaSR = maniaStars.NoMod //assuming NM
					}
					memory.MenuData.Bm.Stats.BeatmapSR = cast.ToFloat32(fmt.Sprintf("%.2f", float32(maniaSR)))
					memory.MenuData.Bm.Stats.FullSR = memory.MenuData.Bm.Stats.BeatmapSR
					memory.MenuData.PP.PpSS = int32(calculateManiaPP(float64(memory.MenuData.Bm.Stats.MemoryOD), maniaSR, maniaHitObjects, 1000000.0)) // LiveSR not implemented yet
				}
			}
			if memory.GameplayData.GameMode == 3 {
				if maniaSR > 0 {
					memory.GameplayData.PP.PPifFC = int32(calculateManiaPP(float64(memory.MenuData.Bm.Stats.MemoryOD), maniaSR, maniaHitObjects, 1000000.0)) //PP if SS
					if memory.GameplayData.Score >= 500000 {
						memory.GameplayData.PP.Pp = int32(calculateManiaPP(float64(memory.MenuData.Bm.Stats.MemoryOD), maniaSR, maniaHitObjects, float64(memory.GameplayData.Score)))
					} else {
						memory.GameplayData.PP.Pp = 0
					}
				}
			}
		}

		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}
}

var (
	tempWiekuFileName string
	tempWiekuMods     int32
	beatMap           *beatmap.BeatMap
	attribs           []osu.Attributes
)

func wiekuCalcCrutch(path string, combo int16, h300 int16, h100 int16, h50 int16, h0 int16) (int32, error) {
	if tempWiekuFileName != path || tempWiekuMods != memory.MenuData.Mods.AppliedMods {
		tempWiekuFileName = path
		tempWiekuMods = memory.MenuData.Mods.AppliedMods

		osuFile, err := os.Open(path)
		if err != nil {
			return 0, fmt.Errorf("Failed to calc via wieku calculator, falling back to oppai, ERROR: %w", err)
		}
		defer osuFile.Close()

		beatMap, err = beatmap.ParseFromReader(osuFile)
		if err != nil {
			return 0, fmt.Errorf("Failed to calc via wieku calculator, falling back to oppai, ERROR: %w", err)
		}

		beatMap.Difficulty.SetMods(difficulty.Modifier(memory.MenuData.Mods.AppliedMods))
		attribs = osu.CalculateStep(beatMap.HitObjects, beatMap.Difficulty)
	}

	ppWieku := &osu.PPv2{}

	currAttrib := int(math.Max(0, float64(h300+h100+h50+h0-1)))

	if len(attribs)-1 < currAttrib {
		return 0, nil //rade condition hell
	}

	ppWieku.PPv2x(attribs[currAttrib], int(combo), int(h300), int(h100), int(h50), int(h0), beatMap.Difficulty)

	return cast.ToInt32(ppWieku.Results.Total), nil
}
