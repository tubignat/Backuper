package main

import (
	"backuper/api"
	"backuper/settings"
	"fmt"
	"time"
)

func main() {
	start()
	readCommands()
}

func start() {
	keeper := settings.NewKeeper("config.json")
	fmt.Println("Program started...")
	settings := keeper.GetRelevantSettings()
	if settings.Yandex.Token == "" || settings.Yandex.Expires.Before(time.Now()) {
		api.Authenticate(settings.Yandex.ApplicationID)
	}

	watcher = newSystemWatcher(settings.Files)
	stop = make(chan bool)
	stub = api.Stub{}

	go watcher.watchAsync(getBackupFunc(stub), stop)
}

func getBackupFunc(apiClient api.Client) func(filename string) {
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
