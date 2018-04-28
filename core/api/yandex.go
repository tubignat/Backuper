package api

import (
	"backuper/core/common"
	"backuper/core/logging"
	"backuper/core/net"
	"backuper/core/oauth"
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
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

// YandexAPIClient is an implementation of api.Client interface for Yandex.Disk API
type YandexAPIClient struct {
	httpClient *net.ClientWrapper
	Token      string
	Expires    time.Time
}

// NewYandexAPIClient creates a new instance of YandexApiClient.
// Calls authentication if user is not authenticated yet
func NewYandexAPIClient(applicationID string) *YandexAPIClient {
	if !common.IsExist("yandex_token") {
		authenticate(applicationID)
	}

	var token oauth.Token
	content := common.ReadFile(YandexTokenFileName)
	if error := json.Unmarshal(*content, &token); error != nil {
		logging.Error("Could not read the token from a file ", error)
		return nil
	}

	if time.Now().After(token.Expires) {
		authenticate(applicationID)
	}

	logging.Debug("Token readed", token.Value)
	return &YandexAPIClient{
		httpClient: net.NewHttpClientWrapper(baseURL, getAuthHeader(token.Value)),
		Token:      token.Value,
		Expires:    token.Expires,
	}
}

// Backup stores a file on Yandex.Disk. If file is already exist there, it will be overwritten
func (client *YandexAPIClient) Backup(filename string) BackupResult {
	name := getFileNameWithoutVolumeName(filename)
	result := client.requestUploadURL(name)
	uploadFile(*result, name, client.Token)
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

func (client *YandexAPIClient) requestUploadURL(filename string) *string {
	data := url.Values{}
	data.Set("path", folder+filename)
	data.Set("overwrite", "true")

	var response uploadURLResponse
	client.httpClient.GET("/upload", data, &response)

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

func getAuthHeader(token string) net.Header {
	return net.Header{Key: "Authorization", Value: "OAuth " + token}
}

func getFileNameWithoutVolumeName(filename string) string {
	absolutePath, err := filepath.Abs(filename)
	if err != nil {
		logging.Error(err)
	}
	volume := filepath.VolumeName(absolutePath)
	return strings.TrimPrefix(absolutePath, volume)
}
