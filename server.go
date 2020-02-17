package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/Andoryuuta/kiwi"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
)

var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")

var osuBase uintptr
var osuStatus uint16
var currentBeatmapData uintptr
var playContainer uintptr
var playContainerBase uintptr
var serverBeatmapString string
var outStrLoop string

func Cmd(cmd string, shell bool) []byte {

	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out

	}
}

func OsuStatusAddr() uintptr { //in hopes to deprecate this
	x := Cmd("scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;48 83 F8 04 73 1E;list;exit'", true)
	outStr := cast.ToString(x)
	outStr = strings.Replace(outStr, " ", "", -1)

	input := outStr
	if input == "" {
		log.Fatalln("osuStatus addr fail")
	}
	output := (input[3:])
	yosuBase := firstN(output, 8)
	osuBaseString := "0x" + yosuBase
	osuBaseUINT32 := cast.ToUint32(osuBaseString)
	osuBase = uintptr(osuBaseUINT32)
	if osuBase == 0 {
		log.Fatalln("could not find osuStatusAddr, is osu! running?")
	}
	//println(CurrentBeatmapFolderString())
	return osuBase
}
func OsuBaseAddr() uintptr { //in hopes to deprecate this
	x := Cmd("scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;F8 01 74 04 83;list;exit'", true)
	outStr := cast.ToString(x)
	outStr = strings.Replace(outStr, " ", "", -1)

	input := outStr
	if input == "" {
		log.Fatalln("OsuBase addr fail")
	}
	output := (input[3:])
	yosuBase := firstN(output, 8)
	osuBaseString := "0x" + yosuBase
	osuBaseUINT32 := cast.ToUint32(osuBaseString)
	osuBase = uintptr(osuBaseUINT32)
	//println(CurrentBeatmapFolderString())
	if osuBase == 0 {
		log.Fatalln("Could not find OsuBaseAddr, is osu! running?")
	}
	return osuBase
}
func OsuplayContainer() uintptr { //in hopes to deprecate this
	x := Cmd("scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;85 C9 74 1F 8D 55 F0 8B 01;list;exit'", true)
	outStr := cast.ToString(x)
	outStr = strings.Replace(outStr, " ", "", -1)

	input := outStr
	if input == "" {
		log.Fatalln("osuplayContainer addr fail")
	}
	output := (input[3:])
	yosuBase := firstN(output, 8)
	osuBaseString := "0x" + yosuBase

	osuBaseUINT32 := cast.ToUint32(osuBaseString)
	osuBase = uintptr(osuBaseUINT32)
	if osuBase == 0 {
		log.Fatalln("Could not find osuplayContainer address, is osu! running?")
	}
	//println(CurrentBeatmapFolderString())
	return osuBase
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	if procerr != nil { //TODO: refactor
		ws.WriteMessage(1, []byte("osu!.exe not found"))
		log.Fatalln("is osu! running? (osu! process was not found)")
	}
	StaticOsuStatusAddr := OsuStatusAddr() //we should only check for this address once.
	osuStatusOffset, err := proc.ReadUint32(StaticOsuStatusAddr - 0x4)
	if err != nil {
		ws.WriteMessage(1, []byte("osu!status offset was not found"))
		log.Fatalln("osu!status offset was not found, are you sure that osu!stable is running? If so, please report this to GitHub!")
	}
	uintptrOsuStatus := uintptr(osuStatusOffset)
	osuStatusValue, err := proc.ReadUint16(uintptrOsuStatus)
	if err != nil {
		ws.WriteMessage(1, []byte("osu!status value was not found"))
		log.Fatalln("osu!status value was not found, are you sure that osu!stable is running? If so, please report this to GitHub!")
	}
	for osuStatusValue != 5 {
		log.Println("please go to songselect in order to proceed!")
		osuStatusValue, err = proc.ReadUint16(uintptrOsuStatus)
		if err != nil {
			log.Fatalln("is osu! running? (osu! status was not found)")
		}
		ws.WriteMessage(1, []byte("osu! is not in SongSelect!"))

		time.Sleep(500 * time.Millisecond)

	}
	fmt.Println("it seems that the client is in song select, you are good to go!")

	log.Println("Client Connected")

	osuBase = OsuBaseAddr()
	currentBeatmapData = (osuBase - 0xC)
	playContainer = OsuplayContainer()
	playContainerBase = (playContainer - 0x4)
	if err != nil {
		log.Fatalln("is osu! running? (osu! status offset was not found)")
	}

	for {
		var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
		if procerr != nil {
			log.Fatalln("is osu! running? (osu! process was not found, terminating...)")
		}
		osuStatusValue, err := proc.ReadUint16(uintptrOsuStatus)
		if err != nil {
			log.Println("osu! status could not be found...", err)

		}
		osuStatus = osuStatusValue
		type PlayContainer struct {
			CurrentHit300c     int16   `json:"300"`
			CurrentHit100c     int16   `json:"100"`
			CurrentHit50c      int16   `json:"50"`
			CurrentHitMiss     int16   `json:"miss"`
			CurrentAccuracy    float64 `json:"accuracy"`
			CurrentScore       int32   `json:"score"`
			CurrentCombo       int32   `json:"combo"`
			CurrentGameMode    int32   `json:"gameMode"`
			CurrentAppliedMods int32   `json:"appliedMods"`
		}
		type EverythingInMenu struct {
			CurrentState                uint16  `json:"osuState"`
			CurrentBeatmapID            uint32  `json:"bmID"`
			CurrentBeatmapSetID         uint32  `json:"bmSetID"`
			CurrentBeatmapCS            float32 `json:"CS"`
			CurrentBeatmapAR            float32 `json:"AR"`
			CurrentBeatmapOD            float32 `json:"OD"`
			CurrentBeatmapHP            float32 `json:"HP"`
			CurrentBeatmapString        string  `json:"bmInfo"`
			CurrentBeatmapFolderString  string  `json:"bmFolder"`
			CurrentBeatmapOsuFileString string  `json:"pathToBM"`
		}

		type EverythingInMenu2 struct { //order sets here
			D EverythingInMenu `json:"menuContainer"`
			P PlayContainer    `json:"gameplayContainer"`
		}

		PlayContainerStruct := PlayContainer{
			CurrentHit300c: CurrentHit300c(), CurrentHit100c: CurrentHit100c(), CurrentHit50c: CurrentHit50c(), CurrentHitMiss: CurrentHitMiss(),
			CurrentScore: CurrentScore(), CurrentAccuracy: CurrentAccuracy(), CurrentCombo: CurrentCombo(), CurrentGameMode: CurrentGameMode(), CurrentAppliedMods: CurrentAppliedMods(),
		}

		//println(ValidCurrentBeatmapFolderString())
		if strings.HasSuffix(ValidCurrentBeatmapOsuFileString(), ".osu") == false {
			println(".osu ends with ???")
		}
		if strings.HasSuffix(ValidCurrentBeatmapString(), "]") == false {
			println("beatmapstring ends with ???")
		}
		MenuContainerStruct := EverythingInMenu{CurrentState: osuStatus,
			CurrentBeatmapID:            CurrentBeatmapID(),
			CurrentBeatmapSetID:         CurrentBeatmapSetID(),
			CurrentBeatmapString:        ValidCurrentBeatmapString(),
			CurrentBeatmapFolderString:  ValidCurrentBeatmapFolderString(),  //TODO: fix strings  //kind of fixed via monkaW functions
			CurrentBeatmapOsuFileString: ValidCurrentBeatmapOsuFileString(), //TODO: fix strings  //kind of fixed via monkaW functions
			//CurrentBeatmapOsuFileString: runesStr,

			CurrentBeatmapAR: CurrentBeatmapAR(),
			CurrentBeatmapOD: CurrentBeatmapOD(),
			CurrentBeatmapCS: CurrentBeatmapCS(),
			CurrentBeatmapHP: CurrentBeatmapHP(),
		}
		group := EverythingInMenu2{
			P: PlayContainerStruct,
			D: MenuContainerStruct,
		}
		b, err := json.Marshal(group)
		if err != nil {
			fmt.Println("error:", err)
		}

		ws.WriteMessage(1, []byte(b)) //sending data to the client

		//if err != nil {
		//	log.Println(err)
		//}
		time.Sleep(100 * time.Millisecond)

	}

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {

	setupRoutes()
	log.Fatal(http.ListenAndServe(":8085", nil))
}

func firstN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func CurrentBeatmapID() uint32 {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapIDBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("CurrentBeatmapID Base level failure")
		return 0
	}
	beatmapIDFirstLevel, err := proc.ReadUint32(uintptr(beatmapIDBase))
	if err != nil {
		log.Println("CurrentBeatmapID First level pointer failure")
		return 0
	}
	currentBeatmapID, err := proc.ReadUint32(uintptr(beatmapIDFirstLevel + 0xC4))
	if err != nil {
		log.Println("CurrentBeatmapID result pointer failure")
		return 0
	}
	return currentBeatmapID
}
func CurrentBeatmapSetID() uint32 {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapSetIDBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("CurrentBeatmapSetID Base level failure")
		return 0
	}
	beatmapSetIDFirstLevel, err := proc.ReadUint32(uintptr(beatmapSetIDBase))
	if err != nil {
		log.Println("CurrentBeatmapSetID First level pointer failure")
		return 0
	}
	currentSetBeatmapID, err := proc.ReadUint32(uintptr(beatmapSetIDFirstLevel + 0xC8))
	if err != nil {
		log.Println("CurrentBeatmapSetID result pointer failure")
		return 0
	}
	return currentSetBeatmapID
}
func CurrentBeatmapString() string {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapStringBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("BeatMapString Base level failure")
		return "-2"
	}
	beatmapStringFirstLevel, err := proc.ReadUint32(uintptr(beatmapStringBase))
	if err != nil {
		log.Println("BeatMapString First level pointer failure")
		return "-3"
	}
	beatmapStringSecondLevel, err := proc.ReadUint32(uintptr(beatmapStringFirstLevel + 0x7C))
	if err != nil {
		log.Println("BeatMapString Second level pointer failure")
		return "-4"
	}
	beatmapStringResult, err := proc.ReadBytes(uintptr(beatmapStringSecondLevel+0x8), 500) // fix repeating
	if err != nil {
		log.Println("BeatMapString Third level pointer failure")
		return "-5"
	}
	beatmapString := cast.ToString(beatmapStringResult)
	return beatmapString
}

func CurrentBeatmapFolderString() string {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapFolderStringBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("BeatMapFolderString Base level failure")
		return "-2"
	}
	beatmapFolderStringFirstLevel, err := proc.ReadUint32(uintptr(beatmapFolderStringBase))
	if err != nil {
		log.Println("BeatMapFolderString First level pointer failure")
		return "-3"
	}
	beatmapFolderStringSecondLevel, err := proc.ReadUint32(uintptr(beatmapFolderStringFirstLevel + 0x74))
	if err != nil {
		log.Println("BeatMapFolderString Second level pointer failure")
		return "-4"
	}
	beatmapFolderStringResult, err := proc.ReadBytes(uintptr(beatmapFolderStringSecondLevel+0x8), 500) // fix repeating
	if err != nil {
		log.Println("BeatMapFolderString result level pointer failure")
		return "-5"
	}
	beatmapString := string(beatmapFolderStringResult)
	return beatmapString
}
func CurrentBeatmapOsuFileString() string {
	proc := proc
	//currentBeatmapData := currentBeatmapData

	beatmapFolderStringBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("BeatMapOsuFileString Base level failure")
		return "-2"
	}
	beatmapFolderStringFirstLevel, err := proc.ReadUint32(uintptr(beatmapFolderStringBase))
	if err != nil {
		log.Println("BeatMapOsuFileString First level pointer failure")
		return "-3"
	}
	beatmapFolderStringSecondLevel, err := proc.ReadUint32(uintptr(beatmapFolderStringFirstLevel + 0x8C))
	if err != nil {
		log.Println("BeatMapOsuFileString Second level pointer failure")
		return "-4"
	}
	beatmapFolderStringResult, err := proc.ReadBytes(uintptr(beatmapFolderStringSecondLevel+0x8), 500) // fix repeating
	if err != nil {
		log.Println("BeatMapOsuFileString result level pointer failure")
		return "-5"
	}
	beatmapString := cast.ToString(beatmapFolderStringResult)
	return beatmapString
}
func CurrentBeatmapFolderBase64() []byte {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapFolderStringBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("BeatMapFolderString Base level failure")
		return []byte("-1")
	}
	beatmapFolderStringFirstLevel, err := proc.ReadUint32(uintptr(beatmapFolderStringBase))
	if err != nil {
		log.Println("BeatMapFolderString First level pointer failure")
		return []byte("-2")
	}
	beatmapFolderStringSecondLevel, err := proc.ReadUint32(uintptr(beatmapFolderStringFirstLevel + 0x74))
	if err != nil {
		log.Println("BeatMapFolderString Second level pointer failure")
		return []byte("-3")
	}
	beatmapFolderStringResult, err := proc.ReadBytes(uintptr(beatmapFolderStringSecondLevel+0x8), 500) // fix repeating
	if err != nil {
		log.Println("BeatMapFolderString result level pointer failure")
		return []byte("-5")
	}
	//beatmapString := cast.ToString(beatmapFolderStringResult)
	return beatmapFolderStringResult
}
func CurrentBeatmapOsuFileBase64() []byte {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapFolderStringBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("BeatMapOsuFileString Base level failure")
		return []byte("-2")
	}
	beatmapFolderStringFirstLevel, err := proc.ReadUint32(uintptr(beatmapFolderStringBase))
	if err != nil {
		log.Println("BeatMapOsuFileString First level pointer failure")
		return []byte("-3")
	}
	beatmapFolderStringSecondLevel, err := proc.ReadUint32(uintptr(beatmapFolderStringFirstLevel + 0x8C))
	if err != nil {
		log.Println("BeatMapOsuFileString Second level pointer failure")
		return []byte("-4")
	}
	beatmapFolderStringResult, err := proc.ReadBytes(uintptr(beatmapFolderStringSecondLevel+0x8), 500) // fix repeating
	if err != nil {
		log.Println("BeatMapOsuFileString result level pointer failure")
		return []byte("-5")
	}
	//beatmapString := cast.ToString(beatmapFolderStringResult)
	return beatmapFolderStringResult
}
func CurrentBeatmapAR() float32 {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapSetIDBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("AR Base level failure")
		return -2
	}
	beatmapSetIDFirstLevel, err := proc.ReadUint32(uintptr(beatmapSetIDBase))
	if err != nil {
		log.Println("AR First level pointer failure")
		return -3
	}
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(beatmapSetIDFirstLevel + 0x2C))
	if err != nil {
		log.Println("AR result level pointer failure")
		return -5
	}
	return currentSetBeatmapID
}
func CurrentBeatmapCS() float32 {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapSetIDBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("CS Base level failure")
		return -2
	}
	beatmapSetIDFirstLevel, err := proc.ReadUint32(uintptr(beatmapSetIDBase))
	if err != nil {
		log.Println("CS First level pointer failure")
		return -3
	}
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(beatmapSetIDFirstLevel + 0x30))
	if err != nil {
		log.Println("CS result level pointer failure")
		return -4
	}
	return currentSetBeatmapID
}
func CurrentBeatmapHP() float32 {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapSetIDBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("HP Base level failure")
		return -2
	}
	beatmapSetIDFirstLevel, err := proc.ReadUint32(uintptr(beatmapSetIDBase))
	if err != nil {
		log.Println("HP First level pointer failure")
		return -3
	}
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(beatmapSetIDFirstLevel + 0x34))
	if err != nil {
		log.Println("HP result level pointer failure")
		return -5
	}
	return currentSetBeatmapID
}
func CurrentBeatmapOD() float32 {
	proc := proc
	currentBeatmapData := currentBeatmapData

	beatmapSetIDBase, err := proc.ReadUint32(currentBeatmapData)
	if err != nil {
		log.Println("OD Base level failure")
		return -2
	}
	beatmapSetIDFirstLevel, err := proc.ReadUint32(uintptr(beatmapSetIDBase))
	if err != nil {
		log.Println("OD first level pointer failure")
		return -3
	}
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(beatmapSetIDFirstLevel + 0x38))
	if err != nil {
		log.Println("OD result level pointer failure")
		return -5
	}
	return currentSetBeatmapID
}

// ------------------- PlayContainer
func CurrentAppliedMods() int32 {
	proc := proc
	playContainer := playContainerBase
	if osuStatus != 2 {
		return -1
	}
	comboBase, err := proc.ReadUint32(playContainer)
	if err != nil {
		log.Println("CurrentCombo Base level failure")
		return -2
	}
	comboFirstLevel, err := proc.ReadUint32(uintptr(comboBase))
	if err != nil {
		log.Println("CurrentCombo First level pointer failure")
		return -3
	}
	comboSecondLevel, err := proc.ReadUint32(uintptr(comboFirstLevel) + 0x38)
	if err != nil {
		log.Println("CurrentCombo Second level pointer failure")
		return -4

	}
	currentCombo, err := proc.ReadInt32(uintptr(comboSecondLevel + 0x1C))
	if err != nil {
		log.Println("CurrentCombo result pointer failure")
		return -5
	}
	xorVal1, err := proc.ReadInt32(uintptr(currentCombo + 0xC))
	if err != nil {
		log.Println("CurrentCombo result pointer failure")
		return -6
	}
	xorVal2, err := proc.ReadInt32(uintptr(currentCombo + 0x8))
	if err != nil {
		log.Println("CurrentCombo result pointer failure")
		return -7
	}
	val := xorVal2 - xorVal1
	return val
}
func CurrentCombo() int32 {
	proc := proc
	playContainer := playContainerBase
	if osuStatus != 2 {
		return -1
	}
	comboBase, err := proc.ReadUint32(playContainer)
	if err != nil {
		log.Println("CurrentCombo Base level failure")
		return -2
	}
	comboFirstLevel, err := proc.ReadUint32(uintptr(comboBase))
	if err != nil {
		log.Println("CurrentCombo First level pointer failure")
		return -3
	}
	comboSecondLevel, err := proc.ReadUint32(uintptr(comboFirstLevel) + 0x38)
	if err != nil {
		log.Println("CurrentCombo Second level pointer failure")
		return -4

	}
	currentCombo, err := proc.ReadInt32(uintptr(comboSecondLevel + 0x90))
	if err != nil {
		log.Println("CurrentCombo result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentHit100c() int16 {
	proc := proc
	playContainer := playContainerBase
	if osuStatus != 2 {
		return -1
	}
	comboBase, err := proc.ReadUint32(playContainer)
	if err != nil {
		log.Println("CurrentHit100c Base level failure")
		return -2
	}
	comboFirstLevel, err := proc.ReadUint32(uintptr(comboBase))
	if err != nil {
		log.Println("CurrentHit100c First level pointer failure")
		return -3
	}
	comboSecondLevel, err := proc.ReadUint32(uintptr(comboFirstLevel) + 0x38)
	if err != nil {
		log.Println("CurrentHit100c Second level pointer failure")
		return -4
	}
	currentCombo, err := proc.ReadInt16(uintptr(comboSecondLevel + 0x84)) //2 bytes
	if err != nil {
		log.Println("CurrentHit100c result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentHit300c() int16 {
	proc := proc
	playContainer := playContainerBase
	if osuStatus != 2 {
		return -1
	}
	comboBase, err := proc.ReadUint32(playContainer)
	if err != nil {
		log.Println("CurrentHit300c Base level failure")
		return -2
	}
	comboFirstLevel, err := proc.ReadUint32(uintptr(comboBase))
	if err != nil {
		log.Println("CurrentHit300c First level pointer failure")
		return -3
	}
	comboSecondLevel, err := proc.ReadUint32(uintptr(comboFirstLevel) + 0x38)
	if err != nil {
		log.Println("CurrentHit300c Second level pointer failure")
		return -4
	}
	current300, err := proc.ReadInt16(uintptr(comboSecondLevel + 0x86)) //2 bytes
	if err != nil {
		log.Println("CurrentHit300c result pointer failure")
		return -5
	}
	//currentgeki, err := proc.ReadInt16(uintptr(comboSecondLevel + 0x8A)) //2 bytes
	//current300Result := current300 + currentgeki // thats not how the game works /shrug
	return current300
}
func CurrentHit50c() int16 {
	proc := proc
	playContainer := playContainerBase
	if osuStatus != 2 {
		return -1
	}

	comboBase, err := proc.ReadUint32(playContainer)
	if err != nil {
		log.Println("CurrentHit50c Base level failure")
		return -2
	}
	comboFirstLevel, err := proc.ReadUint32(uintptr(comboBase))
	if err != nil {
		log.Println("CurrentHit50c First level pointer failure")
		return -3
	}
	comboSecondLevel, err := proc.ReadUint32(uintptr(comboFirstLevel) + 0x38)
	if err != nil {
		log.Println("CurrentHit50c Second level pointer failure")
		return -4
	}
	currentCombo, err := proc.ReadInt16(uintptr(comboSecondLevel + 0x88)) //2 bytes
	if err != nil {
		log.Println("CurrentHitMiss result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentHitMiss() int16 {
	proc := proc
	playContainer := playContainerBase
	if osuStatus != 2 {
		return -1
	}

	comboBase, err := proc.ReadUint32(playContainer)
	if err != nil {
		log.Println("CurrentHitMiss Base level failure")
		return -2
	}
	comboFirstLevel, err := proc.ReadUint32(uintptr(comboBase))
	if err != nil {
		log.Println("CurrentHitMiss First level pointer failure")
		return -2
	}
	comboSecondLevel, err := proc.ReadUint32(uintptr(comboFirstLevel) + 0x38)
	if err != nil {
		log.Println("CurrentHitMiss Second level pointer failure")
		return -2
	}
	currentCombo, err := proc.ReadInt16(uintptr(comboSecondLevel + 0x8E)) //2 bytes
	if err != nil {
		log.Println("CurrentHitMiss result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentScore() int32 {
	proc := proc
	playContainer := playContainerBase
	if osuStatus != 2 {
		return -1
	}

	comboBase, err := proc.ReadUint32(playContainer)
	if err != nil {
		log.Println("CurrentScore Base level failure")
		return -2
	}
	comboFirstLevel, err := proc.ReadUint32(uintptr(comboBase))
	if err != nil {
		log.Println("CurrentScore First level pointer failure")
		return -3
	}
	comboSecondLevel, err := proc.ReadUint32(uintptr(comboFirstLevel) + 0x38)
	if err != nil {
		log.Println("CurrentScore Second level pointer failure")
		return -4
	}
	currentCombo, err := proc.ReadInt32(uintptr(comboSecondLevel + 0x74))
	if err != nil {
		log.Println("CurrentScore result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentGameMode() int32 {
	proc := proc
	playContainer := playContainerBase
	if osuStatus != 2 {
		return -1
	}

	comboBase, err := proc.ReadUint32(playContainer)
	if err != nil {
		log.Println("GameMode Base level failure")
		return -2
	}
	comboFirstLevel, err := proc.ReadUint32(uintptr(comboBase))
	if err != nil {
		log.Println("GameMode First level pointer failure")
		return -3
	}
	comboSecondLevel, err := proc.ReadUint32(uintptr(comboFirstLevel) + 0x38)
	if err != nil {
		log.Println("GameMode Second level pointer failure")
		return -4
	}
	currentCombo, err := proc.ReadInt32(uintptr(comboSecondLevel + 0x64))
	if err != nil {
		log.Println("GameMode result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentAccuracy() float64 {
	proc := proc
	playContainer := playContainerBase
	if osuStatus != 2 {
		return -1
	}

	comboBase, err := proc.ReadUint32(playContainer)
	if err != nil {
		log.Println("Accuracy Base level failure")
		return -2
	}
	comboFirstLevel, err := proc.ReadUint32(uintptr(comboBase))
	if err != nil {
		log.Println("Accuracy First level pointer failure")
		return -3
	}
	comboSecondLevel, err := proc.ReadUint32(uintptr(comboFirstLevel) + 0x48)
	if err != nil {
		log.Println("Accuracy Second level pointer failure")
		return -4
	}
	currentCombo, err := proc.ReadFloat64(uintptr(comboSecondLevel + 0x14))
	if err != nil {
		log.Println("Accuracy result pointer failure")
		return -5
	}
	return currentCombo
}

//monkaW section
func ValidCurrentBeatmapFolderString() string {
	validCurrentBeatmapFolderString := strings.ToValidUTF8(CurrentBeatmapFolderString(), "")
	t := strings.Replace(validCurrentBeatmapFolderString, "\u0000", "", -1)
	strParts := strings.Split(t, "\u0018")
	//fmt.Println(strParts[0])

	if strings.Contains(strParts[0], "@Ou") == true {
		strParts = strings.Split(strParts[0], "@Ou")

	}
	if strings.Contains(strParts[0], " \u000f") == true {
		strParts = strings.Split(strParts[0], " \u000f")

	}
	//if strings.Contains(strParts[0], "?") == true {
	//		strParts = strings.Split(strParts[0], "?")
	//	fmt.Println("?")

	//	}
	if strings.Contains(strParts[0], "W\u000e") == true {
		//fmt.Println("da")
		strParts = strings.Split(strParts[0], "W\u000e")

	}
	if strings.Contains(strParts[0], "\u0001") == true {
		strParts = strings.Split(strParts[0], "\u0001")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "P\u000f") == true {
		strParts = strings.Split(strParts[0], "P\u000f")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "!\u000f") == true {
		strParts = strings.Split(strParts[0], "!\u000f")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "\u001d") == true {
		strParts = strings.Split(strParts[0], "\u001d")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "\u0019") == true {
		strParts = strings.Split(strParts[0], "\u0019")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "-\u000e") == true {
		strParts = strings.Split(strParts[0], "-\u000e")
		//fmt.Println("\u0001 true")

	}

	return strParts[0]
}
func ValidCurrentBeatmapString() string {
	validCurrentBeatmapFolderString := strings.ToValidUTF8(CurrentBeatmapString(), "")
	t := strings.Replace(validCurrentBeatmapFolderString, "\u0000", "", -1)
	strParts := strings.Split(t, "\u0018")
	//fmt.Println(strParts[0])

	if strings.Contains(strParts[0], "@Ou") == true {
		strParts = strings.Split(strParts[0], "@Ou")

	}
	if strings.Contains(strParts[0], " \u000f") == true {
		strParts = strings.Split(strParts[0], " \u000f")

	}
	//if strings.Contains(strParts[0], "?") == true {
	//		strParts = strings.Split(strParts[0], "?")
	//	fmt.Println("?")

	//	}
	if strings.Contains(strParts[0], "W\u000e") == true {
		//fmt.Println("da")
		strParts = strings.Split(strParts[0], "W\u000e")

	}
	if strings.Contains(strParts[0], "\u0001") == true {
		strParts = strings.Split(strParts[0], "\u0001")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "P\u000f") == true {
		strParts = strings.Split(strParts[0], "P\u000f")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "!\u000f") == true {
		strParts = strings.Split(strParts[0], "!\u000f")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "\u001d") == true {
		strParts = strings.Split(strParts[0], "\u001d")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "\u0019") == true {
		strParts = strings.Split(strParts[0], "\u0019")
		//fmt.Println("\u0001 true")

	}
	if strings.Contains(strParts[0], "-\u000e") == true {
		strParts = strings.Split(strParts[0], "\u000e")
		//fmt.Println("\u0001 true")

	}

	return strParts[0]
}
func ValidCurrentBeatmapOsuFileString() string {
	validCurrentBeatmapFolderString := strings.ToValidUTF8(CurrentBeatmapOsuFileString(), "")
	t := strings.Replace(validCurrentBeatmapFolderString, "\u0000", "", -1)
	strParts := strings.Split(t, "\u0018")

	if strings.Contains(strParts[0], "@Ou") == true {
		strParts = strings.Split(strParts[0], "@Ou")

	}
	if strings.Contains(strParts[0], " \u000f") == true {
		strParts = strings.Split(strParts[0], " \u000f")

	}
	if strings.Contains(strParts[0], "\u000f") == true {
		strParts = strings.Split(strParts[0], "\u000f")

	}
	if strings.Contains(strParts[0], "W\u000e") == true {
		strParts = strings.Split(strParts[0], "W\u000e")

	}
	if strings.Contains(strParts[0], "\u0001") == true {
		strParts = strings.Split(strParts[0], "\u0001")
	}
	if strings.Contains(strParts[0], "P\u000f") == true {
		strParts = strings.Split(strParts[0], "P\u000f")

	}
	if strings.Contains(strParts[0], "!\u000f") == true {
		strParts = strings.Split(strParts[0], "!\u000f")

	}
	if strings.Contains(strParts[0], "\u001c") == true {
		strParts = strings.Split(strParts[0], "\u001c")
	}
	if strings.Contains(strParts[0], "\u001d") == true {
		strParts = strings.Split(strParts[0], "\u001d")

	}
	if strings.Contains(strParts[0], "\u0019") == true {
		strParts = strings.Split(strParts[0], "\u0019")

	}
	if strings.Contains(strParts[0], "-\u000e") == true {
		strParts = strings.Split(strParts[0], "-\u000e")

	}

	return strParts[0]
}
