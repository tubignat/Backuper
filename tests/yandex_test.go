package tests

import (
	"backuper/core/api"
	"backuper/core/common"
	"testing"
)

const (
	applicationID = "201581d56b834eaa8f4b59942be5f7be"
)

// TODO: automatically move yandex_token file into tests directory after authentication succeeded
func TestYandex(t *testing.T) {
	filename := "file.txt"
	text := []byte("Hello! this is a test file")
	common.WriteFile(filename, &text)
	client := api.NewYandexApiClient(applicationID)
	client.Backup(filename)
	// TODO: add assertion that file has successfully uploaded to Yandex.Disk
}
