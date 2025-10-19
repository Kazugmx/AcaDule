package simpleform

import "github.com/charmbracelet/huh"

func Ask(question string) (result string) {
	huh.NewInput().
		Title(question).
		Value(&result).
		Run()
	return
}
