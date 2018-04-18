package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func newSystemWatcher(filenames []string) *SystemWatcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
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
		watchAsyncInternal(watcher, onFileChanged, quit)
	}
}

func watchAsyncInternal(watcher *SystemWatcher, onFileChanged func(filename string), quit chan bool) {
	select {
	case event := <-watcher.watcher.Events:
		fmt.Printf("File %s has been changed, operation is %s\n", event.Name, event.Op)
		onFileChanged(event.Name)
	case err := <-watcher.watcher.Errors:
		fmt.Println("Error has occured", err)
	case <-quit:
		fmt.Println("Stop watching...")
		break
	}
}

func addFileWatchers(watcher *SystemWatcher) {
	for _, file := range watcher.filenames {
		watcher.watcher.Add(file)
		fmt.Printf("File %s is being watched now...\n", file)
	}
}
