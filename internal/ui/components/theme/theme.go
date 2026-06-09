package theme

import (
	"charm.land/lipgloss/v2"
)

var (
	grey = lipgloss.Lighten(lipgloss.BrightBlack, .1)

	MainView      = lipgloss.NewStyle()
	PrimaryText   = lipgloss.NewStyle().Foreground(lipgloss.Magenta)
	SecondaryText = lipgloss.NewStyle().Foreground(grey)
	AccentText    = lipgloss.NewStyle().Foreground(lipgloss.BrightYellow)
)
