package api

import (
	"backuper/core/common"
	"backuper/core/logging"
	"encoding/json"
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
	if !common.IsExist("yandex_token") {
		authenticate(applicationId)
	}

	content := common.ReadFile(YandexTokenFileName)
	if error := json.Unmarshal(*content, &oauth); error != nil {
		logging.Error("Could not read the token from a file ", error)
		return nil
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

func authenticate(applicationID string) {
	logging.Debug("Receiving Yandex token...")
	url := YandexOAuthURL + "&client_id=" + applicationID
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	common.WaitUntil(func() bool {
		return common.IsExist(YandexTokenFileName)
	}, 5*time.Minute)
	logging.Debug("Yandex token received...")
}

func (client *YandexApiClient) Backup(filename string) BackupResult {
	logging.Debug("Yandex client is backing up...")
	return BackupResult{Status: Success}
}
