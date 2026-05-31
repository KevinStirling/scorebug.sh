package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
	"github.com/KevinStirling/scorebug.sh/ui/components/schedule"
	"github.com/KevinStirling/scorebug.sh/ui/components/scorebug"
	"github.com/KevinStirling/scorebug.sh/ui/components/theme"
	"github.com/KevinStirling/scorebug.sh/ui/screen"
)

type Model struct {
	schedule schedule.Model
}

func NewModel() Model {
	c := mlbstats.New()
	return Model{
		schedule: schedule.NewModel(c),
	}
}

func (m Model) Init() tea.Cmd {
	return m.schedule.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "l":
			m.schedule.ActiveTab = 0
		case "s":
			m.schedule.ActiveTab = 1
		case "f":
			m.schedule.ActiveTab = 2

		}
	case tea.WindowSizeMsg:
		log.Info("window width", msg.Width)

		if pages, err := screen.GetSchedulePageSize(scorebug.SB_HEIGHT, scorebug.SB_MARGIN, msg.Height); err != nil {
			log.Fatal("failed to determine scheudle page size for provided SB_HEIGHT", "SB_HEIGHT", scorebug.SB_HEIGHT)
		} else {
			m.schedule.Paginator.PerPage = pages
		}
	}

	// forward ui inputs to the child components
	var cmd tea.Cmd
	m.schedule, cmd = m.schedule.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	help := theme.SecondaryText.Render("\n\n n/p ←/→ page • q: quit • l: live • s: scheduled • f: final\n")
	content := lipgloss.JoinVertical(lipgloss.Top, m.schedule.View(), help)

	v := tea.NewView(theme.Divider.Render(content))
	v.AltScreen = true
	return v
}
