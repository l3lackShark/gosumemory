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
	PlayTime             uint32  `json:"bmTime"`
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
	IsReady              bool
}

//GameplayValues inside osu!memory
type GameplayValues struct {
	PlayContainer38 uint32
	AppliedMods     int32   `json:"appliedMods"`
	Hit300c         int16   `json:"300"`
	Hit100c         int16   `json:"100"`
	Hit50c          int16   `json:"50"`
	HitMiss         int16   `json:"miss"`
	Accuracy        float64 `json:"accuracy"`
	Score           int32   `json:"score"`
	Combo           int32   `json:"combo"`
	GameMode        int32   `json:"gameMode"`
	MaxCombo        int32   `json:"maxCombo"`
	Pp              string  `json:"pp"`
	PPifFC          float64 `json:"ppIfFC"`
}

//MenuData contains raw values taken from osu! memory
var MenuData = InMenuValues{}

//GameplayData contains raw values taken from osu! memory
var GameplayData = GameplayValues{}
