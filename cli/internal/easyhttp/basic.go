package easyhttp

import (
	"net"
	"net/http"
	"time"
)

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
)

func Do(request *http.Request) (*http.Response, error) {
	return httpClient.Do(request)
}

func Get(url string) (*http.Response, error) {
	return httpClient.Get(url)
}
