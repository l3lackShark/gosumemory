package patterns

//Patterns is Base osu signatures stuct
type Patterns struct {
	status        string
	bpm           string
	base          string
	inMenuMods    string
	playTime      string
	playContainer string
}

//StaticAddresses (should be updated every client restart)
type StaticAddresses struct {
	Status        uint32
	BPM           uint32
	Base          uint32
	InMenuMods    uint32
	PlayTime      uint32
	PlayContainer uint32
}

//OsuSignatures are the main sigs used by the program
var osuSignatures = Patterns{
	status:        "48 83 F8 04 73 1E",             //-0x4
	bpm:           "8B 40 08 89 86 4C 01 00 00 C6", //-0x4
	base:          "F8 01 74 04 83 65",
	inMenuMods:    "C8 FF ?? ?? ?? ?? ?? 81 0D ?? ?? ?? ?? 00 08 00 00",
	playTime:      "5E 5F 5D C3 A1 ?? ?? ?? ?? 89 ?? 04",
	playContainer: "85 C9 74 1F 8D 55 F0 8B 01",
}
