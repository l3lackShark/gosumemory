package web

import (
	"encoding/json"
	"log"
	"time"

	"github.com/l3lackShark/gosumemory/memory"
)

//SetupStructure sets up ws and json output
func SetupStructure() {
	var err error
	type wsStruct struct { //order sets here
		A memory.InMenuValues   `json:"menu"`
		B memory.GameplayValues `json:"gameplay"`
	}
	for {
		group := wsStruct{
			A: memory.MenuData,
			B: memory.GameplayData,
		}

		JSONByte, err = json.Marshal(group)
		if err != nil {
			log.Println("error:", err)
		}
		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}

}
