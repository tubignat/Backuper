package api

import "fmt"

type Stub struct {
}

func (stub Stub) Backup(filename string) BackupResult {
	fmt.Println("Backuping...")
	return BackupResult{Success}
}
