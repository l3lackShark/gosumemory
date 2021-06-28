//+build windows

package pp

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

/*
#include "oppai.c"
*/
import "C"

func wCharPtrToString(p *C.wchar_t) string {
	return windows.UTF16PtrToString((*uint16)(p))
}

func wCharPtrFromString(s string) (*C.wchar_t, error) {
	p, err := windows.UTF16PtrFromString(s)
	return (*C.wchar_t)(p), err
}

func calcpp(data *(C.ezpp_t), path string) error {
	osu, err := wCharPtrFromString(path)
	if err != nil {
		return fmt.Errorf("%s, %e", "UTF16 wchar_t* convert err", err)
	}
	if rc := C.ezpp_win(*data, osu); rc < 0 {
		return errors.New(C.GoString(C.errstr(rc)))
	}
	return nil
}
