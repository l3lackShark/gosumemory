package pp

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/tcolgate/mp3"
)

func hashFileMD5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}
func calculateMP3Time() (int32, error) {
	if !strings.HasSuffix(strings.ToLower(memory.MenuData.Bm.Path.FullMP3Path), ".mp3") {
		pp.Println("Expected mp3, got something else. Aborting mp3 time calculation. GOT: ", memory.MenuData.Bm.Path.FullMP3Path)
		return 0, nil
	}
	var t int64
	r, err := os.Open(memory.MenuData.Bm.Path.FullMP3Path)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {

		if err := d.Decode(&f, &skipped); err != nil {
			if err != nil {
				break
			}
		}

		t = t + f.Duration().Milliseconds()
	}

	return int32(t), nil
}

func minMax(array []int) (int, int) {
	if len(array) < 1 {
		return 0, 0
	}
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func calculateAccuracy(h300 float32, h100 float32, h50 float32, h0 float32) float32 {
	return 100 * (h50*50 + h100*100 + h300*300) / (h50*300 + h100*300 + h300*300 + h0*300)
}
