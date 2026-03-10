package scorebug

import "charm.land/lipgloss/v2"

var (
	purple    = lipgloss.Color("#5F5FAF")
	cellStyle = lipgloss.NewStyle().Padding(0, 1)
	outsCol   = lipgloss.NewStyle().BorderLeft(false).Width(3).Align(lipgloss.Center)
)
