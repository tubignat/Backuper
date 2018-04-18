package api

import (
	"os/exec"
)

const (
	YandexOAuthURL = "https://oauth.yandex.com/authorize?response_type=code"
)

type YandexDiskApiClient struct {
}

func Authenticate(applicationId string) {
	url := YandexOAuthURL + "&client_id=" + applicationId
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
}
