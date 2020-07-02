package db

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/k0kubun/pp"
)

type osudb struct {
	buildVer          int32
	songsFolderSize   int32
	isAccountUnlocked bool
	Nickname          string
	BmInfo            []beatmapInfo
}

type beatmapInfo struct {
	Artist                    string
	artistU                   string
	Title                     string
	titleU                    string
	Creator                   string
	Difficulty                string
	audioName                 string
	md5                       string
	Filename                  string
	rankedStatus              int8
	NumHitCircles             int16
	NumSliders                int16
	NumSpinners               int16
	dateTime                  uint64
	approachRate              float32
	circleSize                float32
	hpDrain                   float32
	overallDifficulty         float32
	sliderVelocity            float64 //double
	starRatingOsu             []starRating
	starRatingTaiko           []starRating
	starRatingCtb             []starRating
	StarRatingMania           []starRating
	drainTime                 int32
	totalTime                 int32
	previewTime               int32
	timingPoints              []timingPoint
	beatmapID                 int32
	beatmapSetID              int32
	threadID                  int32
	gradeOsu                  int8
	gradeTaiko                int8
	gradeCtb                  int8
	gradeMania                int8
	localOffset               int16
	stackLeniency             float32
	gameMode                  int8
	songSource                string
	songTags                  string
	onlineOffset              int16
	fontTitle                 string //?
	isUnplayed                bool
	lastPlayed                uint64
	isOsz2                    bool
	folderFromSongs           string
	lastCheckedAgainstOsuRepo uint64
	isBmSoundIgnored          bool
	isBmSkinIgnored           bool
	isBmStoryBoardDisabled    bool
	isBmVideoDisabled         bool
	isVisualOverride          bool
	lastClosedEditor          int32
	maniaScrollSpeed          uint8
}

type starRating struct {
	zeroXeight uint8
	BitMods    int32
	zeroXzerod uint8
	StarRating float64 //double
}

type timingPoint struct {
	msPerBeat            float64 //double
	songOffset           float64 //double
	inheritedTimingPoint bool
}

//OsuDB is a structure representation of osu!.db file
var OsuDB osudb

//InitDB initializes osu database and gets data within it
func InitDB() error {
	pp.Println("Parsing osu!db...")
	folder, err := os.Getwd()
	if err != nil {

	}
	file, err := os.Open(filepath.Join(folder, "osu!.db"))
	osuDB := bufio.NewReader(file)
	defer file.Close()
	binary.Read(osuDB, binary.LittleEndian, &OsuDB.buildVer)
	binary.Read(osuDB, binary.LittleEndian, &OsuDB.songsFolderSize)
	binary.Read(osuDB, binary.LittleEndian, &OsuDB.isAccountUnlocked)
	var dateTime int64
	binary.Read(osuDB, binary.LittleEndian, &dateTime)
	OsuDB.Nickname, err = readDBString(osuDB)
	if err != nil {
		log.Println("Database  parse error: ", err)
	}
	OsuDB.BmInfo, err = readDBArray(osuDB)
	if err != nil {
		panic(err)
	}

	pp.Println("Done parsing osu!db")

	return nil
}

func readVarUint(r io.Reader, n uint) (uint64, error) {
	if n > 64 {
		panic(errors.New("leb128: n must <= 64"))
	}
	p := make([]byte, 1)
	var res uint64
	var shift uint
	for {
		_, err := io.ReadFull(r, p)
		if err != nil {
			return 0, err
		}
		b := uint64(p[0])
		switch {
		// note: can not use b < 1<<n, when n == 64, 1<<n will overflow to 0
		case b < 1<<7 && b <= 1<<n-1:
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		default:
			return 0, errors.New("leb128: invalid uint")
		}
	}
}

func readDBString(osuDB io.Reader) (string, error) {
	var checkByte byte
	err := binary.Read(osuDB, binary.LittleEndian, &checkByte)
	if err != nil {
		return "", err
	}
	switch checkByte {
	case 0x00:
		return "", nil
	case 0x0b:
		strlen, err := readVarUint(osuDB, 32)
		if err != nil {
			return "", err
		}
		stringBytes := make([]byte, int(strlen))
		_, err = io.ReadFull(osuDB, stringBytes)
		if err != nil {
			return "", err
		}
		return string(stringBytes[:]), nil

	default:
		return "", errors.New("string parse error")
	}
}
func readDBArray(osuDB io.Reader) ([]beatmapInfo, error) {
	var arrLength int32
	err := binary.Read(osuDB, binary.LittleEndian, &arrLength)
	if err != nil {
		return nil, err
	}
	if arrLength == -1 {
		return nil, nil
	}
	beatmapsArray := make([]beatmapInfo, int(arrLength))
	for i := 0; i < int(arrLength); i++ {
		beatmapsArray[i], err = readBeatmapInfo(osuDB)
		if err != nil {
			return nil, err
		}
	}
	return beatmapsArray, nil
}
func readBeatmapInfo(osuDB io.Reader) (beatmapInfo, error) {

	data := beatmapInfo{}
	var err error
	data.Artist, err = readDBString(osuDB)
	data.artistU, err = readDBString(osuDB)
	data.Title, err = readDBString(osuDB)
	data.titleU, err = readDBString(osuDB)
	data.Creator, err = readDBString(osuDB)
	data.Difficulty, err = readDBString(osuDB)
	data.audioName, err = readDBString(osuDB)
	data.md5, err = readDBString(osuDB)
	data.Filename, err = readDBString(osuDB)
	err = binary.Read(osuDB, binary.LittleEndian, &data.rankedStatus)
	err = binary.Read(osuDB, binary.LittleEndian, &data.NumHitCircles)
	err = binary.Read(osuDB, binary.LittleEndian, &data.NumSliders)
	err = binary.Read(osuDB, binary.LittleEndian, &data.NumSpinners)
	err = binary.Read(osuDB, binary.LittleEndian, &data.dateTime)
	err = binary.Read(osuDB, binary.LittleEndian, &data.approachRate)
	err = binary.Read(osuDB, binary.LittleEndian, &data.circleSize)
	err = binary.Read(osuDB, binary.LittleEndian, &data.hpDrain)
	err = binary.Read(osuDB, binary.LittleEndian, &data.overallDifficulty)
	err = binary.Read(osuDB, binary.LittleEndian, &data.sliderVelocity)
	var lengthList int32 // should move this into a separate functuion and use reflections to set values
	err = binary.Read(osuDB, binary.LittleEndian, &lengthList)
	if lengthList >= 1 {
		data.starRatingOsu = make([]starRating, int(lengthList))
		for i := 0; i < int(lengthList); i++ {
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingOsu[i].zeroXeight)
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingOsu[i].BitMods)
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingOsu[i].zeroXzerod)
			if data.starRatingOsu[i].zeroXzerod != 0x0d || data.starRatingOsu[i].zeroXeight != 0x08 {
				pp.Println("Star rating parse err.")
				data := beatmapInfo{}
				return data, errors.New("Star rating parse err")
			}
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingOsu[i].StarRating)
			if err != nil {
				data := beatmapInfo{}
				return data, err
			}
		}
	}

	var lengthListTaiko int32 // should move this into a separate functuion and use reflections to set values
	err = binary.Read(osuDB, binary.LittleEndian, &lengthListTaiko)
	if lengthListTaiko >= 1 {
		data.starRatingTaiko = make([]starRating, int(lengthListTaiko))
		for i := 0; i < int(lengthListTaiko); i++ {
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingTaiko[i].zeroXeight)
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingTaiko[i].BitMods)
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingTaiko[i].zeroXzerod)
			if data.starRatingTaiko[i].zeroXzerod != 0x0d || data.starRatingTaiko[i].zeroXeight != 0x08 {
				pp.Println("Star rating parse err. (taiko)")
				data := beatmapInfo{}
				return data, errors.New("Star rating parse err (taiko)")
			}
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingTaiko[i].StarRating)
			if err != nil {
				data := beatmapInfo{}
				return data, err
			}
		}
	}

	var lengthListCtb int32 // should move this into a separate functuion and use reflections to set values
	err = binary.Read(osuDB, binary.LittleEndian, &lengthListCtb)
	if lengthListCtb >= 1 {
		data.starRatingCtb = make([]starRating, int(lengthListCtb))
		for i := 0; i < int(lengthListCtb); i++ {
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingCtb[i].zeroXeight)
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingCtb[i].BitMods)
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingCtb[i].zeroXzerod)
			if data.starRatingCtb[i].zeroXzerod != 0x0d || data.starRatingCtb[i].zeroXeight != 0x08 {
				pp.Println("Star rating parse err. (ctb)")
				data := beatmapInfo{}
				return data, errors.New("Star rating parse err (ctb)")
			}
			err = binary.Read(osuDB, binary.LittleEndian, &data.starRatingCtb[i].StarRating)
			if err != nil {
				data := beatmapInfo{}
				return data, err
			}
		}
	}
	var lengthListMania int32 // should move this into a separate functuion and use reflections to set values
	err = binary.Read(osuDB, binary.LittleEndian, &lengthListMania)
	if lengthListMania >= 1 {
		data.StarRatingMania = make([]starRating, int(lengthListMania))
		for i := 0; i < int(lengthListMania); i++ {
			err = binary.Read(osuDB, binary.LittleEndian, &data.StarRatingMania[i].zeroXeight)
			err = binary.Read(osuDB, binary.LittleEndian, &data.StarRatingMania[i].BitMods)
			err = binary.Read(osuDB, binary.LittleEndian, &data.StarRatingMania[i].zeroXzerod)
			if data.StarRatingMania[i].zeroXzerod != 0x0d || data.StarRatingMania[i].zeroXeight != 0x08 {
				pp.Println("Star rating parse err. (Mania)")
				data := beatmapInfo{}
				return data, errors.New("Star rating parse err (Mania)")
			}
			err = binary.Read(osuDB, binary.LittleEndian, &data.StarRatingMania[i].StarRating)
			if err != nil {
				data := beatmapInfo{}
				return data, err
			}
		}
	}
	err = binary.Read(osuDB, binary.LittleEndian, &data.drainTime)
	err = binary.Read(osuDB, binary.LittleEndian, &data.totalTime)
	err = binary.Read(osuDB, binary.LittleEndian, &data.previewTime)

	var lengthTimingPoints int32
	err = binary.Read(osuDB, binary.LittleEndian, &lengthTimingPoints)
	if lengthTimingPoints >= 1 {
		data.timingPoints = make([]timingPoint, int(lengthTimingPoints))
		for i := 0; i < int(lengthTimingPoints); i++ {
			err = binary.Read(osuDB, binary.LittleEndian, &data.timingPoints[i].msPerBeat)
			err = binary.Read(osuDB, binary.LittleEndian, &data.timingPoints[i].songOffset)
			err = binary.Read(osuDB, binary.LittleEndian, &data.timingPoints[i].inheritedTimingPoint)
			if err != nil {
				data := beatmapInfo{}
				return data, err
			}
		}
	}
	err = binary.Read(osuDB, binary.LittleEndian, &data.beatmapID)
	err = binary.Read(osuDB, binary.LittleEndian, &data.beatmapSetID)
	err = binary.Read(osuDB, binary.LittleEndian, &data.threadID)
	err = binary.Read(osuDB, binary.LittleEndian, &data.gradeOsu)
	err = binary.Read(osuDB, binary.LittleEndian, &data.gradeTaiko)
	err = binary.Read(osuDB, binary.LittleEndian, &data.gradeCtb)
	err = binary.Read(osuDB, binary.LittleEndian, &data.gradeMania)
	err = binary.Read(osuDB, binary.LittleEndian, &data.localOffset)
	err = binary.Read(osuDB, binary.LittleEndian, &data.stackLeniency)
	err = binary.Read(osuDB, binary.LittleEndian, &data.gameMode)
	data.songSource, err = readDBString(osuDB)
	data.songTags, err = readDBString(osuDB)
	err = binary.Read(osuDB, binary.LittleEndian, &data.onlineOffset)
	data.fontTitle, err = readDBString(osuDB)
	err = binary.Read(osuDB, binary.LittleEndian, &data.isUnplayed)
	err = binary.Read(osuDB, binary.LittleEndian, &data.lastPlayed)
	err = binary.Read(osuDB, binary.LittleEndian, &data.isOsz2)
	data.folderFromSongs, err = readDBString(osuDB)
	err = binary.Read(osuDB, binary.LittleEndian, &data.lastCheckedAgainstOsuRepo)
	err = binary.Read(osuDB, binary.LittleEndian, &data.isBmSoundIgnored)
	err = binary.Read(osuDB, binary.LittleEndian, &data.isBmSkinIgnored)
	err = binary.Read(osuDB, binary.LittleEndian, &data.isBmStoryBoardDisabled)
	err = binary.Read(osuDB, binary.LittleEndian, &data.isBmVideoDisabled)
	err = binary.Read(osuDB, binary.LittleEndian, &data.isVisualOverride)
	err = binary.Read(osuDB, binary.LittleEndian, &data.lastClosedEditor)
	err = binary.Read(osuDB, binary.LittleEndian, &data.maniaScrollSpeed)

	if err != nil {
		data := beatmapInfo{}
		return data, err
	}
	return data, nil
}
