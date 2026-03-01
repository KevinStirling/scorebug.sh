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
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	client := mlbstats.New()
	m := schedule.NewModel(client)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		log.Fatal("failed to start", "error", err)
	}
}
