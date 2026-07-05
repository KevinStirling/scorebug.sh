package playfeed

import (
	"slices"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/KevinStirling/scorebug.sh/internal/ui/components/schedule"
)

const maxWidth = 90

type Model struct {
	list list.Model
	game *schedule.GameSelectedMsg
}

// SetSize sizes the feed, clamping the width to maxPlayFeedWidth. All width
// updates should go through here so the cap can't be bypassed.
func (m *Model) SetSize(width, height int) {
	w := min(width, maxWidth)
	m.list.SetSize(w, max(height, rowHeight))
	m.list.Styles.TitleBar = m.list.Styles.TitleBar.Width(w)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New() Model {
	l := list.New([]list.Item{}, itemDelegate{}, 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Title = "Play-by-Play"
	l.Styles.Title = lipgloss.NewStyle()
	l.Styles.TitleBar = titleBar
	return Model{list: l}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case schedule.GameSelectedMsg:
		m.game = &msg
		m.list.SetItems(m.build())
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}

func (m Model) build() []list.Item {
	if m.game == nil {
		return nil
	}
	plays := slices.Clone(m.game.Bug.Feed.LiveData.Plays.AllPlays)
	slices.Reverse(plays)
	items := make([]list.Item, 0, len(plays))
	for _, p := range plays {
		desc := p.Result.Description
		if desc == "" {
			desc = "Play in progress..."
		}
		items = append(items, item{
			inningSt:    p.About.HalfInning,
			inning:      p.About.Inning,
			eventType:   p.Result.EventType,
			description: desc,
		})
	}
	return items
}
