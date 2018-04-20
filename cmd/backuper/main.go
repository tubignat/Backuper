package main

import (
	"backuper/core/api"
	"backuper/core/settings"
	"fmt"
	"log"
)

func main() {
	go start()
	readCommands()
}

func start() {
	log.Print("Program started...")
	keeper := settings.NewKeeper("config.json")
	settings := keeper.GetRelevantSettings()
	yandexClient := api.NewYandexApiClient(settings.Yandex.ApplicationID)
	watcher := newSystemWatcher(settings.Files)
	stop := make(chan bool)
	backupFunc := getBackupFunc(yandexClient)

	go watcher.watchAsync(backupFunc, stop)
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
