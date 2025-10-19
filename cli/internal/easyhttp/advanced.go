package easyhttp

import (
	"bytes"
	"net/http"
)

func PostJsonBytes(url string, body []byte) (*http.Response, error) {
	return httpClient.Post(url, "application/json", bytes.NewBuffer(body))
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
