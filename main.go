package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/l3lackShark/gosumemory/pp"
	"github.com/l3lackShark/gosumemory/web"

	"github.com/l3lackShark/gosumemory/memory"
)

func main() {
	fmt.Println("owo")
	updateTimeFlag := flag.Int("update", 100, "How fast should we update the values? (in milliseconds)")
	flag.Parse()
	memory.UpdateTime = *updateTimeFlag
	go memory.Init()
	go web.SetupStructure()
	go web.HTTPServer()
	go web.SetupRoutes()
	go pp.GetData()
	http.ListenAndServe(":8085", nil)

}
