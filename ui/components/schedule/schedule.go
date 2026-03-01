package schedule

import (
	"log"
	"strings"
	"time"

	"github.com/KevinStirling/scorebug.sh/data"
	"github.com/KevinStirling/scorebug.sh/internal/snapshots"
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
	client    ScheduleClient
	games     []data.ScoreBug
	date      *time.Time
	paginator paginator.Model
	err       error
}

type tickMsg time.Time

func NewModel(client ScheduleClient) Model {
	now := time.Now()
	d := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	sched, err := client.Schedule(&d)
	if err != nil {
		log.Fatal("failed to fetch schedule", "error", err)
	}

	snaps, err := snapshots.Build(client, sched)
	if err != nil {
		log.Fatal("failed to build snapshots", "error", err)
	}

	bugs := data.BuildScoreBugs(snaps)

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(pages(len(bugs), p.PerPage))
	return Model{
		client:    client,
		paginator: p,
		games:     bugs,
		date:      &d,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.checkServer(), tickAfter(10*time.Second))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case []data.ScoreBug:
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
func renderSchedule(bugs []data.ScoreBug) []string {
	out := make([]string, 0, len(bugs))
	for _, bug := range bugs {
		out = append(out, scorebug.Render(bug))
	}
	return out
}

func tickAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) checkServer() tea.Cmd {
	return func() tea.Msg {
		sched, err := m.client.Schedule(m.date)
		if err != nil {
			log.Fatal("failed to refresh schedule", "error", err)
		}

		snaps, err := snapshots.Build(m.client, sched)
		if err != nil {
			log.Fatal("failed to build snapshots", "error", err)
		}

		return data.BuildScoreBugs(snaps)
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
