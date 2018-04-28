package api

import (
	"backuper/core/common"
	"backuper/core/logging"
	"backuper/core/net"
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// YandexOAuthURL represents URL of Yandex OAuth service
	YandexOAuthURL = "https://oauth.yandex.com/authorize?response_type=code"

	// YandexTokenFileName is a name for the file that should contain authentication info.
	// New instance of YandexAPIClient will look for this file and read a token if one exists,
	// otherwise client will start authentication process and wait until the file appears
	YandexTokenFileName = "yandex_token"

	baseURL = "https://cloud-api.yandex.net/v1/disk/resources"
	folder  = "backuper_app/"
)

// OAuthResponse is a struct that represents Yandex.API OAuth response
type OAuthResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// YandexAPIClient is an implementation of api.Client interface for Yandex.Disk API
type YandexAPIClient struct {
	httpClient *net.ClientWrapper
	Token      string
	Expires    time.Time
}

// NewYandexAPIClient creates a new instance of YandexApiClient.
// Calls authentication if user is not authenticated yet
func NewYandexAPIClient(applicationID string) *YandexAPIClient {
	var oauth OAuthResponse
	if !common.IsExist("yandex_token") {
		authenticate(applicationID)
	}

	content := common.ReadFile(YandexTokenFileName)
	if error := json.Unmarshal(*content, &oauth); error != nil {
		logging.Error("Could not read the token from a file ", error)
		return nil
	}

	logging.Debug("Token readed", oauth.AccessToken)
	return &YandexAPIClient{
		httpClient: net.NewHttpClientWrapper(baseURL, getAuthHeader(oauth.AccessToken)),
		Token:      oauth.AccessToken,
		// this made for the testing purposes only
		// TODO:
		//       # Need to store correct token expiration date
		//       # Check if token is still fresh. Otherwise refresh it
		Expires: time.Now().Add(time.Second * time.Duration(oauth.ExpiresIn)),
	}
}

// Backup stores a file on Yandex.Disk. If file is already exist there, it will be overwritten
func (client *YandexAPIClient) Backup(filename string) BackupResult {
	client.createBackuperFolderIfNotExist()
	result := client.requestUploadURL(filename)
	uploadFile(*result, filename, client.Token)
	logging.Debug(result)
	return BackupResult{Status: Success}
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

type uploadURLResponse struct {
	Href      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

func (yandexAPIClient *YandexAPIClient) requestUploadURL(filename string) *string {
	data := url.Values{}
	data.Set("path", folder+filename)
	data.Set("overwrite", "true")

	var response uploadURLResponse
	yandexAPIClient.httpClient.GET("/upload", data, &response)

	return &response.Href
}

func uploadFile(uploadingURL, filename, token string) {
	client := net.NewHttpClientWrapper(uploadingURL, getAuthHeader(token))
	file, err := os.Open(filename)
	if err != nil {
		logging.Error(err)
	}

	client.PUT("", url.Values{}, file)
}

func (yandexAPIClient *YandexAPIClient) createBackuperFolderIfNotExist() {
	data := url.Values{}
	data.Set("path", folder)

	yandexAPIClient.httpClient.PUT("", data, strings.NewReader(""))
}

func getAuthHeader(token string) net.Header {
	return net.Header{Key: "Authorization", Value: "OAuth " + token}
}
