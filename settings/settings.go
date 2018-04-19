package settings

import (
	"time"
)

type Settings struct {
	Files  []string
	Yandex YandexDiskSettings
}

type YandexDiskSettings struct {
	ApplicationID string
	Password      string
	Token         string
	Expires       time.Time
}
