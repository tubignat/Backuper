package settings

import (
	"backuper/common"
	"encoding/json"
	"log"

	"github.com/fsnotify/fsnotify"
)

type Keeper struct {
	ConfigFileName string
	Settings       *Settings
	HasChanges     bool
}

func NewKeeper(configFilename string) *Keeper {
	return &Keeper{
		ConfigFileName: configFilename,
		Settings:       nil,
		HasChanges:     false,
	}
}

func (keeper *Keeper) GetRelevantSettings() *Settings {
	if keeper.Settings == nil || keeper.HasChanges {
		keeper.Settings = loadSettings(keeper.ConfigFileName)
	}
	return keeper.Settings
}

func (keeper *Keeper) TrackConfig(filename string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go watchFileRoutine(watcher, &keeper.HasChanges)
}

func watchFileRoutine(watcher *fsnotify.Watcher, hasChanges *bool) {
	for {
		select {
		case <-watcher.Events:
			*hasChanges = true
		case err := <-watcher.Errors:
			log.Panic(err)
		}
	}
}
func loadSettings(filename string) *Settings {
	var settings Settings
	content := common.ReadFile(filename)
	if error := json.Unmarshal(*content, &settings); error != nil {
		log.Panic(error)
	}
	return &settings
}
