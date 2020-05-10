package web

import (
	"encoding/json"
	"log"
	"time"

	"github.com/l3lackShark/gosumemory/memory"
	"github.com/l3lackShark/gosumemory/values"
)

//SetupStructure sets up ws and json output
func SetupStructure() {
	var err error
	type wsStruct struct { //order sets here
		A values.InMenuValues   `json:"menuContainer"`
		B values.GameplayValues `json:"gameplayContainer"`
	}
	for {
		group := wsStruct{
			A: values.MenuData,
			B: values.GameplayData,
		}

		JSONByte, err = json.Marshal(group)
		if err != nil {
			log.Println("error:", err)
		}
		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}

}
