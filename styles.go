package main

import (
	"sync"

	"github.com/charmbracelet/lipgloss"
)

var ansiColors = []int{1, 9, 20, 22, 34, 52, 88, 100, 124, 130, 160, 200, 202, 240}

var usernameStyles sync.Map

var colorPickCounter int

func getUsernameStyle(username string) lipgloss.Style {
	value, ok := usernameStyles.Load(username)
	if ok {
		return value.(lipgloss.Style)
	} else {
		style := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.ANSIColor(ansiColors[colorPickCounter]))

		colorPickCounter++
		if colorPickCounter >= len(ansiColors) {
			colorPickCounter = 0
		}

		usernameStyles.Store(username, style)

		return style
	}
}
