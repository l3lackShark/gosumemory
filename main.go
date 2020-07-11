package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/l3lackShark/gosumemory/db"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/l3lackShark/gosumemory/pp"
	"github.com/l3lackShark/gosumemory/web"
)

func main() {
	fmt.Println("gosumemory v0.x-alpha")
	updateTimeFlag := flag.Int("update", 100, "How fast should we update the values? (in milliseconds)")
	songsFolderFlag := flag.String("path", "auto", `Path to osu! Songs directory ex: /mnt/ps3drive/osu\!/Songs`)
	flag.Parse()
	memory.UpdateTime = *updateTimeFlag
	memory.SongsFolderPath = *songsFolderFlag
	if runtime.GOOS != "windows" && memory.SongsFolderPath == "auto" {
		log.Fatalln("Please specify path to osu!Songs (see --help)")
	}
	go memory.Init()
	err := db.InitDB()
	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	fmt.Println("WARNING: Mania pp calcualtion is experimental and only works if you choose mania gamemode in the SongSelect!")

	go web.SetupStructure()
	go web.HTTPServer()
	go web.SetupRoutes()
	go pp.GetData()
	go pp.GetFCData()
	go pp.GetEditorData()
	http.ListenAndServe("127.0.0.1:8085", nil) //This duplicate fileserver is for backwards compatibility only and will be removed in the future.

}
