package cmd

import (
	"acadule-cli/internal/acaduleapi"
	"acadule-cli/internal/config"
	"acadule-cli/internal/simpleform"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	mailAddr string
)

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

		password := simpleform.AskPassword(true)
		if password == nil {
			fmt.Println("Password doesn't match.")
			os.Exit(1)
		}

		// create request data
		registerReq := acaduleapi.RegisterRequest{
			Username: username,
			Password: *password,
			Mail:     mailAddr,
		}

		// do register request
		res, statusCode, err := acaduleapi.Register(
			apiURL,
			registerReq,
		)
		if err != nil {
			slog.Error("Failed to request api", slog.Any("error", err))
			os.Exit(1)
		}

		// check failure and handle it
		if fmt.Sprintf("%v", res.Status) != "true" {
			fmt.Println("Register failed. Status code:", statusCode, "\nid:", res.ID)
			// on Fail
			if statusCode == http.StatusInternalServerError {
				fmt.Println("Internal server error occurred.")
			} else {
				fmt.Println("Status:", res.Status)
				fmt.Println("Username or mail address is already registered.")
			}
			os.Exit(1)
		}

		// on Success
		fmt.Println("Register result:", res.Status)

		_ = config.Save(cfg)
		fmt.Println("Register success. Config saved to", config.GetConfigPath())

		// login challenge
		fmt.Println("Login challenge in progress...")
		loginReq := acaduleapi.LoginRequest{
			Username: username,
			Password: *password,
		}
		loginRes, statusCode, err := acaduleapi.Login(apiURL, loginReq)

		if statusCode != http.StatusOK {
			fmt.Println("Login failed. Status code:", statusCode)
			fmt.Println("Response:", *loginRes)
			os.Exit(1)
		}

		cfg.Token = loginRes.Token
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
