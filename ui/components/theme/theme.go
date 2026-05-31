package theme

import (
	"charm.land/lipgloss/v2"
)

// TODO build the main app theme here to use globally
var (
	grey          = lipgloss.Lighten(lipgloss.BrightBlack, .1)
	Divider       = lipgloss.NewStyle().Padding(0, 1)
	SecondaryText = lipgloss.NewStyle().Foreground(grey)
)
