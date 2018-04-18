package main

import (
	"fmt"
)

func main() {
	fmt.Println("Program started...")
	watcher := newSystemWatcher([]string{"C:\\Connector\\stuff\\file"})
	quit := make(chan bool)
	stub := stub{}
	go watcher.watchAsync(getBackupFunc(stub), quit)
	readCommands(quit)
}

func getBackupFunc(apiClient ApiClient) func(filename string) {
	return func(filename string) {
		apiClient.Backup(filename)
	}
}

func readCommands(quit chan bool) {
	for {
		command := ""
		fmt.Scanln(&command)
		if command == "exit" {
			quit <- true
			return
		}
	}
}
