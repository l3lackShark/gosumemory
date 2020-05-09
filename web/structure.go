package web

import (
	"encoding/json"
	"fmt"

	"github.com/l3lackShark/gosumemory/values"
)

//SetupStructure sets up ws and json output
func SetupStructure() {
	var err error
	type wsStruct struct { //order sets here
		A values.InMenuValues   `json:"menuContainer"`
		B values.GameplayValues `json:"gameplayContainer"`
	}

	group := wsStruct{
		A: values.MenuData,
		B: values.GameplayData,
	}
	JSONByte, err = json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
}
