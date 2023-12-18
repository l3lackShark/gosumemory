package updater

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/skratchdot/open-golang/open"
)

const version = "1.3.9"

// DoSelfUpdate updates the application
func DoSelfUpdate() {
	fmt.Println("Checking Updates... (can take some time if you have bad routing to GitHub)")
	name, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	v := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(v, "l3lackShark/gosumemory")
	if err != nil {
		log.Println("Binary update failed:", err)
		return
	}
	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		log.Println("Current binary is the latest version", version)
		log.Println("Release notes:\n", latest.ReleaseNotes)
		full, _ := os.Executable()
		path, executable := filepath.Split(full)
		oldName := filepath.Join(path, "."+executable+".old")
		os.Remove(oldName)
	} else {
		log.Println("Successfully updated to version", latest.Version)
		log.Println("Release notes:\n", latest.ReleaseNotes)
		time.Sleep(3 * time.Second)
		open.Start(name)
		os.Exit(0)
	}
}
