//+build windows

package db

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/MakeNowJust/hotkey"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/skratchdot/open-golang/open"
)

func writeOSR() {
	hkey := hotkey.New()
	for {
		hkey.Register(0, hotkey.F2, func() {
			if memory.DynamicAddresses.IsReady == true && memory.GameplayData.GameMode != 3 && memory.GameplayData.IsFailed == 1 && tempBeatmapFailTime != memory.GameplayData.FailTime && memory.GameplayData.FailTime != 0 {
				tempBeatmapFailTime = memory.GameplayData.FailTime
				fmt.Println("Writing replay file...")
				name := fmt.Sprintf("FailedReplays/%s - %s - %s [%s] (%s) %s.osr", memory.GameplayData.Name+"(Failed)", memory.MenuData.Bm.Metadata.Artist, memory.MenuData.Bm.Metadata.Title, memory.MenuData.Bm.Metadata.Version, strings.ReplaceAll(time.Now().Format(time.RFC1123), ":", "-"), gamemodeToStr(memory.GameplayData.GameMode))
				fmt.Println(name)
				file, err := os.Create(name)
				if err != nil {
					fmt.Println(err)
				}
				replayWriter := bufio.NewWriter(file)
				OsrStruct := convertMemoryDataToOSRStruct()
				v := reflect.ValueOf(OsrStruct)
				values := make([]interface{}, v.NumField())
				for i := 0; i < v.NumField(); i++ {
					values[i] = v.Field(i).Interface()
					switch v.Field(i).Kind() {
					case reflect.String:
						replayWriter.WriteByte(0x0B) //please never exceed 255 (TODO: proper strings handler)
						replayWriter.WriteByte(byte(len(v.Field(i).String())))
						replayWriter.WriteString(v.Field(i).String())
					case reflect.Uint8:
						writeUint8(replayWriter, uint8(v.Field(i).Uint()))
					case reflect.Uint16:
						writeUint16(replayWriter, uint16(v.Field(i).Uint()))
					case reflect.Int32:
						writeInt32(replayWriter, int32(v.Field(i).Int()))
					case reflect.Bool:
						writeBool(replayWriter, v.Field(i).Bool())
					case reflect.Int64:
						writeInt64(replayWriter, v.Field(i).Int())
					case reflect.Slice:
						replayWriter.Write(v.Field(i).Bytes())
					default:
						log.Fatalln("Unsupported struct type!")
					}

				}
				replayWriter.Flush()
				file.Close()
				fmt.Println("Finished writing replay file!")
				err = open.Start(filepath.Join(name))
				if err != nil {
					fmt.Println("Replay open err: ", err)
				}
			}

		})
		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}
}
