package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	configFilename = "config.json"
)

type Settings struct {
	ApplicationID string
	Files         []string
}

func loadSettings() *Settings {
	var settings *Settings
	content, error := ioutil.ReadFile(configFilename)
	if error != nil {
		fmt.Println(error)
	}
	error = json.Unmarshal(content, &settings)
	if error != nil {
		fmt.Println(error)
	}
	return settings
}

func watchConfig(onConfigChanged func()) {
	watcher, _ := fsnotify.NewWatcher()
	watcher.Add(configFilename)
	for {
		waitForChanges(watcher)
		onConfigChanged()
	}
}

// This is a crutch for the Windows's FileSystemWatcher - sometimes events raise twice
// Here is an explanation:
// https://stackoverflow.com/questions/36563396/how-to-populate-an-array-with-channels-in-go
func waitForChanges(watcher *fsnotify.Watcher) {
	<-watcher.Events
	select {
	case <-watcher.Events:
	case <-time.After(time.Second):
	}
}
