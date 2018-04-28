package main

/*
 *	 oauth.exe is used for handling custom URL protocol backuper://
 *   APIs that use OAuth redirect to this protocol so the app can
 *   handle a request and get a verification code
 *
 *   oauth.exe must be located in the same directory as backuper.exe
 */

import (
	"backuper/core/api"
	"backuper/core/common"
	"backuper/core/logging"
	"backuper/core/oauth"
	"backuper/core/settings"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	directory := os.Args[1]
	callback := os.Args[2]
	logging.Configure(directory + "\\ouath_logs.txt")
	keeper := settings.NewKeeper(directory+"\\config.json", make(chan bool))
	HandleOAuthRequest(callback, directory, keeper.GetRelevantSettings())
}

func HandleOAuthRequest(callback, directory string, settings *settings.Settings) {
	logging.Info("Got a request: ", callback)
	request := strings.TrimPrefix(callback, "backuper://")
	if strings.Contains(request, api.YandexTokenFileName) {
		handleYandexAuth(request, directory, &settings.Yandex)
	}
}

func handleYandexAuth(request, directory string, settings *settings.YandexDiskSettings) {
	logging.Debug("Handling yandex auth request")
	log.Print("Handling yandex auth request")
	code := strings.TrimPrefix(request, api.YandexTokenFileName+"/?code=")
	data := formEncodedURLValues(code, settings.ApplicationID, settings.Password)
	response, err := http.Post("https://oauth.yandex.com/token", "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		logging.Error("Request failed... ", err)
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	logging.Debug("Got the response", string(body))
	if response.StatusCode != 200 {
		logging.Error("Something went wrong. Response is ", string(body))
		return
	}
	path := directory + "\\" + api.YandexTokenFileName
	token := getTokenBody(&body)

	common.WriteFile(path, token)
	logging.Debug("Yandex auth handling succeeded. File ", path, " is written")
}

func getTokenBody(responseBody *[]byte) *[]byte {
	var responseStruct oauth.APIResponse
	if err := json.Unmarshal(*responseBody, &responseStruct); err != nil {
		logging.Error(err)
	}
	token := oauth.Token{
		Value:   responseStruct.AccessToken,
		Expires: time.Now().Add(time.Second * time.Duration(responseStruct.ExpiresIn)),
	}
	result, err := json.Marshal(&token)
	if err != nil {
		logging.Error(err)
	}
	return &result
}

func formEncodedURLValues(code, applicationId, password string) string {
	logging.Debug("Start marshaling URL values for the request", code, applicationId, password)
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", applicationId)
	data.Set("client_secret", password)
	return data.Encode()
}
