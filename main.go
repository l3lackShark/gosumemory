package main

import (
	"os"

	"github.com/k0kubun/pp"

	"github.com/l3lackShark/gosumemory/opengl"
)

// func main() {
// 	fmt.Println("gosumemory v0.x-alpha")
// 	updateTimeFlag := flag.Int("update", 100, "How fast should we update the values? (in milliseconds)")
// 	songsFolderFlag := flag.String("path", "D:\\osu!\\Songs", `Path to osu! Songs directory ex: /mnt/ps3drive/osu\!/Songs`)
// 	flag.Parse()
// 	memory.UpdateTime = *updateTimeFlag
// 	memory.SongsFolderPath = *songsFolderFlag
// 	if runtime.GOOS != "windows" && memory.SongsFolderPath == "auto" {
// 		log.Fatalln("Please specify path to osu!Songs (see --help)")
// 	}
// 	go memory.Init()
// 	go web.SetupStructure()
// 	go web.HTTPServer()
// 	go web.SetupRoutes()
// 	go pp.GetData()
// 	http.ListenAndServe(":8085", nil)

// }
func main() {
	if err := opengl.Init(); err != nil {
		pp.Println(err)
		os.Exit(1)
	}

}
