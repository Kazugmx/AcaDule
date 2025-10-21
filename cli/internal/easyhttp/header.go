package easyhttp

import "net/http"

func AuthBearer(req *http.Request, token string) {
	req.Header.Add("Authorization", "Bearer "+token)
}

func AcceptJson(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}
