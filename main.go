package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/l3lackShark/gosumemory/values"

	"github.com/l3lackShark/gosumemory/patterns"
)

func main() {
	fmt.Println("owo")
	updateTimeFlag := flag.Int("update", 100, "How fast should we update the values? (in milliseconds)")
	flag.Parse()
	patterns.UpdateTime = *updateTimeFlag
	go patterns.Init()
	for {
		fmt.Println(values.MenuData.InnerBGPath)
		time.Sleep(500 * time.Millisecond)
	}
	//fmt.Println(patterns.OsuStaticAddresses)

}
