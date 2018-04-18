package main

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
)

func newSystemWatcher(filenames []string) *SystemWatcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error has occured: ", err)
	}
	return &SystemWatcher{filenames, watcher}
}

type SystemWatcher struct {
	filenames []string
	watcher   *fsnotify.Watcher
}

func (watcher *SystemWatcher) watchAsync(onFileChanged func(filename string), quit chan bool) {
	fmt.Println("Start watching...")
	addFileWatchers(watcher)
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
		fmt.Printf("File %s has been changed, operation is %s\n", event.Name, event.Op)
		onFileChanged(event.Name)
		return false
	case err := <-watcher.watcher.Errors:
		fmt.Println("Error has occured", err)
		return false
	case <-quit:
		watcher.watcher.Close()
		fmt.Println("Stop watching...")
		return true
	}
}

func addFileWatchers(watcher *SystemWatcher) {
	for _, file := range watcher.filenames {
		watcher.watcher.Add(file)
		fmt.Printf("File %s is being watched now...\n", file)
	}
}
