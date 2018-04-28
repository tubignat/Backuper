package net

import (
	"backuper/core/logging"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ClientWrapper struct {
	client  *http.Client
	baseURL string
	headers []Header
}

type Header struct {
	Key   string
	Value string
}

func NewHttpClientWrapper(baseURL string, headers ...Header) *ClientWrapper {
	return &ClientWrapper{
		client:  &http.Client{},
		headers: headers,
		baseURL: baseURL,
	}
}

func (client *ClientWrapper) GET(url string, params url.Values, responseStruct interface{}) {
	fullURL := client.baseURL + url + "?" + params.Encode()

	logging.Debug("Sending request to ", fullURL)

	request, err := http.NewRequest("GET", fullURL, strings.NewReader(""))
	if err != nil {
		logging.Error(err)
		return
	}
	client.setHeaders(request)

	response, err := client.client.Do(request)
	if err != nil {
		logging.Error(err)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	logging.Debug("Got a response", string(body))

	err = json.Unmarshal(body, responseStruct)
	if err != nil {
		logging.Error(err)
		return
	}
}

func (client *ClientWrapper) PUT(url string, params url.Values, reader io.Reader) {
	fullURL := client.baseURL + "?" + params.Encode()
	logging.Debug("Sending request to ", fullURL)
	request, err := http.NewRequest("PUT", fullURL, reader)
	if err != nil {
		logging.Error(err)
	}
	client.setHeaders(request)

	response, err := client.client.Do(request)
	if err != nil {
		logging.Error(err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	logging.Debug("Got a response", string(body))
}

func (client *ClientWrapper) setHeaders(request *http.Request) {
	for _, header := range client.headers {
		request.Header.Add(header.Key, header.Value)
	}
}
