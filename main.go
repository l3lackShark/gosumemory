package main

import (
	"fmt"
	"time"

	"github.com/l3lackShark/gosumemory/values"

	"github.com/l3lackShark/gosumemory/patterns"
)

func main() {
	fmt.Println("owo")
	// for i := 0; i < 500; i++ {
	// 	fmt.Println(patterns.ResolveOsuStatus())
	// }
	go patterns.Init()
	for {
		fmt.Println(values.OsuData.OsuStatus)
		fmt.Println(values.OsuData.BeatMapID)
		time.Sleep(500 * time.Millisecond)
	}
	//fmt.Println(patterns.OsuStaticAddresses)

}
