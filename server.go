package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Andoryuuta/kiwi"
	"github.com/gorilla/websocket"
	parser "github.com/natsukagami/go-osu-parser"
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
var baseDir string = "/home/blackshark/drives/ps3drive/osu!/Songs"
var playTimeBase uintptr
var playTime uintptr
var currentBeatmapDataBase uint32
var currentBeatmapDataFirtLevel uint32
var playContainerBaseAddr uint32
var playContainerFirstlevel uint32
var playContainer38 uint32
var fullPathToOsu string
var osuFileStdIN string
var minBPM float64
var maxBPM float64
var ourTime []int
var lastObjectInt int
var lastObject string

func Cmd(cmd string, shell bool) []byte {

	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			println("some error found", err)
		}
		return out
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			println("some error found2", err)
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
		log.Fatalln("osu! is probably not fully loaded, please load the game up and try again!")
	}
	output := (input[3:])
	yosuBase := firstN(output, 8)
	check := strings.Contains(yosuBase, ",")
	if check == true {
		yosuBase = strings.Replace(yosuBase, ",", "", -1)
	}
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
	check := strings.Contains(yosuBase, ",")
	if check == true {
		yosuBase = strings.Replace(yosuBase, ",", "", -1)
	}
	osuBaseString := "0x" + yosuBase
	osuBaseUINT32 := cast.ToUint32(osuBaseString)
	osuBase = uintptr(osuBaseUINT32)
	//println(CurrentBeatmapFolderString())
	if osuBase == 0 {
		log.Fatalln("Could not find OsuBaseAddr, is osu! running?")
	}
	return osuBase
}

func OsuPlayTimeAddr() uintptr { //in hopes to deprecate this
	x := Cmd("scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;5E 5F 5D C3 A1 ?? ?? ?? ?? 89 ?? 04;list;exit'", true)
	outStr := cast.ToString(x)
	outStr = strings.Replace(outStr, " ", "", -1)

	input := outStr
	if input == "" {
		log.Fatalln("OsuBase addr fail")
	}
	output := (input[3:])
	yosuBase := firstN(output, 8)
	check := strings.Contains(yosuBase, ",")
	if check == true {
		yosuBase = strings.Replace(yosuBase, ",", "", -1)
	}
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
	check := strings.Contains(yosuBase, ",")
	if check == true {
		yosuBase = strings.Replace(yosuBase, ",", "", -1)
	}
	osuBaseString := "0x" + yosuBase

	osuBaseUINT32 := cast.ToUint32(osuBaseString)
	osuBase = uintptr(osuBaseUINT32)
	if osuBase == 0 {
		log.Fatalln("Could not find osuplayContainer address, is osu! running?", "Address was:", osuBaseString)
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
	playTimeBase = OsuPlayTimeAddr()
	playContainer = OsuplayContainer()
	playContainerBase = (playContainer - 0x4)
	playTime = (playTimeBase + 0x5)

	if err != nil {
		log.Fatalln("is osu! running? (osu! status offset was not found)")
	}
	var tempCurrentBeatmapOsu string

	for {

		var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
		for procerr != nil {
			log.Println("is osu! running? (osu! process was not found, terminating...)")
			proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
		}
		osuStatusValue, err := proc.ReadUint16(uintptrOsuStatus)
		if err != nil {
			log.Println("osu! status could not be found...", err)

		}
		osuStatus = osuStatusValue
		//init base stuff
		currentBeatmapDataBase, err = proc.ReadUint32(currentBeatmapData)
		if err != nil {
			log.Println("currentBeatmapDataBase Base level failure")
		}
		currentBeatmapDataFirtLevel, err = proc.ReadUint32(uintptr(currentBeatmapDataBase))
		if err != nil {
			log.Println("currentBeatmapDataFirtLevel First level pointer failure")
		}
		if osuStatus == 2 {
			playContainerBaseAddr, err = proc.ReadUint32(playContainerBase)
			if err != nil {
				log.Println("playContainerBaseAddr Base level failure")
			}
			playContainerFirstlevel, err = proc.ReadUint32(uintptr(playContainerBaseAddr))
			if err != nil {
				log.Println("playContainerFirstlevel pointer failure")
			}
			playContainer38, err = proc.ReadUint32(uintptr(playContainerFirstlevel) + 0x38)
			if err != nil {
				log.Println("playContainer38 pointer failure")

			}
		}

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
			CurrentPlayTime             int32   `json:"bmCurrentTime"`
			MinBPM                      float64 `json:"bmMinBPM"`
			MaxBPM                      float64 `json:"bmMaxBPM"`
		}

		type EverythingInMenu2 struct { //order sets here
			D EverythingInMenu `json:"menuContainer"`
			P PlayContainer    `json:"gameplayContainer"`
		}

		PlayContainerStruct := PlayContainer{
			CurrentHit300c:     CurrentHit300c(),
			CurrentHit100c:     CurrentHit100c(),
			CurrentHit50c:      CurrentHit50c(),
			CurrentHitMiss:     CurrentHitMiss(),
			CurrentScore:       CurrentScore(),
			CurrentAccuracy:    CurrentAccuracy(),
			CurrentCombo:       CurrentCombo(),
			CurrentGameMode:    CurrentGameMode(),
			CurrentAppliedMods: CurrentAppliedMods(),
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
			CurrentBeatmapFolderString:  ValidCurrentBeatmapFolderString(),
			CurrentBeatmapOsuFileString: ValidCurrentBeatmapOsuFileString(),
			CurrentBeatmapAR:            CurrentBeatmapAR(),
			CurrentBeatmapOD:            CurrentBeatmapOD(),
			CurrentBeatmapCS:            CurrentBeatmapCS(),
			CurrentBeatmapHP:            CurrentBeatmapHP(),
			CurrentPlayTime:             CurrentPlayTime(),

			MinBPM: minBPM,
			MaxBPM: maxBPM,
		}
		if osuStatusValue == 2 {

			//fmt.Println(ourTime[0])
			for _, hitObjectTime := range ourTime {

				if int32(hitObjectTime) >= MenuContainerStruct.CurrentPlayTime {

					lastObjectInt = SliceIndex(len(ourTime), func(i int) bool { return ourTime[i] == hitObjectTime })
					lastObject = cast.ToString(lastObjectInt)
					fmt.Println(PP())
					//fmt.Println(hitObjectTime)
					break

				}
			}

		}

		if MenuContainerStruct.CurrentBeatmapOsuFileString != tempCurrentBeatmapOsu {
			ourTime = nil

			tempCurrentBeatmapOsu = MenuContainerStruct.CurrentBeatmapOsuFileString
			fullPathToOsu = fmt.Sprintf(baseDir + "/" + MenuContainerStruct.CurrentBeatmapFolderString + "/" + MenuContainerStruct.CurrentBeatmapOsuFileString)

			j, err := ioutil.ReadFile(fullPathToOsu)
			if err != nil {
				fmt.Println("osu file was not found2")
			}
			osuFileStdIN = string(j)
			if strings.Contains(osuFileStdIN, "[HitObjects]") == true {
				splitted := strings.Split(osuFileStdIN, "[HitObjects]")[1]
				newline := strings.Split(splitted, "\n")

				for i := 1; i < len(newline)-1; i++ {
					if len(newline[i]) > 0 {
						elements := strings.Split(newline[i], ",")[2]
						elementsInt := cast.ToInt(elements)
						ourTime = append(ourTime, elementsInt)

					}
				}
			}

			if strings.HasSuffix(fullPathToOsu, ".osu") == true {
				//fmt.Println(fullPathToOsu)
				file, err := os.Open(fullPathToOsu)
				if err != nil {
					log.Println(err, "in error")
					defer file.Close()
				}
				defer file.Close()
				scanner := bufio.NewScanner(file)
				var bgString string
				for scanner.Scan() {

					//fmt.Println(scanner.Text())
					if strings.Contains(scanner.Text(), ".jpg") == true {
						bg := strings.Split(scanner.Text(), "\"")
						bgString = (bg[1])
						break
						//log.Fatalln(scanner.Text())
					}
					if strings.Contains(scanner.Text(), ".png") == true {
						bg := strings.Split(scanner.Text(), "\"")
						bgString = (bg[1])
						//log.Fatalln(scanner.Text())
						break
					}
					if strings.Contains(scanner.Text(), ".JPG") == true {
						bg := strings.Split(scanner.Text(), "\"")
						bgString = (bg[1])
						//log.Fatalln(scanner.Text())
						break
					}
					if strings.Contains(scanner.Text(), ".PNG") == true {
						bg := strings.Split(scanner.Text(), "\"")
						bgString = (bg[1])
						break
						//log.Fatalln(scanner.Text())
					} else {
						bgString = ""
					}
				}
				if err := scanner.Err(); err != nil {
					log.Println(err)
				}
				var fullPathToBG string = fmt.Sprintf(baseDir + "/" + MenuContainerStruct.CurrentBeatmapFolderString + "/" + bgString)
				var fullBGCommand string = fmt.Sprintf("ln -nsf " + "\"" + fullPathToBG + "\"" + " " + "$PWD" + "/bg.png")
				Cmd((fullBGCommand), true)

				minBPM = CurrentBeatmapMinBPM()
				maxBPM = CurrentBeatmapMaxBPM()
				//fmt.Println(OsuHitobjects())
			} else {
				fmt.Println("osu file was not found")
			}

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

func CurrentBeatmapID() uint32 { //currentbeatmapdata

	currentBeatmapID, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0xC4))
	if err != nil {
		log.Println("CurrentBeatmapID result pointer failure")
		return 0
	}
	return currentBeatmapID
}
func CurrentBeatmapSetID() uint32 {
	currentSetBeatmapID, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0xC8))
	if err != nil {
		log.Println("CurrentBeatmapSetID result pointer failure")
		return 0
	}
	return currentSetBeatmapID
}
func CurrentBeatmapString() string {
	beatmapStringSecondLevel, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0x7C))
	if err != nil {
		log.Println("BeatMapString Second level pointer failure")
		return "-4"
	}
	beatmapStringSize, err := proc.ReadUint32(uintptr(beatmapStringSecondLevel + 0x4))
	if err != nil {
		log.Println("BeatMapString Third level pointer failure")
		return "-6"
	}

	beatmapStringResult, err := proc.ReadBytes(uintptr(beatmapStringSecondLevel+0x8), cast.ToInt(beatmapStringSize*2)) // fix repeating
	if err != nil {
		log.Println("BeatMapString Third level pointer failure")
		return "-5"
	}
	beatmapString := string(beatmapStringResult)
	beatmapValidString := strings.ToValidUTF8(beatmapString, "")
	return beatmapValidString
}

func CurrentBeatmapFolderString() string {

	beatmapFolderStringSecondLevel, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0x74))
	if err != nil {
		log.Println("BeatMapFolderString Second level pointer failure")
		return "-4"
	}

	beatmapFolderStringSize, err := proc.ReadUint32(uintptr(beatmapFolderStringSecondLevel + 0x4))
	if err != nil {
		log.Println("BeatMapFolderStringSize Third level pointer failure")
		return "-6"
	}

	beatmapStringResult, err := proc.ReadBytes(uintptr(beatmapFolderStringSecondLevel+0x8), cast.ToInt(beatmapFolderStringSize*2)) // fix repeating
	if err != nil {
		log.Println("BeatMapFolderString Third level pointer failure")
		return "-5"
	}
	beatmapString := string(beatmapStringResult)
	beatmapValidString := strings.ToValidUTF8(beatmapString, "")
	return beatmapValidString
}
func CurrentBeatmapOsuFileString() string {

	beatmapFolderStringSecondLevel, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0x8C))
	if err != nil {
		log.Println("BeatMapOsuFileString Second level pointer failure")
		return "-4"
	}

	beatmapStringSize, err := proc.ReadUint32(uintptr(beatmapFolderStringSecondLevel + 0x4))
	if err != nil {
		log.Println("BeatMapString Third level pointer failure")
		return "-6"
	}

	beatmapStringResult, err := proc.ReadBytes(uintptr(beatmapFolderStringSecondLevel+0x8), cast.ToInt(beatmapStringSize*2))
	if err != nil {
		log.Println("BeatMapString Third level pointer failure")
		return "-5"
	}
	beatmapString := string(beatmapStringResult)
	beatmapValidString := strings.ToValidUTF8(beatmapString, "")
	return beatmapValidString
}
func CurrentBeatmapAR() float32 {
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(currentBeatmapDataFirtLevel + 0x2C))
	if err != nil {
		log.Println("AR result level pointer failure")
		return -5
	}
	return currentSetBeatmapID
}
func CurrentBeatmapCS() float32 {
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(currentBeatmapDataFirtLevel + 0x30))
	if err != nil {
		log.Println("CS result level pointer failure")
		return -4
	}
	return currentSetBeatmapID
}
func CurrentBeatmapHP() float32 {
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(currentBeatmapDataFirtLevel + 0x34))
	if err != nil {
		log.Println("HP result level pointer failure")
		return -5
	}
	return currentSetBeatmapID
}
func CurrentBeatmapOD() float32 {
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(currentBeatmapDataFirtLevel + 0x38))
	if err != nil {
		log.Println("OD result level pointer failure")
		return -5
	}
	return currentSetBeatmapID
}

// ------------------- PlayContainer
func CurrentAppliedMods() int32 {
	if osuStatus != 2 {
		return -1
	}
	currentCombo, err := proc.ReadInt32(uintptr(playContainer38 + 0x1C))
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
	val := xorVal2 ^ xorVal1
	return val
}
func CurrentCombo() int32 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt32(uintptr(playContainer38 + 0x90))
	if err != nil {
		log.Println("CurrentCombo result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentHit100c() int16 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt16(uintptr(playContainer38 + 0x84)) //2 bytes
	if err != nil {
		log.Println("CurrentHit100c result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentHit300c() int16 {
	if osuStatus != 2 {
		return -1
	}
	current300, err := proc.ReadInt16(uintptr(playContainer38 + 0x86)) //2 bytes
	if err != nil {
		log.Println("CurrentHit300c result pointer failure")
		return -5
	}
	//currentgeki, err := proc.ReadInt16(uintptr(comboSecondLevel + 0x8A)) //2 bytes
	//current300Result := current300 + currentgeki // thats not how the game works /shrug
	return current300
}
func CurrentHit50c() int16 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt16(uintptr(playContainer38 + 0x88)) //2 bytes
	if err != nil {
		log.Println("CurrentHitMiss result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentHitMiss() int16 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt16(uintptr(playContainer38 + 0x8E)) //2 bytes
	if err != nil {
		log.Println("CurrentHitMiss result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentScore() int32 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt32(uintptr(playContainer38 + 0x74))
	if err != nil {
		log.Println("CurrentScore result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentGameMode() int32 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt32(uintptr(playContainer38 + 0x64))
	if err != nil {
		log.Println("GameMode result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentAccuracy() float64 {
	if osuStatus != 2 {
		return -1
	}

	comboSecondLevel, err := proc.ReadUint32(uintptr(playContainerFirstlevel) + 0x48)
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

//not so monkaW section
func ValidCurrentBeatmapFolderString() string {
	validCurrentBeatmapFolderString := strings.ToValidUTF8(CurrentBeatmapFolderString(), "")
	t := strings.Replace(validCurrentBeatmapFolderString, "\u0000", "", -1)
	strParts := strings.Split(t, "\u0018")

	return strParts[0]
}
func ValidCurrentBeatmapString() string {
	validCurrentBeatmapFolderString := strings.ToValidUTF8(CurrentBeatmapString(), "")
	t := strings.Replace(validCurrentBeatmapFolderString, "\u0000", "", -1)
	strParts := strings.Split(t, "\u0018")

	return strParts[0]
}
func ValidCurrentBeatmapOsuFileString() string {
	validCurrentBeatmapFolderString := strings.ToValidUTF8(CurrentBeatmapOsuFileString(), "")
	t := strings.Replace(validCurrentBeatmapFolderString, "\u0000", "", -1)
	strParts := strings.Split(t, "\u0018")

	if strings.Contains(strParts[0], ".osu") == true {
		strParts = strings.Split(strParts[0], ".osu")
		strPartsString := cast.ToString(strParts[0])
		strPartsString = strPartsString + ".osu"
		strParts[0] = strPartsString

	}

	return strParts[0]
}
func CurrentPlayTime() int32 {
	playTimeFirstLevel, err := proc.ReadUint32(playTime)
	if err != nil {
		log.Println("playTime Base level failure")
		return 0
	}
	playTimeValue, err := proc.ReadUint32(uintptr(playTimeFirstLevel))
	if err != nil {
		log.Println("playTime Result level failure")
		return 0
	}

	return cast.ToInt32(playTimeValue)
}
func CurrentBeatmapMinBPM() float64 {
	osuBeatmap, err := parser.ParseFile(fullPathToOsu)
	if err != nil {
		fmt.Println("CurrentBeatmapMinBPM error")
	}
	result := osuBeatmap.BpmMin
	return result
}
func CurrentBeatmapMaxBPM() float64 {
	osuBeatmap, err := parser.ParseFile(fullPathToOsu)
	if err != nil {
		fmt.Println("CurrentBeatmapMaxBPM error")
	}
	result := osuBeatmap.BpmMax
	return result
}
func CurrentBeatmapIntro() int32 {
	osuBeatmap, err := parser.ParseFile(fullPathToOsu)
	if err != nil {
		fmt.Println("CurrentBeatmapMaxBPM error")
	}
	result := cast.ToInt32(osuBeatmap.PreviewTime)
	return result
}
func PP() string {
	//fmt.Println(stdin)
	calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+"-end"+lastObject, true)
	return cast.ToString(calc)
}
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
