// +build linux

package mem

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

func FindProcess(re *regexp.Regexp) ([]Process, error) {
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	var pids []int
	for _, dir := range dirs {
		if pid, err := strconv.Atoi(dir.Name()); err == nil {
			pids = append(pids, pid)
		}
	}
	var procs []Process
	for _, pid := range pids {
		path := fmt.Sprintf("/proc/%d/cmdline", pid)
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		content, err := ioutil.ReadAll(f)
		if err != nil {
			continue
		}

		slices := bytes.SplitN(content, []byte{'\x00'}, 2)
		if !re.Match(slices[0]) {
			continue
		}

		procs = append(procs, process{pid})
	}
	if len(procs) < 1 {
		return nil, ErrNoProcess
	}
	return procs, nil
}

func FindWindow(title string) (syscall.Handle, error) {
	return nil, errors.New("Not implemented!")
}

func GetWindowThreadProcessID(hwnd syscall.Handle) int32 {
	return //not implemented
}

type process struct {
	pid int
}

func (p process) ExecutablePath() (string, error) {
	path := fmt.Sprintf("/proc/%d/exe", p.pid)
	path, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

func (p process) Close() error {
	return nil
}

func (p process) Pid() int {
	return p.pid
}

func (p process) ReadAt(b []byte, off int64) (n int, err error) {
	localIov := [1]unix.Iovec{
		{Base: &b[0]},
	}
	localIov[0].SetLen(len(b))
	remoteIov := [1]unix.RemoteIovec{
		{Base: uintptr(off), Len: len(b)},
	}
	n, err = unix.ProcessVMReadv(p.pid, localIov[:], remoteIov[:], 0)
	logRead(b, n, off, err)
	return n, err
}

func (p process) Maps() ([]Map, error) {
	path := fmt.Sprintf("/proc/%d/maps", p.pid)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var maps []Map
	s := bufio.NewScanner(f)
	for s.Scan() {
		var reg region
		_, err := fmt.Sscanf(s.Text(), "%x-%x",
			&reg.start, &reg.end)
		if err != nil && err != io.EOF {
			return nil, err
		}
		maps = append(maps, reg)
	}
	return maps, nil
}

type region struct {
	start int64
	end   int64
}

func (r region) Start() int64 {
	return r.start
}

func (r region) Size() int64 {
	return r.end - r.start
}
