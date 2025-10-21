// Package cmd /*
package cmd

import (
	"acadule-cli/internal/acaduleapi"
	"acadule-cli/internal/config"
	"acadule-cli/internal/simpleform"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	username string
	apiURL   string
	check    bool
)

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
		password := simpleform.AskPassword(false)
		if password == nil {
			fmt.Println("Failed to get password from input")
			os.Exit(1)
		}

		loginReq := acaduleapi.LoginRequest{
			Username: username,
			Password: *password,
		}
		loginRes, statusCode, err := acaduleapi.Login(apiURL, loginReq)

		if statusCode != http.StatusOK {
			fmt.Println("WARN: Login failed. Status code:", statusCode)
			fmt.Println("Response:", *loginRes)
			os.Exit(1)
		}

		cfg.Token = loginRes.Token

		_ = config.Save(cfg)
		fmt.Println("Login success. Config saved to", config.GetConfigPath())
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	loginCmd.Flags().StringVarP(&apiURL, "api-url", "a", "dev", "API URL")
	loginCmd.Flags().BoolVarP(&check, "check", "c", false, "Check config")
}
