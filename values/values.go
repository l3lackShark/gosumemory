package values

//InMenuValues inside osu!memory
type InMenuValues struct {
	BeatmapAddr          uint32
	OsuStatus            uint32  `json:"osuState"`
	BeatmapID            uint32  `json:"bmID"`
	BeatmapSetID         uint32  `json:"bmSetID"`
	BeatmapCS            float32 `json:"CS"`
	BeatmapAR            float32 `json:"AR"`
	BeatmapOD            float32 `json:"OD"`
	BeatmapHP            float32 `json:"HP"`
	BeatmapString        string  `json:"bmInfo"`
	BeatmapFolderString  string  `json:"bmFolder"`
	BeatmapOsuFileString string  `json:"pathToBM"`
	HitObjectStats       string  `json:"bmStats"`
	PlayTime             int32   `json:"bmTime"`
	InnerBGPath          string  `json:"innerBG"`
	BGPath               string
	AppliedMods          int32  `json:"appliedMods"`
	PpMods               string `json:"appliedModsString"`
	PpSS                 string `json:"ppSS"`
	Pp99                 string `json:"pp99"`
	Pp98                 string `json:"pp98"`
	Pp97                 string `json:"pp97"`
	Pp96                 string `json:"pp96"`
	Pp95                 string `json:"pp95"`
}

//GameplayValues inside osu!memory
type GameplayValues struct {
	CurrentHit300c  int16   `json:"300"`
	CurrentHit100c  int16   `json:"100"`
	CurrentHit50c   int16   `json:"50"`
	CurrentHitMiss  int16   `json:"miss"`
	CurrentAccuracy float64 `json:"accuracy"`
	CurrentScore    int32   `json:"score"`
	CurrentCombo    int32   `json:"combo"`
	CurrentGameMode int32   `json:"gameMode"`
	CurrentMaxCombo int32   `json:"maxCombo"`
	// CurrentPlayerHP         int8    `json:"playerHP"`
	// CurrentPlayerHPSmoothed int8    `json:"playerHPSmoothed"`
	Pp     string `json:"pp"`
	PPifFC string `json:"ppIfFC"`
}

//MenuData contains raw values taken from osu! memory
var MenuData = InMenuValues{}

//GameplayData contains raw values taken from osu! memory
var GameplayData = GameplayValues{}
