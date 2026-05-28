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
	paginator  paginator.Model
	tabs       []string
	tabContent [3][]string
	activeTab  int
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
	p.PerPage = 10

	// TODO fix adaptiveActive color... don't think i'm using it right
	p.ActiveDot = lipgloss.NewStyle().Foreground(adaptiveActive).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(adaptiveInactive).Render("•")
	p.SetTotalPages(len(bugs))
	return Model{
		client:    client,
		paginator: p,
		games:     bugs,
		date:      &d,
		tabs:      []string{"live", "scheduled", "final"},
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
		for i := range m.tabs {
			m.tabContent[i] = renderTab(m.games, i)
		}
		m.syncPaginator()
		return m, nil

	case tickMsg:
		m.syncPaginator()
		return m, tea.Batch(m.checkServer(), tickAfter(10*time.Second))

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "l":
			m.activeTab = 0
		case "s":
			m.activeTab = 1
		case "f":
			m.activeTab = 2
		}
	}
	m.syncPaginator()
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	g := renderTab(m.games, m.activeTab)
	var b strings.Builder
	b.WriteString(renderHeader(m.tabs, m.activeTab))
	start, end := m.paginator.GetSliceBounds(len(g))
	for _, item := range g[start:end] {
		b.WriteString("\n" + item)
	}
	b.WriteString("\n " + m.paginator.View())
	b.WriteString(secondaryText.Render("\n\n n/p ←/→ page • q: quit • l: live • s: scheduled • f: final\n"))

	v := tea.NewView(divider.Render(b.String()))
	v.AltScreen = true

	return v
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
	content := m.tabContent[m.activeTab]
	if len(content) == 0 {
		m.paginator.TotalPages = 1
		m.paginator.Page = 0
		return
	}
	m.paginator.SetTotalPages(len(content))
	if m.paginator.Page > m.paginator.TotalPages-1 {
		m.paginator.Page = m.paginator.TotalPages - 1
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

	return primaryText.Render("\nscorebug.sh  ") +
		strings.Join(parts, secondaryText.Render(" • ")) +
		secondaryText.Render("\n"+strings.Repeat("‾", scorebug.SB_WIDTH))
}
