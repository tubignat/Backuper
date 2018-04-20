package main

import (
	"backuper/core/logging"
	"time"

	"github.com/fsnotify/fsnotify"
)

func newSystemWatcher(filenames []string) *SystemWatcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logging.Error("Error has occured: ", err)
		return nil
	}
	return &SystemWatcher{filenames, watcher}
}

type SystemWatcher struct {
	filenames []string
	watcher   *fsnotify.Watcher
}

func (watcher *SystemWatcher) watchAsync(onFileChanged func(filename string), quit chan bool) {
	logging.Debug("Start watching...")
	addFileWatchers(watcher)
	defer watcher.watcher.Close()
	for {
		isClosed := watchAsyncInternal(watcher, onFileChanged, quit)
		if isClosed {
			return
		}
	}
}

func watchAsyncInternal(watcher *SystemWatcher, onFileChanged func(filename string), quit chan bool) bool {
	select {
	case event := <-watcher.watcher.Events:
		select {
		case event = <-watcher.watcher.Events:
		case <-time.After(time.Second):
		}
		logging.Debug("File %s has been changed, operation is %s\n", event.Name, event.Op)
		onFileChanged(event.Name)
		return false
	case err := <-watcher.watcher.Errors:
		logging.Error("Error has occured", err)
		return false
	case <-quit:
		logging.Debug("Stop watching...")
		return true
	}
}

func addFileWatchers(watcher *SystemWatcher) {
	for _, file := range watcher.filenames {
		watcher.watcher.Add(file)
		logging.Debug("File", file, "is being watched now...")
	}
}
