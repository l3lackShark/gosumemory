package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/l3lackShark/gosumemory/db"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/l3lackShark/gosumemory/pp"
	"github.com/l3lackShark/gosumemory/updater"
	"github.com/l3lackShark/gosumemory/web"
)

func main() {
	updateTimeFlag := flag.Int("update", 100, "How fast should we update the values? (in milliseconds)")
	shouldWeUpdate := flag.Bool("autoupdate", true, "Should we auto update the application?")
	isRunningInWINE := flag.Bool("wine", false, "Running under WINE?")
	songsFolderFlag := flag.String("path", "auto", `Path to osu! Songs directory ex: /mnt/ps3drive/osu\!/Songs`)
	flag.Parse()
	memory.UpdateTime = *updateTimeFlag
	memory.SongsFolderPath = *songsFolderFlag
	memory.UnderWine = *isRunningInWINE
	if runtime.GOOS != "windows" && memory.SongsFolderPath == "auto" {
		log.Fatalln("Please specify path to osu!Songs (see --help)")
	}
	if *shouldWeUpdate == true {
		updater.DoSelfUpdate()
	}

	go memory.Init()
	go db.InitOSRWriting()
	err := db.InitDB()
	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	go web.SetupStructure()
	go web.HTTPServer()
	go web.SetupRoutes()
	go pp.GetData()
	go pp.GetFCData()
	fmt.Println("WARNING: Mania pp calcualtion is experimental and only works if you choose mania gamemode in the SongSelect!")
	fmt.Println("Initialization complete, you can now visit http://localhost:24050 or add it as a browser source in OBS")
	pp.GetEditorData()

}
