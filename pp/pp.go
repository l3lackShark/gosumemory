package pp

import (
	"log"
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

//GetData resolves pp values (using cgo)
func GetData() {
	for {
		ez := C.ezpp_new()
		C.ezpp_set_autocalc(ez, 1)
		path := (memory.SongsFolderPath + "/" + memory.MenuData.Bm.Path.BeatmapFolderString + "/" + memory.MenuData.Bm.Path.BeatmapOsuFileString) //TODO: Automatic Songs folder finder
		var tempPath string
		if strings.HasSuffix(path, ".osu") && path != tempPath && memory.DynamicAddresses.IsReady == true {
			tempPath = path
			cpath := C.CString(path)
			defer C.free(unsafe.Pointer(cpath))
			if rc := C.ezpp(ez, cpath); rc < 0 {
				log.Println((C.GoString(C.errstr(rc))))
			}
			switch memory.MenuData.OsuStatus {
			case 2:
				C.ezpp_set_accuracy_percent(ez, C.float(memory.GameplayData.Accuracy))
				C.ezpp_set_end_time(ez, C.float(memory.MenuData.Bm.Time.PlayTime))
				C.ezpp_set_nmiss(ez, C.int(memory.GameplayData.Hits.H0))
				C.ezpp_set_combo(ez, C.int(memory.GameplayData.Combo.Max))
				C.ezpp_set_mods(ez, C.int(memory.GameplayData.Mods.AppliedMods))
				memory.GameplayData.PP.Pp = cast.ToString(float64(C.ezpp_pp(ez)))

			default:
				memory.MenuData.Bm.Metadata.Artist = C.GoString(C.ezpp_artist(ez))
				memory.MenuData.Bm.Metadata.Title = C.GoString(C.ezpp_title(ez))
				memory.MenuData.Bm.Metadata.Version = C.GoString(C.ezpp_version(ez))
				memory.MenuData.Bm.Metadata.Mapper = C.GoString(C.ezpp_creator(ez))
				C.ezpp_set_base_ar(ez, C.float(memory.MenuData.Bm.Stats.BeatmapAR))
				C.ezpp_set_base_od(ez, C.float(memory.MenuData.Bm.Stats.BeatmapOD))
				C.ezpp_set_base_cs(ez, C.float(memory.MenuData.Bm.Stats.BeatmapCS))
				C.ezpp_set_base_hp(ez, C.float(memory.MenuData.Bm.Stats.BeatmapHP))
				memory.MenuData.PP.PpSS = cast.ToString(float64(C.ezpp_pp(ez)))
				C.ezpp_set_accuracy_percent(ez, C.float(99.0))
				memory.MenuData.PP.Pp99 = cast.ToString(float64(C.ezpp_pp(ez)))
				C.ezpp_set_accuracy_percent(ez, C.float(98.0))
				memory.MenuData.PP.Pp98 = cast.ToString(float64(C.ezpp_pp(ez)))
				C.ezpp_set_accuracy_percent(ez, C.float(97.0))
				memory.MenuData.PP.Pp97 = cast.ToString(float64(C.ezpp_pp(ez)))
				C.ezpp_set_accuracy_percent(ez, C.float(96.0))
				memory.MenuData.PP.Pp96 = cast.ToString(float64(C.ezpp_pp(ez)))
				C.ezpp_set_accuracy_percent(ez, C.float(95.0))
				memory.MenuData.PP.Pp95 = cast.ToString(float64(C.ezpp_pp(ez)))

			}

		}
		C.ezpp_free(ez)
		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}

}
