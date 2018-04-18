package api

const (
	Success Status = 1
	Fail    Status = 2
)

type ApiClient interface {
	Backup(filename string) BackupResult
}

type BackupResult struct {
	Status Status
}

type Status int
