package memory

//NewPatterns is Base osu signatures stuct (for tdeo's wrapper)
type NewPatterns struct {
	Base          int64 `sig:"F8 01 74 04 83 65"`                                  //-0xC
	InMenuMods    int64 `sig:"C8 FF ?? ?? ?? ?? ?? 81 0D ?? ?? ?? ?? 00 08 00 00"` //+0x9
	PlayTime      int64 `sig:"5E 5F 5D C3 A1 ?? ?? ?? ?? 89 ?? 04"`                //+0x5
	PlayContainer int64 `sig:"85 C9 74 1F 8D 55 F0 8B 01"`
	LeaderBoard   int64 `sig:"A1 ?? ?? ?? ?? 8B 50 04 8B 0D"` //+0x1
	SongsFolder   int64 `sig:"?? ?? 67 ?? 2F 00 28 00"`
	ChatChecker   int64 `sig:"0A D7 23 3C 00 00 ?? 01"` //-0x20 (value)
}

//StaticAddresses (should be updated every client restart)
type StaticAddresses struct {
	Status        uint32
	BPM           uint32
	Base          uint32
	InMenuMods    uint32
	PlayTime      uint32
	PlayContainer uint32
	LeaderBoard   uint32
	SongsFolder   uint32
	ChatChecker   uint32
}

//osuStaticAddresses (should be updated every client restart)
var osuStaticAddresses = StaticAddresses{}
