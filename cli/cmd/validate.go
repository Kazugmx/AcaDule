package cmd

import (
	"acadule-cli/internal/acaduleapi"
	"acadule-cli/internal/config"
	"acadule-cli/internal/simpleform"
	"fmt"
	"net/http"
	"os"
)

func validateAndUpdateConfig(cfg *config.Config) {

	//  --- config input validation ---
	if apiURL == "dev" {
		apiURL = "http://localhost:8080"
		fmt.Println("API URL is set to development setting: ", apiURL)
	} else if apiURL != "" {
		cfg.ApiURL = apiURL
		fmt.Println("API URL is set to user input: ", apiURL)
	}

	if cfg.Username != "" && username == "" {
		username = cfg.Username
		fmt.Println("username is set from config: ", username)
	} else if cfg.Username == "" && username != "" {
		fmt.Println("username is set from user input: ", username)
		cfg.Username = username
	} else if cfg.Username == "" && username == "" {
		fmt.Println("Error: --username flag is required")
		os.Exit(1)
	}

	// --- config update selection ----
	if cfg.ApiURL != apiURL && cfg.ApiURL != "" {
		fmt.Println("WARN: Existing API URL is different from user input.\nconfig :", cfg.ApiURL, "\nuser input:", apiURL)
		if simpleform.Confirm("Do you want to update config?") {
			cfg.ApiURL = apiURL
			fmt.Println("Updated config:", cfg.ApiURL)
		} else {
			fmt.Println("Continue with current config.", apiURL)
		}
	}
	if cfg.Username != username && cfg.Username != "" {
		fmt.Println("WARN: Existing username is different from user input.\nconfig :", cfg.Username, "\nuser input:", username)
		if simpleform.Confirm("Do you want to update config?") {
			cfg.Username = username
			fmt.Println("Updated config:", cfg.Username)
		} else {
			fmt.Println("Continue with current config.", username)
		}
	}
}

func tryCred(cfg config.Config) {
	fmt.Println("Current config >>")
	fmt.Println("api_url: " + cfg.ApiURL)
	fmt.Println("username: " + cfg.Username)
	fmt.Println("token: " + cfg.Token)

	if simpleform.Confirm("Do you want to try login?") {

		res, statusCode, err := acaduleapi.GetMe(apiURL, cfg)
		if err != nil {
			fmt.Println("Error occurred on login:", err)
			os.Exit(1)
		}

		if statusCode != http.StatusOK {
			fmt.Println("Login failed. Status code:", statusCode)
			fmt.Println("Response:", res)
			os.Exit(1)
		}

		fmt.Println("Login result:", res.Status, "\nExpiry:", res.Expiry)
	}
	os.Exit(0)
}
