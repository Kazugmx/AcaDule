package simpleform

import "github.com/charmbracelet/huh"

func Confirm(title string) (confirm bool) {
	huh.NewConfirm().
		Title(title).
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm).
		Run()

	return
}
