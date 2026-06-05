package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
	"github.com/KevinStirling/scorebug.sh/ui/components/game"
	"github.com/KevinStirling/scorebug.sh/ui/components/schedulelist"
	"github.com/KevinStirling/scorebug.sh/ui/components/scorebug"
	"github.com/KevinStirling/scorebug.sh/ui/components/theme"
	"github.com/KevinStirling/scorebug.sh/ui/screen"
)

type Model struct {
	schedule        schedulelist.Model
	game            game.Model
	containerHeight int
}

func NewModel() Model {
	c := mlbstats.New()
	return Model{
		schedule: schedulelist.NewModel(c),
		game:     game.NewModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.schedule.Init(), m.game.Init())
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
		if pages, err := screen.GetSchedulePageSize(scorebug.SB_HEIGHT, scorebug.SB_MARGIN, msg.Height); err != nil {
			log.Fatal("failed to determine scheudle page size for provided SB_HEIGHT", "SB_HEIGHT", scorebug.SB_HEIGHT)
		} else {
			m.containerHeight = msg.Height
			m.game.ContainerHeight = scorebug.SB_HEIGHT*pages - theme.Margin
			m.game.ContainerWidth = msg.Width - scorebug.SB_WIDTH - game.OffsetVerticalMargin
		}
	}

	var scheduleUpdate, gameUpdate tea.Cmd
	m.schedule, scheduleUpdate = m.schedule.Update(msg)
	m.game, gameUpdate = m.game.Update(msg)
	return m, tea.Batch(scheduleUpdate, gameUpdate)
}

func (m Model) View() tea.View {
	help := theme.SecondaryText.AlignVertical(lipgloss.Bottom).Render("\n\n n/p ←/→ page • q: quit • l: live • s: scheduled • f: final\n")
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, m.schedule.View(), m.game.View())
	content := lipgloss.JoinVertical(lipgloss.Left, mainContent, help)

	v := tea.NewView(theme.MainView.Render(content))
	v.AltScreen = true
	return v
}
