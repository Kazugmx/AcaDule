// Package cmd /*
package cmd

import (
	"acadule-cli/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	username string
	apiURL   string
	check    bool
)

type loginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type tokenTry struct {
	Status string `json:"status"`
	Expiry string `json:"expiry"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to AcaDule platform",
	Long:  `Login to AcaDule platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Println("Error occurred on loading config:", err)
			os.Exit(1)
		}

		validateAndUpdateConfig(&cfg)

		// --- check config ---
		if check {
			tryCred(cfg)
		}

		// --- login input ---

		fmt.Print("Password >")
		password, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Println("\nError occurred on reading password:", err)
			os.Exit(1)
		}
		fmt.Println("")

		loginCred := map[string]string{
			"username": username,
			"password": string(password),
		}
		jsonData, _ := json.Marshal(loginCred)

		// try login&get token
		res, err := http.Post(apiURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
		body, _ := io.ReadAll(res.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)

		if res.StatusCode != http.StatusOK {
			fmt.Println("WARN: Login failed. Status code:", res.StatusCode)
			fmt.Println("Response:", string(body))
			os.Exit(1)
		}

		var loginRes loginResponse
		if err := json.Unmarshal(body, &loginRes); err != nil {
			fmt.Println("Failed to parse response:", err)
			os.Exit(1)
		}
		cfg.Token = loginRes.Token

		_ = config.Save(cfg)
		fmt.Println("Login success. Config saved to", config.GetPath())
	},
}

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
		fmt.Print("Do you want to update config? [y/N] >")
		var input string
		_, _ = fmt.Scanln(&input)
		if input == "y" || input == "Y" {
			cfg.ApiURL = apiURL
			fmt.Println("Updated config:", cfg.ApiURL)
		} else {
			fmt.Println("Continue with current config.", apiURL)
		}
	}
	if cfg.Username != username && cfg.Username != "" {
		fmt.Println("WARN: Existing username is different from user input.\nconfig :", cfg.Username, "\nuser input:", username)
		fmt.Print("Do you want to update config? [y/N] >")
		var input string
		_, _ = fmt.Scanln(&input)
		if input == "y" || input == "Y" {
			cfg.Username = username
			fmt.Println("Updated config:", cfg.Username)
		} else {
			fmt.Println("Continue with current config.", username)
		}
	}
}

func tryCred(cfg config.Config) {
	client := &http.Client{}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Authorization", "Bearer "+cfg.Token)
	fmt.Println(`Current config >>
api_url: ` + cfg.ApiURL + `
username: ` + cfg.Username + `
token: ` + cfg.Token)

	fmt.Print("Do you want to try login? [y/N] >")
	var input string
	_, _ = fmt.Scanln(&input)
	fmt.Println("")
	if input == "y" || input == "Y" {

		req, _ := http.NewRequest("GET", apiURL+"/auth/me", nil)
		req.Header = header
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error occurred on login:", err)
			os.Exit(1)
		}
		result, _ := io.ReadAll(res.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		if res.StatusCode != http.StatusOK {
			fmt.Println("Login failed. Status code:", res.StatusCode)
			fmt.Println("Response:", string(result))
			os.Exit(1)
		}

		var tryResult tokenTry
		_ = json.Unmarshal(result, &tryResult)
		fmt.Println("Login result:", tryResult.Status, "\nExpiry:", tryResult.Expiry)
	}
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	loginCmd.Flags().StringVarP(&apiURL, "api-url", "a", "dev", "API URL")
	loginCmd.Flags().BoolVarP(&check, "check", "c", false, "Check config")
}
