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
	"backuper/core/settings"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	directory := os.Args[1]
	callback := os.Args[2]
	configureLogging(directory + "\\logs.txt")
	keeper := settings.NewKeeper(directory + "\\config.json")
	HandleOAuthRequest(callback, directory, keeper.GetRelevantSettings())
}

func HandleOAuthRequest(callback, directory string, settings *settings.Settings) {
	request := strings.TrimPrefix(callback, "backuper://")
	if strings.Contains(request, api.YandexTokenFileName) {
		handleYandexAuth(request, directory, &settings.Yandex)
	}
}

func handleYandexAuth(request, directory string, settings *settings.YandexDiskSettings) {
	log.Print("Handling yandex oAuth callback")
	code := strings.TrimPrefix(request, api.YandexTokenFileName+"/?code=")
	data := formEncodedURLValues(code, settings.ApplicationID, settings.Password)
	response, err := http.Post("https://oauth.yandex.com/token", "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		log.Panic(err)
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		log.Panic("Something went wrong. Response is ", string(body))
	}
	common.WriteFile(directory+"\\"+api.YandexTokenFileName, &body)
	log.Print("Success. File with token created")
}

func formEncodedURLValues(code, applicationId, password string) string {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", applicationId)
	data.Set("client_secret", password)
	return data.Encode()
}

func configureLogging(path string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
}
