package main

import "fmt"

type stub struct {
}

func (stub stub) Fetch(filename string) FetchResult {
	fmt.Println("Fetching...")
	return FetchResult{Success}
}

func (stub stub) Backup(filename string) BackupResult {
	fmt.Println("Backuping...")
	return BackupResult{Success}
}
