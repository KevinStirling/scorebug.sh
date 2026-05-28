package schedule

import (
	"os"

	"charm.land/lipgloss/v2"
)

var (
	grey          = lipgloss.Lighten(lipgloss.BrightBlack, .1)
	divider       = lipgloss.NewStyle().Padding(0, 1)
	primaryText   = lipgloss.NewStyle().Foreground(lipgloss.Magenta)
	secondaryText = lipgloss.NewStyle().Foreground(grey)
	accentText    = lipgloss.NewStyle().Foreground(lipgloss.BrightYellow)

	hasDark          = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark        = lipgloss.LightDark(hasDark)
	adaptiveInactive = lightDark(lipgloss.Black, lipgloss.BrightBlack)
	adaptiveActive   = lightDark(grey, lipgloss.BrightWhite)
)
