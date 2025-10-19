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
	mailAddr string
)

type registerResponse struct {
	Status bool `json:"status"`
	ID     int  `json:"id"`
}

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register account to AcaDule platform",
	Long:  `Create account to AcaDule platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Println("Error occurred on loading config:", err)
			os.Exit(1)
		}

		validateAndUpdateConfig(&cfg)

		// --- register input ---
		fmt.Print("Password >")
		password, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Println("\nError occurred on reading password:", err)
			os.Exit(1)
		}
		fmt.Println("")
		fmt.Print("Confirm password >")
		confirmPassword, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Println("\nError occurred on reading password:", err)
			os.Exit(1)
		} else if string(password) != string(confirmPassword) {
			fmt.Println("\nPasswords do not match.")
			os.Exit(1)
		}

		httpClient := &http.Client{}
		registerCred := map[string]string{
			"username": username,
			"password": string(password),
			"mail":     mailAddr,
		}
		registerReq, _ := json.Marshal(registerCred)

		req, _ := http.NewRequest("POST", apiURL+"/auth/createUser", bytes.NewBuffer(registerReq))
		req.Header.Set("Content-Type", "application/json") // Added this line
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error occurred on register:", err)
			os.Exit(1)
		}
		fmt.Println("\nRegistering...")
		result, _ := io.ReadAll(res.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)

		var registerRes registerResponse
		if err := json.Unmarshal(result, &registerRes); err != nil {
		}

		if res.StatusCode != http.StatusOK {
			fmt.Println("Register failed. Status code:", res.StatusCode, "\nid:", registerRes.ID)
			if res.StatusCode == http.StatusInternalServerError {
				fmt.Println("Internal server error.")
			} else if registerRes.Status == false {
				fmt.Println("Username or mail address is already registered.")
			}
			fmt.Println("Please try again later.")

			os.Exit(1)
		}

		fmt.Println("Register result:", registerRes.Status)

		_ = config.Save(cfg)
		fmt.Println("Register success. Config saved to", config.GetConfigPath(), "\nLogin challenge in progress...")
		loginCred := map[string]string{
			"username": username,
			"password": string(password),
		}
		loginJsonData, _ := json.Marshal(loginCred)
		loginReq, err := http.NewRequest("POST", apiURL+"/auth/login", bytes.NewBuffer(loginJsonData))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRes, _ := httpClient.Do(loginReq)

		result, _ = io.ReadAll(loginRes.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)

		if loginRes.StatusCode != http.StatusOK {
			fmt.Println("Login failed. Status code:", loginRes.StatusCode)
			fmt.Println("Response:", string(result))
			os.Exit(1)
		}

		var loginResVal loginResponse
		if err := json.Unmarshal(result, &loginResVal); err != nil {
			fmt.Println("Failed to parse response:", err)
		}
		cfg.Token = loginResVal.Token
		_ = config.Save(cfg)
		fmt.Println("Login success. Config saved to", config.GetConfigPath())
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	registerCmd.Flags().StringVarP(&mailAddr, "mail-addr", "m", "", "Mail address")
	registerCmd.Flags().StringVarP(&apiURL, "api-url", "a", "dev", "API URL")

	_ = registerCmd.MarkFlagRequired("username")
	_ = registerCmd.MarkFlagRequired("mail-addr")
}
