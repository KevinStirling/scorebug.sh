package schedule

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/KevinStirling/scorebug.sh/data"
	"github.com/KevinStirling/scorebug.sh/internal/snapshots"
)

type ScorebugItem struct {
	bug data.ScoreBug
}

func (s ScorebugItem) FilterValue() string { return s.bug.HomeAbbr + " " + s.bug.AwayAbbr }

type TabChangedMsg int

func tabChanged(t int) tea.Cmd {
	return func() tea.Msg { return TabChangedMsg(t) }
}

type tickMsg time.Time
type scorebugMsg []data.ScoreBug
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type Model struct {
	list      list.Model
	client    ScheduleClient
	games     []data.ScoreBug
	date      time.Time
	err       error
	tabs      []string
	Keys      ScheduleKeyMap
	ActiveTab int
}

func (m Model) IsFiltering() bool { return m.list.SettingFilter() }

func NewModel(client ScheduleClient) Model {
	now := time.Now()
	d := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	items := make([]list.Item, 0)
	l := list.New(items, scorebugDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	return Model{
		list:      l,
		client:    client,
		date:      d,
		tabs:      []string{"live", "scheduled", "final"},
		Keys:      keys,
		ActiveTab: 0,
	}

}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.checkServer(), tickAfter(10*time.Second))
}

func (m *Model) SetSize(width, height int) {
	h, v := listStyle.GetFrameSize()
	m.list.SetSize(width-h, height-v)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case scorebugMsg:
		m.games = msg
		m.err = nil
		cmd := m.list.SetItems(buildTab(msg, m.ActiveTab))
		return m, cmd
	case errMsg:
		m.err = msg.err
		log.Error("schedule", "error", m.err)
		return m, nil
	case tickMsg:
		return m, tea.Batch(m.checkServer(), tickAfter(10*time.Second))
	case tea.KeyPressMsg:
		if m.list.SettingFilter() {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
		switch {
		case key.Matches(msg, m.Keys.FilterLive):
			m.ActiveTab = 0
		case key.Matches(msg, m.Keys.FilterScheduled):
			m.ActiveTab = 1
		case key.Matches(msg, m.Keys.FilterFinal):
			m.ActiveTab = 2
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
		return m, tea.Batch(
			m.list.SetItems(buildTab(m.games, m.ActiveTab)),
			tabChanged(m.ActiveTab),
		)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return listStyle.Render("error: " + m.err.Error() + "\n(retrying...)")
	}
	return listStyle.Render(m.list.View())
}

func fetchScoreBugs(client ScheduleClient, date time.Time) ([]data.ScoreBug, error) {
	sched, err := client.Schedule(date)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %w", err)
	}

	snaps, err := snapshots.Build(client, sched)
	if err != nil {
		return nil, fmt.Errorf("failed to build snapshots: %w", err)
	}

	return data.BuildScoreBugs(snaps), nil
}

// Returns an array of ScoreBugItems for a given tab
func buildTab(bugs []data.ScoreBug, tab int) []list.Item {
	items := make([]list.Item, 0, len(bugs))
	for _, bug := range bugs {
		switch tab {
		case 0:
			if bug.Status == "Live" {
				items = append(items, ScorebugItem{bug: bug})
			}
		case 1:
			if bug.Status == "Preview" {
				items = append(items, ScorebugItem{bug: bug})
			}
		case 2:
			if bug.Status == "Final" || bug.Status == "Other" {
				items = append(items, ScorebugItem{bug: bug})
			}
		}
	}
	return items
}

func (m Model) checkServer() tea.Cmd {
	return func() tea.Msg {
		bugs, err := fetchScoreBugs(m.client, m.date)
		if err != nil {
			return errMsg{err}
		}
		return scorebugMsg(bugs)
	}
}

func tickAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}
