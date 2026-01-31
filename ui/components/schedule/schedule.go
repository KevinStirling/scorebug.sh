package schedule

import (
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

var divider = lipgloss.NewStyle().Padding(0, 1)
var primaryText = lipgloss.NewStyle().Foreground(lipgloss.Color("253"))
var secondaryText = lipgloss.NewStyle().Foreground(lipgloss.Color("247"))

type Model struct {
	games     data.Schedule
	paginator paginator.Model
	err       error
}

func NewModel() Model {
	games := data.BuildSchedule(data.GetSchedule())
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(games.Games))
	return Model{
		paginator: p,
		games:     games,
	}
}

type tickMsg time.Time

func tickAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func checkServer() tea.Cmd {
	return func() tea.Msg {
		return data.BuildSchedule(data.GetSchedule())
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(checkServer(), tickAfter(10*time.Second))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case data.Schedule:
		m.games = msg
		return m, nil

	case tickMsg:
		return m, tea.Batch(checkServer(), tickAfter(10*time.Second))

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
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
