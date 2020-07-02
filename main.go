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
	err := db.InitDB()
	if err != nil {
		log.Println("osu database parse error!, your osu!.db file is either too old or corrupt!", err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
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
	go web.SetupStructure()
	go web.HTTPServer()
	go web.SetupRoutes()
	go pp.GetData()
	go pp.GetFCData()
	http.ListenAndServe("127.0.0.1:8085", nil) //This duplicate fileserver is for backwards compatibility only and will be removed in the future.

}
