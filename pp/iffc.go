package pp

//TODO: I need to figure out how to use only one calc.

import (
	"errors"
	"math"
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
import "C"

var ezfc C.ezpp_t

type PPfc struct {
	RestSS C.float
	Acc    C.float
}

func readFCData(data *PPfc, ezfc C.ezpp_t, acc C.float) error {
	path := memory.MenuData.Bm.Path.FullDotOsu

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
		C.ezpp_set_mods(ezfc, C.int(memory.MenuData.Mods.AppliedMods))
		totalObj := C.ezpp_nobjects(ezfc)
		totalCombo := C.ezpp_max_combo(ezfc)

		C.ezpp_set_combo(ezfc, C.int(totalCombo)) //since we are not freeing the counter every time we need to clear the combo
		C.ezpp_set_nmiss(ezfc, C.int(0))

		remaining := int16(totalObj) - memory.GameplayData.Hits.H300 - memory.GameplayData.Hits.H100 - memory.GameplayData.Hits.H50 - memory.GameplayData.Hits.H0
		ifRestSSACC := float64(calculateAccuracy(float32(memory.GameplayData.Hits.H300+remaining), float32(memory.GameplayData.Hits.H100), float32(memory.GameplayData.Hits.H50), float32(memory.GameplayData.Hits.H0)))
		ifRestSSACC = math.Round(ifRestSSACC*100) / 100
		C.ezpp_set_accuracy_percent(ezfc, C.float(ifRestSSACC))
		ifRestSS := C.ezpp_pp(ezfc)
		C.ezpp_set_accuracy_percent(ezfc, C.float(acc))
		//C.ezpp_set_score_version(ezfc)
		*data = PPfc{
			RestSS: ifRestSS,
			Acc:    C.ezpp_pp(ezfc),
		}

		//fmt.Println("True: ", ifRestSS, " MaxThisPlay: ", maxThisPlay, " Current: ", memory.GameplayData.PP.Pp, " PossibleMaxCombo: ", possibleMax)

	}

	return nil
}

func GetFCData() {
	ezfc := C.ezpp_new()
	C.ezpp_set_autocalc(ezfc, 1)
	for {

		if memory.DynamicAddresses.IsReady == true {

			switch memory.GameplayData.GameMode {
			case 0, 1:

				if memory.MenuData.OsuStatus == 2 && memory.GameplayData.Combo.Max > 0 {
					var data PPfc
					readFCData(&data, ezfc, C.float(memory.GameplayData.Accuracy))
					if memory.GameplayData.Combo.Max > 0 {
						memory.GameplayData.PP.PPifFC = cast.ToInt32(float64(data.RestSS))
					}
				}
				switch memory.MenuData.OsuStatus {
				case 1, 4, 5, 13, 2:
					if memory.MenuData.OsuStatus == 2 && memory.MenuData.Bm.Time.PlayTime > 150 { //To catch up with the F2-->Enter
						//C.ezpp_free(ezfc)
						time.Sleep(250 * time.Millisecond)
						continue
					}

					var data PPfc
					readFCData(&data, ezfc, 100.0)
					memory.MenuData.PP.PpSS = cast.ToInt32(float64(data.Acc))
					readFCData(&data, ezfc, 99.0)
					memory.MenuData.PP.Pp99 = cast.ToInt32(float64(data.Acc))
					readFCData(&data, ezfc, 98.0)
					memory.MenuData.PP.Pp98 = cast.ToInt32(float64(data.Acc))
					readFCData(&data, ezfc, 97.0)
					memory.MenuData.PP.Pp97 = cast.ToInt32(float64(data.Acc))
					readFCData(&data, ezfc, 96.0)
					memory.MenuData.PP.Pp96 = cast.ToInt32(float64(data.Acc))
					readFCData(&data, ezfc, 95.0)
					memory.MenuData.PP.Pp95 = cast.ToInt32(float64(data.Acc))
				}
				//C.ezpp_free(ezfc)
			}

		}

		time.Sleep(250 * time.Millisecond)
	}
}
