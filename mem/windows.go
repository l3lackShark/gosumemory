// +build windows

package mem

import (
	"regexp"
	"syscall"
	"unsafe"

	windows "github.com/elastic/go-windows"
	xsyscall "golang.org/x/sys/windows"
)

var (
	modkernel32        = xsyscall.NewLazySystemDLL("kernel32.dll")
	procVirtualQueryEx = modkernel32.NewProc("VirtualQueryEx")
)

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

func FindProcess(re *regexp.Regexp) (Process, error) {
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
			return process{pid, handle}, nil
		}
	}
	return process{}, ErrNoProcess
}

type process struct {
	pid uint32
	h   syscall.Handle
}

func (p process) Close() error {
	return syscall.CloseHandle(p.h)
}

func (p process) Pid() int {
	return int(p.pid)
}

func (p process) ReadAt(b []byte, off int64) (n int, err error) {
	un, err := windows.ReadProcessMemory(p.h, uintptr(off), b)
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
