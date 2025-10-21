package acaduleapi

import (
	"acadule-cli/internal/config"
	"acadule-cli/internal/easyhttp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type MeResponse struct {
	Status string `json:"status"`
	Expiry string `json:"expiry"`
}

func GetMe(apiUrl string, config config.Config) (response MeResponse, statusCode int, err error) {
	req, err := http.NewRequest(http.MethodGet, apiUrl+"/auth/me", nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)
	res, err := easyhttp.Do(req)
	if err != nil {
		fmt.Println("Error occurred on login:", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	statusCode = res.StatusCode

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(result, &response)
	return
}
