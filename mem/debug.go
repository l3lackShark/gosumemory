package mem

import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
)

var Debug = false
var indent int32 = 0

func beginDebug() {
	if !Debug {
		return
	}

	var callers [8]uintptr
	n := runtime.Callers(2, callers[:])

	frames := runtime.CallersFrames(callers[:n])
	var stack []string
	for {
		frame, more := frames.Next()
		if strings.HasPrefix(frame.Function, "runtime") {
			break
		}
		stack = append(stack,
			fmt.Sprintf("%s:%d", frame.Function, frame.Line))
		if !more {
			break
		}
	}
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}
	for _, s := range stack {
		log("%s\n", s)
		atomic.AddInt32(&indent, 1)
	}
}

func pushDebug() int32 {
	if !Debug {
		return 0
	}

	return atomic.AddInt32(&indent, 1) - 1
}

func popDebug(i int32) {
	if !Debug {
		return
	}

	atomic.StoreInt32(&indent, i)
}

func endDebug() {
	if !Debug {
		return
	}

	atomic.StoreInt32(&indent, 0)
}

func log(format string, args ...interface{}) {
	if !Debug {
		return
	}

	for i := int32(0); i < atomic.LoadInt32(&indent)*4; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf(format, args...)
}

func logRead(b []byte, n int, off int64, err error) {
	if !Debug {
		return
	}

	if err == nil {
		var arr string
		if n < 16 {
			arr = fmt.Sprintf("%v", b[:n])
		} else {
			arr = "[...]"
		}
		log("Read(0x%x, %d): %s\n", uint64(off), n, arr)
	} else {
		log("Read(0x%x, %d): %v\n", uint64(off), n, err)
	}
}
