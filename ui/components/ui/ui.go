package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
	"github.com/KevinStirling/scorebug.sh/ui/components/game"
	"github.com/KevinStirling/scorebug.sh/ui/components/header"
	"github.com/KevinStirling/scorebug.sh/ui/components/schedule"
	"github.com/KevinStirling/scorebug.sh/ui/components/theme"
)

type Model struct {
	width, height int

	header   header.Model
	schedule schedule.Model
	game     game.Model
	help     help.Model
}

func NewModel() Model {
	c := mlbstats.New()
	return Model{
		schedule: schedule.New(c),
		game:     game.New(),
		help:     help.New(),
		header:   header.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.schedule.Init(), m.game.Init())
}

func (m Model) layout() Model {
	m.help.SetWidth(m.width)
	// add header here, and swap schedule.Keys with context aware keys
	chrome := lipgloss.Height(theme.MainView.Render(m.header.Render())) + lipgloss.Height(m.help.View(m.schedule.Keys))

	bodyH := m.height - chrome
	leftW := m.width / 2
	rightW := m.width - leftW

	m.schedule.SetSize(leftW, bodyH)
	m.game.SetSize(rightW, bodyH)
	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.schedule.IsFiltering() {
			break
		}
		switch {
		case key.Matches(msg, m.schedule.Keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m.layout(), nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m.layout(), nil
	case schedule.TabChangedMsg:
		m.header.ActiveTab = int(msg)
	}

	var scheduleUpdate, gameUpdate tea.Cmd
	m.schedule, scheduleUpdate = m.schedule.Update(msg)
	m.game, gameUpdate = m.game.Update(msg)
	return m, tea.Batch(scheduleUpdate, gameUpdate)
}

func (m Model) View() tea.View {
	header := lipgloss.JoinHorizontal(lipgloss.Top, m.header.Render())
	body := lipgloss.JoinHorizontal(lipgloss.Top, m.schedule.View(), m.game.View())
	main := lipgloss.JoinVertical(lipgloss.Left, header, body)
	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left,
		main, m.help.View(m.schedule.Keys)))
	v.AltScreen = true
	return v
}
