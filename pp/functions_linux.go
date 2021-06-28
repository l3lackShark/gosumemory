//+build linux

package pp

/*
#include "oppai.c"
*/
import "C"

func calcpp(data *(C.ezpp_t), path string) error {
	cpath := C.CString(path)

	defer C.free(unsafe.Pointer(cpath))
	if rc := C.ezpp(ez, cpath); rc < 0 {
		return errors.New(C.GoString(C.errstr(rc)))
	}
	return nil
}
