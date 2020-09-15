package pp

//TODO: I need to figure out how to use only one calc.

//#cgo LDFLAGS: -lm
//#cgo CPPFLAGS: -DOPPAI_STATIC_HEADER
//#include <stdlib.h>
//#include "oppai.c"
import "C"
import (
	"errors"
	"math"
	"strings"
	"time"
	"unsafe"

	"github.com/l3lackShark/gosumemory/memory"
	"github.com/spf13/cast"
)

var ezmax C.ezpp_t

type PPmax struct {
	MaxThisPlay C.float
}

func readMaxData(data *PPmax, ezmax C.ezpp_t) error {
	path := memory.MenuData.Bm.Path.FullDotOsu

	if strings.HasSuffix(path, ".osu") && memory.DynamicAddresses.IsReady == true {
		cpath := C.CString(path)

		defer C.free(unsafe.Pointer(cpath))
		if rc := C.ezpp(ezmax, cpath); rc < 0 {
			return errors.New(C.GoString(C.errstr(rc)))
		}
		C.ezpp_set_base_ar(ezmax, C.float(memory.MenuData.Bm.Stats.BeatmapAR))
		C.ezpp_set_base_od(ezmax, C.float(memory.MenuData.Bm.Stats.BeatmapOD))
		C.ezpp_set_base_cs(ezmax, C.float(memory.MenuData.Bm.Stats.BeatmapCS))
		C.ezpp_set_base_hp(ezmax, C.float(memory.MenuData.Bm.Stats.BeatmapHP))
		C.ezpp_set_mods(ezmax, C.int(memory.MenuData.Mods.AppliedMods))
		totalObj := C.ezpp_nobjects(ezmax)
		totalCombo := C.ezpp_max_combo(ezmax)

		remaining := int16(totalObj) - memory.GameplayData.Hits.H300 - memory.GameplayData.Hits.H100 - memory.GameplayData.Hits.H50 - memory.GameplayData.Hits.H0
		ifRestSSACC := float64(calculateAccuracy(float32(memory.GameplayData.Hits.H300+remaining), float32(memory.GameplayData.Hits.H100), float32(memory.GameplayData.Hits.H50), float32(memory.GameplayData.Hits.H0)))
		ifRestSSACC = math.Round(ifRestSSACC*100) / 100
		C.ezpp_set_accuracy_percent(ezmax, C.float(ifRestSSACC))

		//Get Possible max combo in the current play
		var possibleMax float64
		//var lessThanMaxCombo bool
		if memory.GameplayData.Hits.H0+memory.GameplayData.Hits.HSB > 0 {
			possibleMax = math.Max(float64(totalCombo-currMaxCombo), float64(memory.GameplayData.Combo.Max))
			//lessThanMaxCombo = true
		} else {
			possibleMax = float64(totalCombo)
			//lessThanMaxCombo = false
		}

		if memory.MenuData.OsuStatus == 2 {
			C.ezpp_set_nmiss(ezmax, C.int(memory.GameplayData.Hits.H0))
			C.ezpp_set_combo(ezmax, C.int(possibleMax))
		}

		maxThisPlay := C.ezpp_pp(ezmax)
		*data = PPmax{
			MaxThisPlay: maxThisPlay,
		}
		// type test struct {
		// 	expectedAccuracy float64
		// 	RemainingHitObj  int16
		// 	lessThanMaxCombo bool
		// 	maxPossibleCombo int32
		// 	maxThisPlayPP    int32
		// 	currentPP        int32
		// }
		//testing := test{ifRestSSACC, remaining, lessThanMaxCombo, int32(possibleMax), int32(maxThisPlay), memory.GameplayData.PP.Pp}
		//pp.Println(testing)

	}

	return nil
}

func GetMaxData() {
	ezmax := C.ezpp_new()
	C.ezpp_set_autocalc(ezmax, 1)
	for {

		if memory.DynamicAddresses.IsReady == true {

			switch memory.GameplayData.GameMode {
			case 0, 1:

				if memory.MenuData.OsuStatus == 2 && memory.GameplayData.Combo.Max > 0 {
					var data PPmax
					readMaxData(&data, ezmax)
					memory.GameplayData.PP.PPMaxThisPlay = cast.ToInt32(float64(data.MaxThisPlay))
				}
			}
		}

		time.Sleep(250 * time.Millisecond)
	}
}
