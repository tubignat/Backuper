package api

import (
	"os/exec"
	"time"
)

const (
	YandexOAuthURL = "https://oauth.yandex.com/authorize?response_type=code"
)

type YandexDiskApiClient struct {
	Token   string
	Expires time.Time
}

func Authenticate(applicationId string) {
	url := YandexOAuthURL + "&client_id=" + applicationId
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
}
