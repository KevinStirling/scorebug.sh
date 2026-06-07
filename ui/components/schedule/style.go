package schedule

import (
	"os"

	"charm.land/lipgloss/v2"
	"github.com/KevinStirling/scorebug.sh/ui/components/scorebug"
	"github.com/KevinStirling/scorebug.sh/ui/components/theme"
)

var (
	grey          = lipgloss.Lighten(lipgloss.BrightBlack, .1)
	divider       = lipgloss.NewStyle().Width(scorebug.SB_WIDTH + theme.Margin)
	secondaryText = lipgloss.NewStyle().Foreground(grey)

	hasDark          = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark        = lipgloss.LightDark(hasDark)
	adaptiveInactive = lightDark(lipgloss.Black, lipgloss.BrightBlack)
	adaptiveActive   = lightDark(grey, lipgloss.BrightWhite)
)
