package easyhttp

import (
	"bytes"
	"net/http"
)

func PostJsonBytes(url string, body []byte) (*http.Response, error) {
	return httpClient.Post(url, "application/json", bytes.NewBuffer(body))
}

func PostJsonWithBearer(url, token string, body []byte) (response *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	// add header
	AcceptJson(req)
	AuthBearer(req, token)

	response, err = Do(req)
	return
}

func GetJsonWithBearer(url, token string) (res *http.Response, err error) {
	// create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	// add headers
	AcceptJson(req)
	AuthBearer(req, token)

	// request
	res, err = Do(req)
	return
}
