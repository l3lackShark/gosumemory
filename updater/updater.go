package updater

import (
	"log"
	"os"
	"path/filepath"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/skratchdot/open-golang/open"
)

const version = "1.2.5"

//DoSelfUpdate updates the application
func DoSelfUpdate() {
	name, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	v := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(v, "l3lackShark/cpol-replays")
	if err != nil {
		log.Println("Binary update failed:", err)
		return
	}
	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		log.Println("Current binary is the latest version", version)
		full, _ := os.Executable()
		path, executable := filepath.Split(full)
		oldName := filepath.Join(path, "."+executable+".old")
		os.Remove(oldName)
	} else {
		log.Println("Successfully updated to version", latest.Version)
		log.Println("Release note:\n", latest.ReleaseNotes)

		open.Start(name)
		os.Exit(0)
	}
}
