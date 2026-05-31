package theme

import (
	"charm.land/lipgloss/v2"
)

// TODO build the main app theme here to use globally
var (
	grey          = lipgloss.Lighten(lipgloss.BrightBlack, .1)
	MainView      = lipgloss.NewStyle()
	SecondaryText = lipgloss.NewStyle().Foreground(grey)
	Margin        = 2
)
