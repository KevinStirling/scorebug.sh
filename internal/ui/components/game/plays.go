package game

import (
	"fmt"
	"io"
	"slices"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/KevinStirling/scorebug.sh/internal/ui/components/schedule"
)

type item struct {
	description string
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

const playRowHeight = 2

func (d itemDelegate) Height() int                             { return playRowHeight }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok || i.description == "" {
		return
	}

	style := playFeedEven
	if index%2 == 1 {
		style = playFeedOdd
	}

	width := m.Width()
	str := style.
		Width(width).
		Height(playRowHeight).
		MaxHeight(playRowHeight).
		Render(i.description)

	if _, err := fmt.Fprint(w, str); err != nil {
		log.Error("encountered error writing list item", "error", err)
	}
}

const maxPlayFeedWidth = 90

type PlayFeed struct {
	list list.Model
	game *schedule.GameSelectedMsg
}

// SetSize sizes the feed, clamping the width to maxPlayFeedWidth. All width
// updates should go through here so the cap can't be bypassed.
func (m *PlayFeed) SetSize(width, height int) {
	m.list.SetSize(min(width, maxPlayFeedWidth), height)
}

func (m PlayFeed) Init() tea.Cmd {
	return nil
}

func NewPlayfeed() PlayFeed {
	l := list.New([]list.Item{}, itemDelegate{}, 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.Title = "Play-by-Play"
	l.Styles.TitleBar = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
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
	plays := slices.Clone(m.game.Bug.Feed.LiveData.Plays.AllPlays)
	slices.Reverse(plays)
	items := make([]list.Item, 0, len(plays))
	for _, p := range plays {
		items = append(items, item{
			description: p.Result.Description,
		})
	}
	return items
}
