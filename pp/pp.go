package pp

import (
	"errors"
	"strings"
	"time"
	"unsafe"

	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/spf13/cast"
)

//#cgo LDFLAGS: -lm
//#cgo CPPFLAGS: -DOPPAI_STATIC_HEADER
//#include <stdlib.h>
//#include "oppai.c"
// import "C"

//#cgo LDFLAGS: -lm
//#cgo CPPFLAGS: -DOPPAI_STATIC_HEADER
//#include <stdlib.h>
//#include "oppai.c"
import "C"

var ez C.ezpp_t

type PP struct {
	Total         C.float
	FC            C.float
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

func readData(data *PP, ez C.ezpp_t) error {
	path := (memory.SongsFolderPath + "/" + memory.MenuData.Bm.Path.BeatmapFolderString + "/" + memory.MenuData.Bm.Path.BeatmapOsuFileString) //TODO: Automatic Songs folder finder

	if strings.HasSuffix(path, ".osu") && memory.DynamicAddresses.IsReady == true {
		cpath := C.CString(path)

		defer C.free(unsafe.Pointer(cpath))
		if rc := C.ezpp(ez, cpath); rc < 0 {
			return errors.New(C.GoString(C.errstr(rc)))
		}
		C.ezpp_set_base_ar(ez, C.float(memory.MenuData.Bm.Stats.BeatmapAR))
		C.ezpp_set_base_od(ez, C.float(memory.MenuData.Bm.Stats.BeatmapOD))
		C.ezpp_set_base_cs(ez, C.float(memory.MenuData.Bm.Stats.BeatmapCS))
		C.ezpp_set_base_hp(ez, C.float(memory.MenuData.Bm.Stats.BeatmapHP))
	}
	C.ezpp_set_accuracy_percent(ez, C.float(memory.GameplayData.Accuracy))
	C.ezpp_set_mods(ez, C.int(memory.GameplayData.Mods.AppliedMods))

	//C.ezpp_set_score_version(ez)
	C.ezpp_set_end_time(ez, C.float(memory.MenuData.Bm.Time.PlayTime))
	C.ezpp_set_combo(ez, C.int(memory.GameplayData.Combo.Max))
	C.ezpp_set_nmiss(ez, C.int(memory.GameplayData.Hits.H0))
	*data = PP{
		Total:         C.ezpp_pp(ez),
		StarRating:    C.ezpp_stars(ez),
		AimStars:      C.ezpp_aim_stars(ez),
		SpeedStars:    C.ezpp_speed_stars(ez),
		AimPP:         C.ezpp_aim_pp(ez),
		SpeedPP:       C.ezpp_speed_pp(ez),
		Accuracy:      C.ezpp_accuracy_percent(ez),
		N300:          C.ezpp_n300(ez),
		N100:          C.ezpp_n100(ez),
		N50:           C.ezpp_n50(ez),
		NMiss:         C.ezpp_nmiss(ez),
		AR:            C.ezpp_ar(ez),
		CS:            C.ezpp_cs(ez),
		OD:            C.ezpp_od(ez),
		HP:            C.ezpp_hp(ez),
		Artist:        C.GoString(C.ezpp_artist(ez)),
		ArtistUnicode: C.GoString(C.ezpp_artist_unicode(ez)),
		Title:         C.GoString(C.ezpp_title(ez)),
		TitleUnicode:  C.GoString(C.ezpp_title_unicode(ez)),
		Version:       C.GoString(C.ezpp_version(ez)),
		Creator:       C.GoString(C.ezpp_creator(ez)),
		NCircles:      C.ezpp_ncircles(ez),
		NSliders:      C.ezpp_nsliders(ez),
		NSpinners:     C.ezpp_nspinners(ez),
		ODMS:          C.ezpp_odms(ez),
		Mode:          C.ezpp_mode(ez),
		Combo:         C.ezpp_combo(ez),
		MaxCombo:      C.ezpp_max_combo(ez),
		Mods:          C.ezpp_mods(ez),
		ScoreVersion:  C.ezpp_score_version(ez),
	}
	return nil
}

func GetData() {
	ez := C.ezpp_new()
	defer C.ezpp_free(ez)
	C.ezpp_set_autocalc(ez, 1)

	for {
		if memory.DynamicAddresses.IsReady == true && memory.GameplayData.Combo.Max >= 1 && memory.MenuData.OsuStatus == 2 {
			var data PP
			err := readData(&data, ez)
			if err != nil {
				pp.Println(err)
			}
			//	pp.Println(data.Total)
			if cast.ToInt32(float64(data.Total)) >= 0 && memory.MenuData.OsuStatus == 2 {
				memory.GameplayData.PP.Pp = cast.ToInt32(float64(data.Total))
			}
			if memory.MenuData.OsuStatus == 5 {
				memory.GameplayData.Hits.H100 = 0
				memory.GameplayData.Hits.H50 = 0
				memory.GameplayData.Hits.H0 = 0
				memory.GameplayData.PP.Pp = 0

			}
		}

		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}
}
