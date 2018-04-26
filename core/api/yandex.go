package api

import (
	"backuper/core/common"
	"backuper/core/logging"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	YandexOAuthURL      = "https://oauth.yandex.com/authorize?response_type=code"
	YandexTokenFileName = "yandex_token"
)

type OAuthResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
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

	logging.Debug("Token readed", oauth.AccessToken)
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
	createBackuperFolderIfNotExist(client.Token)
	result := requestUploadURL(filename, client.Token)
	uploadFile(*result, filename, client.Token)
	logging.Debug(result)
	return BackupResult{Status: Success}
}

type uploadURLResponse struct {
	Href      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

var baseURL = "https://cloud-api.yandex.net/v1/disk/resources"
var folder = "backuper_app/"

func requestUploadURL(filename, token string) *string {
	client := &http.Client{}
	data := url.Values{}
	data.Set("path", folder+filename)
	data.Set("overwrite", "true")
	url := baseURL + "/upload?" + data.Encode()
	logging.Debug("Sending request to ", url)
	request, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		logging.Error(err)
		return nil
	}
	request.Header.Add("Authorization", "OAuth "+token)
	response, err := client.Do(request)
	if err != nil {
		logging.Error(err)
		return nil
	}
	body, _ := ioutil.ReadAll(response.Body)
	logging.Debug("Got a response", string(body))

	var responseStruct uploadURLResponse
	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		logging.Error(err)
		return nil
	}
	return &responseStruct.Href
}

func uploadFile(uploadingURL, filename, token string) {
	client := &http.Client{}
	file, err := os.Open(filename)
	if err != nil {
		logging.Error(err)
	}
	request, err := http.NewRequest("PUT", uploadingURL, file)
	if err != nil {
		logging.Error(err)
	}
	request.Header.Add("Authorization", "OAuth "+token)
	response, err := client.Do(request)
	if err != nil {
		logging.Error(err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	logging.Debug("Got a response", string(body))
}

func createBackuperFolderIfNotExist(token string) {
	client := &http.Client{}
	data := url.Values{}
	data.Set("path", folder)
	url := baseURL + "?" + data.Encode()
	logging.Debug("Sending request to ", url)
	request, err := http.NewRequest("PUT", url, strings.NewReader(""))
	if err != nil {
		logging.Error(err)
	}
	request.Header.Add("Authorization", "OAuth "+token)
	response, err := client.Do(request)
	if err != nil {
		logging.Error(err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	logging.Debug("Got a response", string(body))
}
