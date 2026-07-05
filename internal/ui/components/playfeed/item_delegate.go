package playfeed

import (
	"fmt"
	"io"
	"strconv"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/KevinStirling/scorebug.sh/data"
)

type item struct {
	inningSt    string
	inning      int
	eventType   string
	description string
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

const rowHeight = 3

func (d itemDelegate) Height() int                             { return rowHeight }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok || i.description == "" {
		return
	}

	style := evenRow
	if index%2 == 1 {
		style = oddRow
	}

	inning := eventInning.Render(fmt.Sprintf("%s%s ", data.RenderInningState(i.inningSt), strconv.Itoa(i.inning)))
	event := lipgloss.JoinHorizontal(lipgloss.Left, inning, style.Render(i.description))

	width := m.Width()
	str := style.
		Width(width).
		Height(rowHeight).
		MaxHeight(rowHeight).
		Render(event)

	if _, err := fmt.Fprint(w, str); err != nil {
		log.Error("encountered error writing list item", "error", err)
	}
}
