package game

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/KevinStirling/scorebug.sh/internal/ui/components/schedule"
)

type item struct {
	description string
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	if _, err := fmt.Fprint(w, i); err != nil {
		log.Error("encoutnered error writing list item", "error", err)
	}
}

type PlayFeed struct {
	list list.Model
	game *schedule.GameSelectedMsg
}

func (m PlayFeed) Init() tea.Cmd {
	return nil
}

func NewPlayfeed() PlayFeed {
	l := list.New([]list.Item{}, itemDelegate{}, 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.Title = "Play-by-Play"
	return PlayFeed{list: l}
}

func (m PlayFeed) Update(msg tea.Msg) (PlayFeed, tea.Cmd) {
	switch msg := msg.(type) {
	case schedule.GameSelectedMsg:
		m.game = &msg
		m.list.SetItems(m.buildPlayFeed())
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m PlayFeed) View() string {
	return m.list.View()
}

func (m PlayFeed) buildPlayFeed() []list.Item {
	if m.game == nil {
		return nil
	}
	plays := m.game.Bug.Feed.LiveData.Plays.AllPlays
	items := make([]list.Item, 0, len(plays))
	for _, p := range plays {
		items = append(items, item{
			description: p.Result.Description,
		})
	}
	return items
}
