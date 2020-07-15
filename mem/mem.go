package mem

import (
	"errors"
	"io"
)

var (
	ErrNoProcess       = errors.New("no process matching the criteria was found")
	ErrPatternNotFound = errors.New("no memory matched the pattern")
)

type (
	Process interface {
		io.Closer
		io.ReaderAt
		Pid() int
		Maps() ([]Map, error)
	}

	Map interface {
		Start() int64
		Size() int64
	}
)
