package playfeed

import (
	"charm.land/lipgloss/v2"
)

var (
	evenRow  = lipgloss.NewStyle()
	oddRow   = lipgloss.NewStyle().Background(lipgloss.Darken(lipgloss.BrightBlack, .7))
	titleBar = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Green).Foreground(lipgloss.Black)

	eventInning = lipgloss.NewStyle().Foreground(lipgloss.Yellow)
)
