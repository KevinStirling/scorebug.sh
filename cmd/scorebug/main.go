package main

import (
	"fmt"
	"os"

	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
	"github.com/KevinStirling/scorebug.sh/ui/components/schedule"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	client := mlbstats.New()
	m := schedule.NewModel(client)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("oy, ya cooked, mate - %s", err.Error())
		os.Exit(1)
	}
}
