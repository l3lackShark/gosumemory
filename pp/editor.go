package pp

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
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

var ezeditor C.ezpp_t

func readEditorData(data *PP, ezeditor C.ezpp_t, needStrain bool) error {
	path := memory.MenuData.Bm.Path.FullDotOsu

	if strings.HasSuffix(path, ".osu") && memory.DynamicAddresses.IsReady == true {
		cpath := C.CString(path)

		defer C.free(unsafe.Pointer(cpath))
		if rc := C.ezpp(ezeditor, cpath); rc < 0 {
			return errors.New(C.GoString(C.errstr(rc)))
		}
		C.ezpp_set_base_ar(ezeditor, C.float(memory.MenuData.Bm.Stats.BeatmapAR))
		C.ezpp_set_base_od(ezeditor, C.float(memory.MenuData.Bm.Stats.BeatmapOD))
		C.ezpp_set_base_cs(ezeditor, C.float(memory.MenuData.Bm.Stats.BeatmapCS))
		C.ezpp_set_base_hp(ezeditor, C.float(memory.MenuData.Bm.Stats.BeatmapHP))
		C.ezpp_set_accuracy_percent(ezeditor, C.float(100.0))
		C.ezpp_set_nmiss(ezeditor, C.int(0))
		if needStrain == true {
			C.ezpp_set_end_time(ezeditor, 0)
			C.ezpp_set_combo(ezeditor, 0)
			C.ezpp_set_nmiss(ezeditor, 0)
			strainArray = nil
			seek := 0
			var window []float64
			var total []float64
			// for seek < int(C.ezpp_time_at(ezeditor, C.ezpp_nobjects(ezeditor)-1)) { //len-1
			for int32(seek) < memory.MenuData.Bm.Time.Mp3Time {
				for obj := 0; obj <= int(C.ezpp_nobjects(ezeditor)-1); obj++ {
					if tempBeatmapFile != memory.MenuData.Bm.Path.BeatmapOsuFileString {
						return nil //Interrupt calcualtion if user has changed the map.
					}
					if int(C.ezpp_time_at(ezeditor, C.int(obj))) >= seek && int(C.ezpp_time_at(ezeditor, C.int(obj))) <= seek+3000 {
						window = append(window, float64(C.ezpp_strain_at(ezeditor, C.int(obj), 0))+float64(C.ezpp_strain_at(ezeditor, C.int(obj), 1)))
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
			memory.MenuData.Bm.Time.FirstObj = int32(C.ezpp_time_at(ezeditor, 0))
			memory.MenuData.Bm.Time.FullTime = int32(C.ezpp_time_at(ezeditor, C.ezpp_nobjects(ezeditor)-1))

		} else {
			C.ezpp_set_mods(ezeditor, C.int(0))
			C.ezpp_set_end_time(ezeditor, C.float(memory.MenuData.Bm.Time.PlayTime))
			C.ezpp_set_combo(ezeditor, C.int(-1))
		}

		*data = PP{
			Total:      C.ezpp_pp(ezeditor),
			Strain:     strainArray,
			AR:         C.ezpp_ar(ezeditor),
			CS:         C.ezpp_cs(ezeditor),
			OD:         C.ezpp_od(ezeditor),
			HP:         C.ezpp_hp(ezeditor),
			StarRating: C.ezpp_stars(ezeditor),
		}
	}
	return nil
}

var tempOsuMD5 string

func GetEditorData() {

	ezeditor := C.ezpp_new()
	C.ezpp_set_autocalc(ezeditor, 1)
	//defer C.ezpp_free(ezeditor)
	var data PP

	for {
		if memory.DynamicAddresses.IsReady == true {
			if memory.MenuData.GameMode == 0 && memory.MenuData.OsuStatus == 1 {
				err := readEditorData(&data, ezeditor, false)
				if err != nil {
					fmt.Println(err)
				}
				memory.GameplayData.PP.Pp = int32(data.Total)
				memory.MenuData.Bm.Stats.BeatmapAR = float32(data.AR)
				memory.MenuData.Bm.Stats.BeatmapCS = float32(data.CS)
				memory.MenuData.Bm.Stats.BeatmapOD = float32(data.OD)
				memory.MenuData.Bm.Stats.BeatmapHP = float32(data.HP)
				memory.MenuData.Bm.Stats.BeatmapSR = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.StarRating)))

				md5, err := hashFileMD5(filepath.Join(memory.SongsFolderPath, memory.MenuData.Bm.Path.BeatmapFolderString, memory.MenuData.Bm.Path.BeatmapOsuFileString))
				if err != nil {
					continue
				}
				if tempOsuMD5 != md5 {
					tempOsuMD5 = md5
					readEditorData(&data, ezeditor, true)
					memory.MenuData.PP.PpStrains = data.Strain
				}

			}

		}
		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}

}
