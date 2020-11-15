package memory

//InMenuValues inside osu!memory
type InMenuValues struct {
	MainMenuValues MainMenuValues `json:"mainMenu"`
	OsuStatus      uint32         `json:"state"`
	SkinFolder     string         `json:"skinFolder"`
	GameMode       int32          `json:"gameMode"`
	ChatChecker    int8           `json:"isChatEnabled"` //bool (1 byte)
	Bm             bm             `json:"bm"`
	Mods           modsM          `json:"mods"`
	PP             ppM            `json:"pp"`
}

type ResultsScreenValues struct {
	Name     string `json:"name"`
	Score    int32  `json:"score"`
	MaxCombo int16  `json:"maxCombo"`
	Mods     modsM  `json:"mods"`
	H300     int16  `json:"300"`
	HGeki    int16  `json:"geki"`
	H100     int16  `json:"100"`
	HKatu    int16  `json:"katu"`
	H50      int16  `json:"50"`
	H0       int16  `json:"0"`
}

type MainMenuValues struct {
	BassDensity float64 `json:"bassDensity"`
}

//InSettingsValues are values represented inside settings class, could be dynamic
type InSettingsValues struct {
	ShowInterface bool `json:"showInterface"` //dynamic in gameplay
}

type TourneyValues struct {
	Manager    tourneyManager `json:"manager"`
	IPCClients []ipcClient    `json:"ipcClients"`
}

type tourneyManager struct {
	IPCState int32            `json:"ipcState"`
	BO       int32            `json:"bestOF"`
	Name     tName            `json:"teamName"`
	Stars    tStars           `json:"stars"`
	Bools    tBools           `json:"bools"`
	Chat     []tourneyMessage `json:"chat"`
	Gameplay tmGameplay       `json:"gameplay"`
}

type tourneyMessage struct {
	Time        string `json:"time"`
	Name        string `json:"name"`
	MessageBody string `json:"messageBody"`
}

type tmGameplay struct {
	Score tScore `json:"score"`
}

type tBools struct {
	ScoreVisible bool `json:"scoreVisible"`
	StarsVisible bool `json:"starsVisible"`
}

type tName struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}
type tStars struct {
	Left  int32 `json:"left"`
	Right int32 `json:"right"`
}
type tScore struct {
	Left  int32 `json:"left"`
	Right int32 `json:"right"`
}

type ipcClient struct {
	SpectatingID int32           `json:"spectatingID"`
	Gameplay     tourneyGameplay `json:"gameplay"`
}

type tourneyGameplay struct {
	GameMode int32       `json:"gameMode"`
	Score    int32       `json:"score"`
	Name     string      `json:"name"`
	Accuracy float64     `json:"accuracy"`
	Hits     tourneyHits `json:"hits"`
	Combo    combo       `json:"combo"`
	Mods     modsM       `json:"mods"`
	Hp       hp          `json:"hp"`
}

type gGrade struct {
	Current  string `json:"current"`
	Expected string `json:"maxThisPlay"`
}

//GameplayValues inside osu!memory
type GameplayValues struct {
	GameMode int32 `json:"gameMode"`
	//BitwiseKeypress int8        `json:"bitwiseKeypress"`
	Name        string      `json:"name"`
	Score       int32       `json:"score"`
	Accuracy    float64     `json:"accuracy"`
	Combo       combo       `json:"combo"`
	Hp          hp          `json:"hp"`
	Hits        hits        `json:"hits"`
	PP          ppG         `json:"pp"`
	Leaderboard leaderboard `json:"leaderboard"`
}

type bm struct {
	Time           tim      `json:"time"`
	BeatmapID      int32    `json:"id"`
	BeatmapSetID   int32    `json:"set"`
	BeatmapMD5     string   `json:"md5"`
	RandkedStatus  int32    `json:"rankedStatus"` //unknown, unsubmitted, pending/wip/graveyard, unused, ranked, approved, qualified
	Metadata       Metadata `json:"metadata"`
	Stats          stats    `json:"stats"`
	Path           path     `json:"path"`
	HitObjectStats string   `json:"-"`
	BeatmapString  string   `json:"-"`
}

type tim struct {
	FirstObj int32 `json:"firstObj"`
	PlayTime int32 `json:"current"`
	FullTime int32 `json:"full"`
	Mp3Time  int32 `json:"mp3"`
}

// Metadata Map data
type Metadata struct {
	Artist  string `json:"artist"`
	Title   string `json:"title"`
	Mapper  string `json:"mapper"`
	Version string `json:"difficulty"`
}

type stats struct {
	BeatmapAR       float32 `json:"AR"`
	BeatmapCS       float32 `json:"CS"`
	BeatmapOD       float32 `json:"OD"`
	BeatmapHP       float32 `json:"HP"`
	BeatmapSR       float32 `json:"SR"`
	BeatmapBPM      bpm     `json:"BPM"`
	BeatmapMaxCombo int32   `json:"maxCombo"`
	FullSR          float32 `json:"fullSR"`
	MemoryAR        float32 `json:"memoryAR"`
	MemoryCS        float32 `json:"memoryCS"`
	MemoryOD        float32 `json:"memoryOD"`
	MemoryHP        float32 `json:"memoryHP"`
}

type bpm struct {
	Minimal int `json:"min"`
	Maximal int `json:"max"`
}

type path struct {
	InnerBGPath          string `json:"full"`
	BeatmapFolderString  string `json:"folder"`
	BeatmapOsuFileString string `json:"file"`
	BGPath               string `json:"bg"`
	AudioPath            string `json:"audio"`
	FullMP3Path          string `json:"-"`
	FullDotOsu           string `json:"-"`
}

type modsM struct {
	AppliedMods int32  `json:"num"`
	PpMods      string `json:"str"`
}

type ppM struct {
	PpSS      int32     `json:"100"`
	Pp99      int32     `json:"99"`
	Pp98      int32     `json:"98"`
	Pp97      int32     `json:"97"`
	Pp96      int32     `json:"96"`
	Pp95      int32     `json:"95"`
	PpStrains []float64 `json:"strains"`
}

type combo struct {
	Current int16 `json:"current"`
	Max     int16 `json:"max"`
	Temp    int16 `json:"-"`
}

type hp struct {
	Normal float64 `json:"normal"`
	Smooth float64 `json:"smooth"`
}

type hits struct {
	H300          int16   `json:"300"`
	HGeki         int16   `json:"geki"`
	H100          int16   `json:"100"`
	HKatu         int16   `json:"katu"`
	H50           int16   `json:"50"`
	H0            int16   `json:"0"`
	H0Temp        int16   `json:"-"`
	HSB           int16   `json:"sliderBreaks"`
	Grade         gGrade  `json:"grade"`
	UnstableRate  float64 `json:"unstableRate"`
	HitErrorArray []int32 `json:"hitErrorArray"`
}

type tourneyHits struct {
	H300          int16   `json:"300"`
	HGeki         int16   `json:"geki"`
	H100          int16   `json:"100"`
	HKatu         int16   `json:"katu"`
	H50           int16   `json:"50"`
	H0            int16   `json:"0"`
	H0Temp        int16   `json:"-"`
	HSB           int16   `json:"sliderBreaks"`
	UnstableRate  float64 `json:"unstableRate"`
	HitErrorArray []int32 `json:"hitErrorArray"`
}

type ppG struct {
	Pp            int32 `json:"current"`
	PPifFC        int32 `json:"fc"`
	PPMaxThisPlay int32 `json:"maxThisPlay"`
}

type dynamicAddresses struct {
	IsReady bool
}

type leaderPlayer struct {
	Name      string `json:"name"`
	Score     int32  `json:"score"`
	Combo     int16  `json:"combo"`
	MaxCombo  int16  `json:"maxCombo"`
	Mods      string `json:"mods"`
	H300      int16  `json:"h300"`
	H100      int16  `json:"h100"`
	H50       int16  `json:"h50"`
	H0        int16  `json:"h0"`
	Team      int32  `json:"team"`
	Position  int32  `json:"position"`
	IsPassing int8   `json:"isPassing"` //bool
}

type leaderboard struct {
	DoesLeaderBoardExists bool           `json:"hasLeaderboard"`
	IsLeaderBoardVisible  bool           `json:"isVisible"`
	OurPlayer             leaderPlayer   `json:"ourplayer"`
	Slots                 []leaderPlayer `json:"slots"`
}

//MenuData contains raw values taken from osu! memory
var MenuData = InMenuValues{}

//GameplayData contains raw values taken from osu! memory
var GameplayData = GameplayValues{}

//ResultsScreenData contains raw values taken from osu! memory
var ResultsScreenData = ResultsScreenValues{}

//SettingsData contains raw values taken from osu! memory
var SettingsData = InSettingsValues{}

//TourneyData contains raw values taken from osu! memory
var TourneyData = TourneyValues{}

//DynamicAddresses are in-between pointers that lead to values
var DynamicAddresses = dynamicAddresses{}
