package main

import (
	"fmt"
	"os"

	"github.com/KevinStirling/scorebug.sh/ui/components/scorebug"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(scorebug.NewModel()).Run(); err != nil {
		fmt.Printf("oy, ya cooked, mate - %s", err.Error())
		os.Exit(1)
	}
}
