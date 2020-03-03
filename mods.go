package main

import (
	"strings"
)

// Mods represents zero or more mods of an osu! score.
type Mods int

// Mods constants.
// Names are taken from official osu! documentation.
const (
	ModsNoFail Mods = 1 << iota
	ModsEasy
	ModsTouchDevice
	ModsHidden
	ModsHardRock
	ModsSuddenDeath
	ModsDoubleTime
	ModsRelax
	ModsHalfTime
	ModsNightcore
	ModsFlashlight
	ModsAutoplay
	ModsSpunOut
	ModsRelax2 // Autopilot
	ModsPerfect
	ModsKey4
	ModsKey5
	ModsKey6
	ModsKey7
	ModsKey8
	ModsFadeIn
	ModsRandom
	ModsCinema
	ModsTargetPractice
	ModsKey9
	ModsKeyCoop
	ModsKey1
	ModsKey3
	ModsKey2
	ModsScoreV2

	ModsKeyMod         Mods = ModsKey1 | ModsKey2 | ModsKey3 | ModsKey4 | ModsKey5 | ModsKey6 | ModsKey7 | ModsKey8 | ModsKey9
	ModsFreeModAllowed Mods = ModsNoFail | ModsEasy | ModsHidden | ModsHardRock | ModsSuddenDeath | ModsFlashlight | ModsFadeIn | ModsRelax | ModsSpunOut | ModsKeyMod
	ModsScoreIncrease  Mods = ModsHidden | ModsHardRock | ModsDoubleTime | ModsFlashlight | ModsFadeIn

//	modsSpeedChanging Mods = ModsDT | ModsHT | ModsNC
//	modsMapChanging   Mods = ModsHR | ModsEZ | modsSpeedChanging
)

// Convenience aliases.
const (
	NF Mods = ModsNoFail
	EZ Mods = ModsEasy
	TD Mods = ModsTouchDevice
	HD Mods = ModsHidden
	HR Mods = ModsHardRock
	SD Mods = ModsSuddenDeath
	DT Mods = ModsDoubleTime
	RX Mods = ModsRelax
	HT Mods = ModsHalfTime
	NC Mods = ModsNightcore
	FL Mods = ModsFlashlight
	SO Mods = ModsSpunOut
	PF Mods = ModsPerfect

	V2 Mods = ModsScoreV2
	AT Mods = ModsAutoplay
	AP Mods = ModsRelax2
	FI Mods = ModsFadeIn
	RD Mods = ModsRandom
	CN Mods = ModsCinema
	TP Mods = ModsTargetPractice
	CO Mods = ModsKeyCoop

	K1 Mods = ModsKey1
	K2 Mods = ModsKey2
	K3 Mods = ModsKey3
	K4 Mods = ModsKey4
	K5 Mods = ModsKey5
	K6 Mods = ModsKey6
	K7 Mods = ModsKey7
	K8 Mods = ModsKey8
	K9 Mods = ModsKey9
)

// This is a slice instead of a map because order matters
var modStrings = []struct {
	mod Mods
	str string
}{
	{NF, "NF"},
	{EZ, "EZ"},
	{TD, "TD"},
	{HD, "HD"},
	{HR, "HR"},
	{SD, "SD"},
	{DT, "DT"},
	{RX, "RX"},
	{HT, "HT"},
	{NC, "NC"},
	{FL, "FL"},
	{AT, "AT"},
	{SO, "SO"},
	{AP, "AP"},
	{PF, "PF"},
	{K4, "4K"},
	{K5, "5K"},
	{K6, "6K"},
	{K7, "7K"},
	{K8, "8K"},
	{FI, "FI"},
	{RD, "RD"},
	{CN, "CN"},
	{TP, "TP"},
	{K9, "9K"},
	{CO, "CO"},
	{K1, "1K"},
	{K3, "3K"},
	{K2, "2K"},
	{V2, "V2"},
}

func (m Mods) String() string {
	if m == 0 {
		return ""
	}

	var s strings.Builder

	for _, v := range modStrings {
		if m&v.mod > 0 {
			s.WriteString(v.str)
		}
	}

	return s.String()
}
