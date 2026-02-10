package schedule

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/KevinStirling/scorebug.sh/data"
	"github.com/KevinStirling/scorebug.sh/ui/components/scorebug"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	SB_WIDTH  = 58
	SB_HEIGHT = 5
)

var (
	divider       = lipgloss.NewStyle().Padding(0, 1)
	primaryText   = lipgloss.NewStyle().Foreground(lipgloss.Color("253"))
	secondaryText = lipgloss.NewStyle().Foreground(lipgloss.Color("247"))
)

type Model struct {
	client    StatsClient
	games     data.Schedule
	date      *time.Time
	paginator paginator.Model
	err       error
}

type tickMsg time.Time

func NewModel(client StatsClient) Model {
	d := time.Date(2025, time.September, 28, 0, 0, 0, 0, time.Local)
	res, err := client.Schedule(&d)
	if err != nil {
		log.Fatal("failed to fetch schedule", "error", err)
	}
	games := data.BuildSchedule(res)
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(pages(len(games.Games), p.PerPage))
	return Model{
		client:    client,
		paginator: p,
		games:     games,
		date:      &d,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.checkServer(), tickAfter(10*time.Second))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case data.Schedule:
		m.games = msg
		return m, nil

	case tickMsg:
		return m, tea.Batch(m.checkServer(), tickAfter(10*time.Second))

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	if m.paginator.Page >= m.paginator.TotalPages-1 {
		m.paginator.Page = m.paginator.TotalPages - 1
	}
	if m.paginator.Page < 0 {
		m.paginator.Page = 0
	}
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	g := renderSchedule(m.games)
	var b strings.Builder
	b.WriteString(primaryText.Render("\n scorebug.sh ") + secondaryText.Render("\n"+strings.Repeat("‾", SB_WIDTH)))
	start, end := m.paginator.GetSliceBounds(len(g))
	for _, item := range g[start:end] {
		b.WriteString("\n" + item)
	}
	b.WriteString("\n " + m.paginator.View())
	b.WriteString(secondaryText.Render("\n\n h/l ←/→ page • q: quit\n"))
	return divider.Render(b.String())
}

// Renders a string slice of scorebugs for a given Schedule type
func renderSchedule(g data.Schedule) []string {
	var bugCells []string
	if len(g.Games) > 0 {
		for _, game := range g.Games {
			bugStr := scorebug.Render(game)
			bugCells = append(bugCells, bugStr)
		}
	}
	return bugCells
}

func tickAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) checkServer() tea.Cmd {
	return func() tea.Msg {
		res, err := m.client.Schedule(m.date)
		if err != nil {
			log.Fatal("failed to referesh schedule", "error", err)
		}
		return data.BuildSchedule(res)
	}
}

func pages(items, perPage int) int {
	if perPage <= 0 {
		return 1
	}
	if items == 0 {
		return 1
	}
	return (items + perPage - 1) / perPage
}
