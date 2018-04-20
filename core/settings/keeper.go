package settings

import (
	"backuper/core/common"
	"backuper/core/logging"
	"encoding/json"

	"github.com/fsnotify/fsnotify"
)

type Keeper struct {
	ConfigFileName string
	Settings       *Settings
	HasChanges     bool
}

func NewKeeper(configFilename string) *Keeper {
	keeper := &Keeper{
		ConfigFileName: configFilename,
		Settings:       nil,
		HasChanges:     false,
	}
	keeper.TrackConfig(configFilename)
	return keeper
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
		logging.Error("Failed to create watcher ", err)
		return
	}
	go watchFileRoutine(watcher, &keeper.HasChanges)
}

func watchFileRoutine(watcher *fsnotify.Watcher, hasChanges *bool) {
	for {
		select {
		case <-watcher.Events:
			*hasChanges = true
			logging.Debug("Changes in the config file noted")
		case err := <-watcher.Errors:
			logging.Error("Watcher throwed an error ", err)
		}
	}
}
func loadSettings(filename string) *Settings {
	logging.Debug("Start loading settings from ", filename, "...")
	var settings Settings
	content := common.ReadFile(filename)
	if error := json.Unmarshal(*content, &settings); error != nil {
		logging.Error(error, "Failed to load settings... ")
	}
	logging.Debug("Settings loaded successfully: ", settings)
	return &settings
}
