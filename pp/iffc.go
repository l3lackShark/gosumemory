package pp

//TODO: I need to figure out how to use only one calc.

import (
	"errors"
	"strings"
	"time"
	"unsafe"

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

var ezfc C.ezpp_t

type PPfc struct {
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

func readFCData(data *PPfc, ezfc C.ezpp_t) error {
	path := (memory.SongsFolderPath + "/" + memory.MenuData.Bm.Path.BeatmapFolderString + "/" + memory.MenuData.Bm.Path.BeatmapOsuFileString) //TODO: Automatic Songs folder finder

	if strings.HasSuffix(path, ".osu") && memory.DynamicAddresses.IsReady == true {
		cpath := C.CString(path)

		defer C.free(unsafe.Pointer(cpath))
		if rc := C.ezpp(ezfc, cpath); rc < 0 {
			return errors.New(C.GoString(C.errstr(rc)))
		}
		C.ezpp_set_base_ar(ezfc, C.float(memory.MenuData.Bm.Stats.BeatmapAR))
		C.ezpp_set_base_od(ezfc, C.float(memory.MenuData.Bm.Stats.BeatmapOD))
		C.ezpp_set_base_cs(ezfc, C.float(memory.MenuData.Bm.Stats.BeatmapCS))
		C.ezpp_set_base_hp(ezfc, C.float(memory.MenuData.Bm.Stats.BeatmapHP))

		C.ezpp_set_accuracy_percent(ezfc, C.float(memory.GameplayData.Accuracy))
		C.ezpp_set_mods(ezfc, C.int(memory.MenuData.Mods.AppliedMods))

		//C.ezpp_set_score_version(ezfc)
		*data = PPfc{
			Total:         C.ezpp_pp(ezfc),
			StarRating:    C.ezpp_stars(ezfc),
			AimStars:      C.ezpp_aim_stars(ezfc),
			SpeedStars:    C.ezpp_speed_stars(ezfc),
			AimPP:         C.ezpp_aim_pp(ezfc),
			SpeedPP:       C.ezpp_speed_pp(ezfc),
			Accuracy:      C.ezpp_accuracy_percent(ezfc),
			N300:          C.ezpp_n300(ezfc),
			N100:          C.ezpp_n100(ezfc),
			N50:           C.ezpp_n50(ezfc),
			NMiss:         C.ezpp_nmiss(ezfc),
			AR:            C.ezpp_ar(ezfc),
			CS:            C.ezpp_cs(ezfc),
			OD:            C.ezpp_od(ezfc),
			HP:            C.ezpp_hp(ezfc),
			Artist:        C.GoString(C.ezpp_artist(ezfc)),
			ArtistUnicode: C.GoString(C.ezpp_artist_unicode(ezfc)),
			Title:         C.GoString(C.ezpp_title(ezfc)),
			TitleUnicode:  C.GoString(C.ezpp_title_unicode(ezfc)),
			Version:       C.GoString(C.ezpp_version(ezfc)),
			Creator:       C.GoString(C.ezpp_creator(ezfc)),
			NCircles:      C.ezpp_ncircles(ezfc),
			NSliders:      C.ezpp_nsliders(ezfc),
			NSpinners:     C.ezpp_nspinners(ezfc),
			ODMS:          C.ezpp_odms(ezfc),
			Mode:          C.ezpp_mode(ezfc),
			Combo:         C.ezpp_combo(ezfc),
			MaxCombo:      C.ezpp_max_combo(ezfc),
			Mods:          C.ezpp_mods(ezfc),
			ScoreVersion:  C.ezpp_score_version(ezfc),
		}

	}

	return nil
}

func GetFCData() {

	for {
		ezfc := C.ezpp_new()
		C.ezpp_set_autocalc(ezfc, 1)
		if memory.DynamicAddresses.IsReady == true && memory.MenuData.OsuStatus == 2 && memory.GameplayData.GameMode == 0 {
			var data PPfc
			readFCData(&data, ezfc)
			if memory.GameplayData.Combo.Max > 0 {
				memory.GameplayData.PP.PPifFC = cast.ToInt32(float64(data.Total))
			}
		}
		C.ezpp_free(ezfc)

		time.Sleep(250 * time.Millisecond)
	}
}
