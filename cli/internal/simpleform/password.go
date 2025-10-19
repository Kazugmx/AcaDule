package simpleform

import (
	"strings"

	"github.com/charmbracelet/huh"
)

func AskPassword(requireConfirm bool) *string {
	var password string
	huh.NewInput().
		Title("Enter password").
		Value(&password).
		EchoMode(huh.EchoModePassword).
		Run()

	if requireConfirm {
		var confirmPassword string
		huh.NewInput().
			Title("Confirm password").
			Value(&confirmPassword).
			EchoMode(huh.EchoModePassword).
			Run()
		if strings.EqualFold(password, confirmPassword) {
			return &password
		} else {
			return nil
		}
	} else {
		return &password
	}
}
