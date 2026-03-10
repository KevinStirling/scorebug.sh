package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
	"github.com/KevinStirling/scorebug.sh/ui/components/schedule"
	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close log file: %v", err)
		}
	}()

	client := mlbstats.New()
	m := schedule.NewModel(client)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal("failed to start", "error", err)
	}
}
