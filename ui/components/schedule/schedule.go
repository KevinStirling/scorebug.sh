package schedule

import (
	"log"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
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
	Paginator  paginator.Model
	tabs       []string
	tabContent [3][]string
	ActiveTab  int
}

type tickMsg time.Time

func NewModel(client ScheduleClient) Model {
	now := time.Now()
	d := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	bugs := fetchScoreBugs(client, &d)

	p := paginator.New()
	p.KeyMap = paginator.KeyMap{
		NextPage: key.NewBinding(key.WithKeys("pgright", "n")),
		PrevPage: key.NewBinding(key.WithKeys("pgleft", "p")),
	}
	p.Type = paginator.Dots

	// TODO fix adaptiveActive color... don't think i'm using it right
	p.ActiveDot = lipgloss.NewStyle().Foreground(adaptiveActive).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(adaptiveInactive).Render("•")
	p.SetTotalPages(len(bugs))
	return Model{
		client:    client,
		Paginator: p,
		games:     bugs,
		date:      &d,
		tabs:      []string{"live", "scheduled", "final"},
		ActiveTab: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.checkServer(), tickAfter(10*time.Second))
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case []data.ScoreBug:
		m.games = msg
		for i := range m.tabs {
			m.tabContent[i] = renderTab(m.games, i)
		}
		m.syncPaginator()
		return m, nil

	case tickMsg:
		m.syncPaginator()
		return m, tea.Batch(m.checkServer(), tickAfter(10*time.Second))

	}
	m.syncPaginator()
	m.Paginator, cmd = m.Paginator.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	g := renderTab(m.games, m.ActiveTab)
	var b strings.Builder
	b.WriteString(renderHeader(m.tabs, m.ActiveTab))
	start, end := m.Paginator.GetSliceBounds(len(g))
	for _, item := range g[start:end] {
		for range scorebug.SB_MARGIN {
			b.WriteString("\n")
		}
		b.WriteString(item)
	}
	b.WriteString("\n")
	b.WriteString(m.Paginator.View())

	return divider.Render(b.String())
}

// Renders a string slice of scorebugs for a given tab
func renderTab(bugs []data.ScoreBug, tab int) []string {
	out := make([]string, 0, len(bugs))
	for _, bug := range bugs {
		switch tab {
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

func (m *Model) syncPaginator() {
	content := m.tabContent[m.ActiveTab]
	if len(content) == 0 {
		m.Paginator.TotalPages = 1
		m.Paginator.Page = 0
		return
	}
	m.Paginator.SetTotalPages(len(content))
	if m.Paginator.Page > m.Paginator.TotalPages-1 {
		m.Paginator.Page = m.Paginator.TotalPages - 1
	}
}

func renderHeader(tabs []string, activeTab int) string {
	parts := make([]string, len(tabs))
	for i, t := range tabs {
		if i == activeTab {
			parts[i] = accentText.Render(t)
		} else {
			parts[i] = accentText.Render(t[:1]) + secondaryText.Render(t[1:])
		}
	}

	return primaryText.Render("\n scorebug.sh  ") +
		strings.Join(parts, secondaryText.Render(" • "))
}
