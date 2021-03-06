package injctr

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

//Injct dll into osu's process
func Injct(pid int) error {
	if runtime.GOOS != "windows" {
		return errors.New("Gameoverlay only works under windows")
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	_, err = os.Stat("gameoverlay")
	if err != nil {
		fmt.Println("[GAMEOVERLAY] Downloading gameoverlay... (can take a while, filesize is around 60MB)")
		err = downloadFile("https://omk.pics/12/wVugx", "gameoverlay.zip")
		if err != nil {
			return err
		}
		unzip("gameoverlay.zip", "gameoverlay")
		err = os.Remove("gameoverlay.zip")
		if err != nil {
			return err
		}
	}

	_, err = os.Stat("gameoverlay\\gosumemoryoverlay.dll")
	if err != nil {
		return err
	}
	_, err = exec.Command("gameoverlay\\a.exe", strconv.Itoa(pid), filepath.Join(exPath, "gameoverlay", "gosumemoryoverlay.dll")).Output()
	if err != nil {
		return err
	}
	fmt.Println("[GAMEOVERLAY] Initialized successfully, see https://github.com/l3lackShark/gosumemory/wiki/GameOverlay for tutorial")
	return nil
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
