package schedule

import (
	"log"
	"strings"
	"time"

	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/KevinStirling/scorebug.sh/data"
	"github.com/KevinStirling/scorebug.sh/internal/snapshots"
	"github.com/KevinStirling/scorebug.sh/ui/components/scorebug"
)

type Model struct {
	client     ScheduleClient
	games      []data.ScoreBug
	date       *time.Time
	paginator  paginator.Model
	tabs       []string
	tabContent []string
	activeTab  int
}

type tickMsg time.Time

func NewModel(client ScheduleClient) Model {
	now := time.Now()
	d := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	bugs := fetchScoreBugs(client, &d)

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	// might not be using the adaptive color right... might need to base it off the state rather than background
	p.ActiveDot = lipgloss.NewStyle().Foreground(adaptiveActive).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(adaptiveInactive).Render("•")
	p.SetTotalPages(len(bugs))
	return Model{
		client:    client,
		paginator: p,
		games:     bugs,
		date:      &d,
		tabs:      []string{"live", "preview", "final"},
		activeTab: 0,
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
		m.tabContent = renderSchedule(m.games, m.activeTab)
		if len(m.tabContent) == 0 {
			m.paginator.TotalPages = 1
			m.paginator.Page = 0
		} else {
			m.paginator.SetTotalPages(len(m.tabContent))
			if m.paginator.Page > m.paginator.TotalPages-1 {
				m.paginator.Page = m.paginator.TotalPages - 1
			}
		}
		return m, nil

	case tickMsg:
		return m, tea.Batch(m.checkServer(), tickAfter(10*time.Second))

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "l":
			m.activeTab = 0
		case "p":
			m.activeTab = 1
		case "f":
			m.activeTab = 2
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

func (m Model) View() tea.View {
	g := renderSchedule(m.games, m.activeTab)
	var b strings.Builder
	b.WriteString(primaryText.Render("\nscorebug.sh  ") + accentText.Render("l") + secondaryText.Render("ive • preview • final") + secondaryText.Render("\n"+strings.Repeat("‾", scorebug.SB_WIDTH)))
	start, end := m.paginator.GetSliceBounds(len(g))
	for _, item := range g[start:end] {
		b.WriteString("\n" + item)
	}
	b.WriteString("\n " + m.paginator.View())
	b.WriteString(secondaryText.Render("\n\n h/l ←/→ page • q: quit\n"))

	v := tea.NewView(divider.Render(b.String()))
	v.AltScreen = true

	return v
}

// Renders a string slice of scorebugs for a given Schedule type
func renderSchedule(bugs []data.ScoreBug, activeTab int) []string {
	out := make([]string, 0, len(bugs))
	for _, bug := range bugs {
		switch activeTab {
		case 0:
			if bug.Status == "Live" {
				out = append(out, scorebug.Render(bug))
			}
		case 1:
			if bug.Status == "Preview" {
				out = append(out, scorebug.Render(bug))
			}
		case 2:
			if bug.Status == "Final" || bug.Status == "Other" {
				out = append(out, scorebug.Render(bug))
			}
		}
	}
	return out
}

func tickAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) checkServer() tea.Cmd {
	return func() tea.Msg {
		return fetchScoreBugs(m.client, m.date)
	}
}

func fetchScoreBugs(client ScheduleClient, date *time.Time) []data.ScoreBug {
	sched, err := client.Schedule(date)
	if err != nil {
		log.Fatal("failed to fetch schedule", "error", err)
	}

	snaps, err := snapshots.Build(client, sched)
	if err != nil {
		log.Fatal("failed to build snapshots", "error", err)
	}

	return data.BuildScoreBugs(snaps)
}
