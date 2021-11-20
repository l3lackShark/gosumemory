package helpers

import "time"

func Sleep(val int) {
	time.Sleep(time.Duration(val) * time.Millisecond)
}
