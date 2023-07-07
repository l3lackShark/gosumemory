package deafen

import (
	"fmt"
	"github.com/l3lackShark/gosumemory/config"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/micmonay/keybd_event"
	"github.com/spf13/cast"
	"strings"
)

var currentlyDeafened bool

var keymap = map[string]int{
	"Q": 16,
	"W": 17,
	"E": 18,
	"R": 19,
	"T": 20,
	"Y": 21,
	"U": 22,
	"I": 23,
	"O": 24,
	"P": 25,
	"A": 30,
	"S": 31,
	"D": 32,
	"F": 33,
	"G": 34,
	"H": 35,
	"J": 36,
	"K": 37,
	"L": 38,
	"Z": 44,
	"X": 45,
	"C": 46,
	"V": 47,
	"B": 48,
	"N": 49,
	"M": 50,
}

func AutoDeafen() {
	percentageToDeafen := cast.ToInt32(config.Config["percentageForDeafen"])
	ppToDeafen := cast.ToInt32(config.Config["ppForDeafen"])

	fmt.Println(fmt.Sprintf("[AUTODEAFEN] Using ALT+%s to deafen at either %d percent or %d pp (whatever comes first), can be changed in config.ini", strings.ToUpper(config.Config["deafenKey"]), percentageToDeafen, ppToDeafen))

	for { // this is probably very spaghetti-esque code but i dont know how else to write this :p
		if memory.MenuData.OsuStatus == 2 {
			if memory.MenuData.Bm.Time.FullTime != 0 {
				var percentageOfMap = (cast.ToFloat32(memory.MenuData.Bm.Time.PlayTime) / cast.ToFloat32(memory.MenuData.Bm.Time.FullTime)) * 100
				if percentageOfMap > cast.ToFloat32(percentageToDeafen) && (cast.ToInt32(memory.GameplayData.PP.Pp) > 1) {
					deafen(true)
				}
			}
			if memory.GameplayData.Score == 0 {
				undeafen(true)
			}
		} else {
			undeafen(true)
		}
	}
}

func deafen(intent bool) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	kb.SetKeys(keymap[strings.ToUpper(config.Config["deafenKey"])])
	kb.HasALT(true)

	if currentlyDeafened {
		return
	} else {
		if intent {
			err = kb.Launching()
			if err != nil {
				panic(err)
			}
			currentlyDeafened = true
		}
	}
}

func undeafen(intent bool) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	kb.SetKeys(keymap[strings.ToUpper(config.Config["deafenKey"])])
	kb.HasALT(true)

	if !currentlyDeafened {
		return
	} else {
		if intent {
			err = kb.Launching()
			if err != nil {
				panic(err)
			}
			currentlyDeafened = false
		}
	}
}
