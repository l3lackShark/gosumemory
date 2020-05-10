package pp

import (
	"log"
	"strings"
	"time"
	"unsafe"

	"github.com/l3lackShark/gosumemory/memory"
	"github.com/l3lackShark/gosumemory/values"
)

//#cgo LDFLAGS:/usr/local/lib/liboppai.so
//#include "oppai.c"
//#include <stdlib.h>
import "C"

//GetData resolves pp values (uses C oppai binary for now, in hopes to switch to cgo)
func GetData() {
	for {
		ez := C.ezpp_new()
		C.ezpp_set_autocalc(ez, 1)
		path := (memory.SongsFolderPath + "/" + values.MenuData.BeatmapFolderString + "/" + values.MenuData.BeatmapOsuFileString) //TODO: Automatic Songs folder finder
		var tempPath string
		if strings.HasSuffix(path, ".osu") && path != tempPath && values.MenuData.IsReady == true {
			tempPath = path
			cpath := C.CString(path)
			defer C.free(unsafe.Pointer(cpath))
			if rc := C.ezpp(ez, cpath); rc < 0 {
				log.Println((C.GoString(C.errstr(rc))))
			}
			switch values.MenuData.OsuStatus {
			case 2:
				C.ezpp_set_accuracy_percent(ez, C.float(values.GameplayData.Accuracy))
				C.ezpp_set_end_time(ez, C.float(values.MenuData.PlayTime))
				C.ezpp_set_nmiss(ez, C.int(values.GameplayData.HitMiss))
				C.ezpp_set_combo(ez, C.int(values.GameplayData.MaxCombo))
				C.ezpp_set_mods(ez, C.int(values.GameplayData.AppliedMods))
				values.GameplayData.Pp = float64(C.ezpp_pp(ez))

			default:
				C.ezpp_set_base_ar(ez, C.float(values.MenuData.BeatmapAR))
				C.ezpp_set_base_od(ez, C.float(values.MenuData.BeatmapOD))
				C.ezpp_set_base_cs(ez, C.float(values.MenuData.BeatmapCS))
				C.ezpp_set_base_hp(ez, C.float(values.MenuData.BeatmapHP))
				values.MenuData.PpSS = float64(C.ezpp_pp(ez))
				C.ezpp_set_accuracy_percent(ez, C.float(99.0))
				values.MenuData.Pp99 = float64(C.ezpp_pp(ez))
				C.ezpp_set_accuracy_percent(ez, C.float(98.0))
				values.MenuData.Pp98 = float64(C.ezpp_pp(ez))
				C.ezpp_set_accuracy_percent(ez, C.float(97.0))
				values.MenuData.Pp97 = float64(C.ezpp_pp(ez))
				C.ezpp_set_accuracy_percent(ez, C.float(96.0))
				values.MenuData.Pp96 = float64(C.ezpp_pp(ez))
				C.ezpp_set_accuracy_percent(ez, C.float(95.0))
				values.MenuData.Pp95 = float64(C.ezpp_pp(ez))

			}
			C.ezpp_free(ez)
		}
		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}

}
