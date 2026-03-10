package schedule

import (
	"os"

	"charm.land/lipgloss/v2"
)

var (
	divider       = lipgloss.NewStyle().Padding(0, 1)
	primaryText   = lipgloss.NewStyle().Foreground(lipgloss.Green)
	secondaryText = lipgloss.NewStyle().Foreground(lipgloss.Black)

	hasDark       = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark     = lipgloss.LightDark(hasDark)
	adaptiveBlack = lightDark(lipgloss.BrightBlack, lipgloss.Black)
)
