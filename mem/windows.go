// +build windows

package mem

import (
	"fmt"
	"regexp"
	"strings"
	"syscall"
	"unsafe"

	windows "github.com/elastic/go-windows"
	xsyscall "golang.org/x/sys/windows"
)

var (
	modkernel32                = xsyscall.NewLazySystemDLL("kernel32.dll")
	user32                     = xsyscall.NewLazySystemDLL("user32.dll")
	procEnumWindows            = user32.NewProc("EnumWindows")
	procGetWindowTextW         = user32.NewProc("GetWindowTextW")
	getWindowThreadProcessID   = user32.NewProc("GetWindowThreadProcessId")
	procVirtualQueryEx         = modkernel32.NewProc("VirtualQueryEx")
	queryFullProcessImageNameW = modkernel32.NewProc("QueryFullProcessImageNameW")
)

func enumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowThreadProcessID(hwnd syscall.Handle) int32 {
	var processID int32
	getWindowThreadProcessID.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&processID)))
	return processID
}

func FindWindow(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := getWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if strings.Contains(syscall.UTF16ToString(b), title) {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	enumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with title '%s' found", title)
	}
	return hwnd, nil
}

func virtualQueryEx(handle syscall.Handle, off int64) (region, error) {
	var reg region
	r1, _, e1 := syscall.Syscall6(
		procVirtualQueryEx.Addr(),
		4,
		uintptr(handle),
		uintptr(off),
		uintptr(unsafe.Pointer(&reg)),
		uintptr(unsafe.Sizeof(reg)),
		0, 0,
	)
	if r1 == 0 {
		if e1 != 0 {
			return region{}, e1
		} else {
			return region{}, syscall.EINVAL
		}
	}
	return reg, nil
}

func queryFullProcessImageName(hProcess syscall.Handle) (string, error) {
	var buf [syscall.MAX_PATH]uint16
	n := uint32(len(buf))
	r1, _, e1 := queryFullProcessImageNameW.Call(
		uintptr(hProcess),
		uintptr(0),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&n)))
	if r1 == 0 {
		if e1 != nil {
			return "", e1
		} else {
			return "", syscall.EINVAL
		}
	}
	return syscall.UTF16ToString(buf[:n]), nil

}

func FindProcess(re *regexp.Regexp) ([]Process, error) {
	var procs []Process
	pids, err := windows.EnumProcesses()
	if err != nil {
		return nil, err
	}
	for _, pid := range pids {
		handle, err := syscall.OpenProcess(
			syscall.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ,
			false, pid)
		if err != nil {
			continue
		}
		name, err := windows.GetProcessImageFileName(handle)
		if err != nil {
			syscall.CloseHandle(handle)
			continue
		}
		if re.MatchString(name) {
			procs = append(procs, process{pid, handle})
		}
	}
	if len(procs) < 1 {
		return nil, ErrNoProcess
	}
	return procs, nil
}

type process struct {
	pid uint32
	h   syscall.Handle
}

func (p process) HandleFromTitle() (string, error) {
	return queryFullProcessImageName(p.h)
}

func (p process) ExecutablePath() (string, error) {
	return queryFullProcessImageName(p.h)
}

func (p process) Close() error {
	return syscall.CloseHandle(p.h)
}

func (p process) Pid() int {
	return int(p.pid)
}

func (p process) ReadAt(b []byte, off int64) (n int, err error) {
	un, err := windows.ReadProcessMemory(p.h, uintptr(off), b)
	logRead(b, int(un), off, err)
	return int(un), err
}

func (p process) Maps() ([]Map, error) {
	lastAddr := int64(0)
	var maps []Map
	for {
		reg, err := virtualQueryEx(p.h, lastAddr)
		if err != nil {
			if lastAddr == 0 {
				return nil, err
			}
			break
		}
		maps = append(maps, reg)
		lastAddr = reg.Start() + reg.Size()
	}
	return maps, nil
}

type region struct {
	baseAddress       uintptr
	allocationBase    uintptr
	allocationProtect int32
	regionSize        int
	state             int32
	protect           int32
	type_             int32
}

func (r region) Start() int64 {
	return int64(r.baseAddress)
}

func (r region) Size() int64 {
	return int64(r.regionSize)
}
