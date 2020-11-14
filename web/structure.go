package web

import (
	"encoding/json"
	"time"

	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/memory"
)

//SetupStructure sets up ws and json output
func SetupStructure() {
	var err error
	type wsStruct struct { //order sets here
		A memory.InSettingsValues    `json:"settings"`
		B memory.InMenuValues        `json:"menu"`
		C memory.GameplayValues      `json:"gameplay"`
		D memory.ResultsScreenValues `json:"resultsScreen"`
		E memory.TourneyValues       `json:"tourney"`
	}
	for {
		group := wsStruct{
			A: memory.SettingsData,
			B: memory.MenuData,
			C: memory.GameplayData,
			D: memory.ResultsScreenData,
			E: memory.TourneyData,
		}

		JSONByte, err = json.Marshal(group)
		if err != nil {
			pp.Println("JSON Marshall error: ", err, group)
		}
		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}

}
