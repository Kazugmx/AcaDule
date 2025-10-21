package acaduleapi

import (
	"acadule-cli/internal/easyhttp"
	"encoding/json"
	"io"
)

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Status any    `json:"status"`
	Token  string `json:"token"`
}

func Login(apiUrl string, request LoginRequest) (response *LoginResponse, statusCode int, err error) {
	marshaledData, err := json.Marshal(request)
	if err != nil {
		return
	}

	res, err := easyhttp.PostJsonBytes(apiUrl+"/auth/login", marshaledData)
	if err != nil {
		return
	}
	defer res.Body.Close()
	statusCode = res.StatusCode

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return
	}

	return
}
