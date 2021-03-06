package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/l3lackShark/config"
)

//Config file
var Config map[string]string

//Init the config file
func Init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	cfg, err := config.SetFile(filepath.Join(exPath, "config.ini"))
	if err == config.ErrDoesNotExist {
		d := []byte(`[Main]
update = 100
path = auto
cgodisable = false
memdebug = false
memcycletest = false
wine = false
		
[Web]
serverip = 127.0.0.1:24050
cors = false
		
[GameOverlay] ; https://github.com/l3lackShark/gosumemory/wiki/GameOverlay
enabled = false
gameWidth = 1920
gameHeight = 1080
overlayURL = http://localhost:24050/InGame2
overlayWidth = 355
overlayHeight = 97
overlayOffsetX = 0
overlayOffsetY = 0
overlayScale = 10
		
`)
		if err := ioutil.WriteFile(filepath.Join(exPath, "config.ini"), d, 0644); err != nil {
			panic(err)
		}
		cfg, err = config.SetFile(filepath.Join(exPath, "config.ini"))
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		log.Fatalln(err)
	}
	Config, err = cfg.Parse()
	if err != nil {
		panic(err)
	}
	if Config["overlayURL"] == "" { //Quck hack to append GameOverlay stuff to existing config, whole system needs revamp
		file, err := os.OpenFile(filepath.Join(exPath, "config.ini"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(fmt.Sprintf("\n[GameOverlay]; https://github.com/l3lackShark/gosumemory/wiki/GameOverlay\nenabled = false\ngameWidth = 1920\ngameHeight = 1080\noverlayURL = http://localhost:24050/InGame2\noverlayWidth = 355\noverlayHeight = 97\noverlayOffsetX = 0\noverlayOffsetY = 0\noverlayScale = 10"))
		if err != nil {
			panic(err)
		}

		Init()
	}
}
