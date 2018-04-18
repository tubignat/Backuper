package main

const (
	Success Status = 1
	Fail    Status = 2
)

type ApiClient interface {
	Fetch(filename string) FetchResult
	Backup(filename string) BackupResult
}

type FetchResult struct {
	Status Status
}

type BackupResult struct {
	Status Status
}

type Status int
