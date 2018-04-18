package main

import (
	"fmt"
)

func main() {
	restart := make(chan bool)
	go watchConfig(func() {
		restart <- true
	})
	go start(restart)
	readCommands()
}

func start(restart chan bool) {
	var watcher *SystemWatcher
	var settings *Settings
	var stop chan bool
	var stub Stub
	for {
		settings = loadSettings()
		Authenticate(settings.ApplicationID)
		fmt.Println("Program started...")

		watcher = newSystemWatcher(settings.Files)
		stop = make(chan bool)
		stub = Stub{}

		go watcher.watchAsync(getBackupFunc(stub), stop)

		<-restart
		stop <- true
		fmt.Println("Restarting...")
	}
}

func getBackupFunc(apiClient ApiClient) func(filename string) {
	return func(filename string) {
		apiClient.Backup(filename)
	}
}

func readCommands() {
	for {
		command := ""
		fmt.Scanln(&command)
		if command == "exit" {
			return
		}
	}
}
