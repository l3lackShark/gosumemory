package memory

//InMenuValues inside osu!memory
type InMenuValues struct {
	OsuStatus uint32 `json:"state"`
	Bm        bm     `json:"bm"`
	Mods      modsM  `json:"mods"`
	PP        ppM    `json:"pp"`
}

//GameplayValues inside osu!memory
type GameplayValues struct {
	GameMode int32   `json:"gameMode"`
	Score    int32   `json:"score"`
	Accuracy float64 `json:"accuracy"`
	Combo    combo   `json:"combo"`
	Hp       hp      `json:"hp"`
	Hits     hits    `json:"hits"`
	Mods     modsG   `json:"mods"`
	PP       ppG     `json:"pp"`
}

type bm struct {
	Time           tim      `json:"time"`
	BeatmapID      uint32   `json:"id"`
	BeatmapSetID   uint32   `json:"set"`
	Metadata       Metadata `json:"metadata"`
	Stats          stats    `json:"stats"`
	Path           path     `json:"path"`
	HitObjectStats string   `json:"bmStats"`
	BeatmapString  string   `json:"bmInfo"`
}

type tim struct {
	PlayTime uint32 `json:"current"`
	//FullTime uint32 `json:"full"`
}

// Metadata Map data
type Metadata struct {
	Artist  string `json:"artist"`
	Title   string `json:"title"`
	Mapper  string `json:"mapper"`
	Version string `json:"difficulty"`
}

type stats struct {
	BeatmapAR float32 `json:"AR"`
	BeatmapCS float32 `json:"CS"`
	BeatmapOD float32 `json:"OD"`
	BeatmapHP float32 `json:"HP"`
}

type path struct {
	InnerBGPath          string `json:"full"`
	BeatmapFolderString  string `json:"folder"`
	BeatmapOsuFileString string `json:"file"`
	BGPath               string `json:"bg"`
}

type modsM struct {
	AppliedMods int32  `json:"num"`
	PpMods      string `json:"str"`
}

type ppM struct {
	PpSS int32 `json:"100"`
	Pp99 int32 `json:"99"`
	Pp98 int32 `json:"98"`
	Pp97 int32 `json:"97"`
	Pp96 int32 `json:"96"`
	Pp95 int32 `json:"95"`
}

type combo struct {
	Current int32 `json:"current"`
	Max     int32 `json:"max"`
}

type hp struct {
	Normal float64 `json:"normal"`
	Smooth float64 `json:"smooth"`
}

type hits struct {
	H300          int16   `json:"300"`
	H100          int16   `json:"100"`
	H50           int16   `json:"50"`
	H0            int16   `json:"0"`
	HitErrorArray []int32 `json:"hitErrorArray"`
}

type modsG struct {
	AppliedMods int32 `json:"num"`
	StrMods     int32 `json:"str"`
}

type ppG struct {
	Pp     int32 `json:"current"`
	PPifFC int32 `json:"fc"`
}

type dynamicAddresses struct {
	PlayContainer38 uint32
	BeatmapAddr     uint32
	IsReady         bool
}

//MenuData contains raw values taken from osu! memory
var MenuData = InMenuValues{}

//GameplayData contains raw values taken from osu! memory
var GameplayData = GameplayValues{}

//DynamicAddresses are in-between pointers that lead to values
var DynamicAddresses = dynamicAddresses{}
