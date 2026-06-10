package game

import "charm.land/lipgloss/v2"

var (
	headerStyle    = lipgloss.NewStyle().Align(lipgloss.Center)
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Margin(1, 1, 2, 1)
)
