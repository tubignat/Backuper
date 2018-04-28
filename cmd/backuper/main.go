package main

import (
	"backuper/core/api"
	"backuper/core/logging"
	"backuper/core/settings"
	"fmt"
	"os/exec"
)

const (
	help = "Available commands:\n    config - opens a config file in notepad\n    exit - closes the program\n"
)

func main() {
	logging.Configure("logs.txt")
	go start()
	readCommands()
}

func start() {
	for {
		logging.Info("Configuring...")

		restart := make(chan bool)
		keeper := settings.NewKeeper("config.json", restart)
		settings := keeper.GetRelevantSettings()
		yandexClient := api.NewYandexAPIClient(settings.Yandex.ApplicationID)
		watcher := newSystemWatcher(settings.Files)
		stop := make(chan bool)
		backupFunc := getBackupFunc(yandexClient)

		go watcher.watchAsync(backupFunc, stop)

		logging.Info("Program started...")

		<-restart
		stop <- true

		logging.Info("Restarting...")
	}
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
		switch command {
		case "help":
			fmt.Print(help)
		case "config":
			if error := exec.Command("notepad", "config.json").Start(); error != nil {
				logging.Error(error)
			}
		case "exit":
			return
		}
	}
}
