package settings

type Settings struct {
	Files  []string
	Yandex YandexDiskSettings
}

type YandexDiskSettings struct {
	ApplicationID string
	Password      string
}
