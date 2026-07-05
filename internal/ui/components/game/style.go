package game

import "charm.land/lipgloss/v2"

var (
	headerStyle    = lipgloss.NewStyle().Align(lipgloss.Center)
	cellStyle      = lipgloss.NewStyle().Padding(0, 1)
	linescoreStyle = lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Center)
	containerStyle = lipgloss.NewStyle().
			Margin(1, 1, 2, 1)

	playFeedEven = lipgloss.NewStyle()
	playFeedOdd  = lipgloss.NewStyle().Background(lipgloss.BrightBlack)
)
