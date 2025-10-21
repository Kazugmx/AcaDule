package acaduleapi

import (
	"acadule-cli/internal/easyhttp"
	"encoding/json"
	"io"
)

type RegisterRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Mail     string `json:"mail,omitempty"`
}

type RegisterResponse struct {
	Status any `json:"status,omitempty"`
	ID     int `json:"id,omitempty"`
}

func Register(apiUrl string, request RegisterRequest) (response *RegisterResponse, statusCode int, err error) {
	// marshal request
	marshaledData, err := json.Marshal(request)
	if err != nil {
		return
	}

	// post create user
	res, err := easyhttp.PostJsonBytes(apiUrl+"/auth/createUser", marshaledData)
	if err != nil {
		return
	}
	defer res.Body.Close()
	statusCode = res.StatusCode

	// read all body
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
