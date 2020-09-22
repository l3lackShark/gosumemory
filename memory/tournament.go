package memory

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//initTournement should be called on tournament manager
func initTournement() error {

	//read tournament.cfg to check how many clients we are expecting
	osuExecutablePath, err := process.ExecutablePath()
	if err != nil {
		return err
	}
	if !strings.Contains(osuExecutablePath, `:\`) {
		log.Println("Automatic executable path finder has failed. The program will now ext. GOT: ", osuExecutablePath)
		return errors.New("osu! executable was not found")
	}
	cfgFile, err := os.Open(filepath.Join(filepath.Dir(osuExecutablePath), "tournament.cfg"))
	if err != nil {
		return err
	}
	defer cfgFile.Close()
	scanner := bufio.NewScanner(cfgFile)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "TeamSize") {
			teamSize, err := strconv.Atoi(scanner.Text()[len(scanner.Text())-1:])
			if err != nil {
				return err
			}
			fmt.Println("Total expected amount of tournament clients:", teamSize*2)

		}
	}
	return nil
}
