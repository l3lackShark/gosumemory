package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/l3lackShark/gosumemory/pp"
	"github.com/l3lackShark/gosumemory/web"

	"github.com/l3lackShark/gosumemory/memory"
)

func main() {
	fmt.Println("gosumemory v0.x-alpha")
	updateTimeFlag := flag.Int("update", 100, "How fast should we update the values? (in milliseconds)")
	songsFolderFlag := flag.String("path", "D:\\osu!\\Songs", `Path to osu! Songs directory ex: /mnt/ps3drive/osu\!/Songs`)
	flag.Parse()
	memory.UpdateTime = *updateTimeFlag
	memory.SongsFolderPath = *songsFolderFlag
	if runtime.GOOS != "windows" && memory.SongsFolderPath == "auto" {
		log.Fatalln("Please specify path to osu!Songs (see --help)")
	}
	go memory.Init()
	go web.SetupStructure()
	go web.HTTPServer()
	go web.SetupRoutes()
	go pp.GetData()
	go pp.GetFCData()
	http.ListenAndServe(":8085", nil)

}
