package api

import (
	"backuper/core/common"
	"encoding/json"
	"log"
	"os/exec"
	"time"
)

const (
	YandexOAuthURL      = "https://oauth.yandex.com/authorize?response_type=code"
	YandexTokenFileName = "yandex_token"
)

type OAuthResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json: "access_token"`
	ExpiresIn    int    `json: "expires_in"`
	RefreshToken string `json: "refresh_token"`
}

type YandexApiClient struct {
	Token   string
	Expires time.Time
}

func NewYandexApiClient(applicationId string) *YandexApiClient {
	var oauth OAuthResponse
	if !common.IsExist("C:\\Connector\\stuff\\go\\src\\backuper\\yandex_token") {
		authenticate(applicationId)
	}

	content := common.ReadFile(YandexTokenFileName)
	if error := json.Unmarshal(*content, &oauth); error != nil {
		log.Fatal(error)
	}
	return &YandexApiClient{
		Token: oauth.AccessToken,
		// this made for the testing purposes only
		// TODO:
		//       # Need to store correct token expiration date
		//       # Check if token is still fresh. Otherwise refresh it
		Expires: time.Now().Add(time.Second * time.Duration(oauth.ExpiresIn)),
	}
}

func authenticate(applicationId string) {
	url := YandexOAuthURL + "&client_id=" + applicationId
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	common.WaitUntil(func() bool {
		return common.IsExist(YandexTokenFileName)
	}, 5*time.Minute)
}

func (client *YandexApiClient) Backup(filename string) BackupResult {
	log.Print("Yandex client is backing up...")
	return BackupResult{Status: Success}
}
