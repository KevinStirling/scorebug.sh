package schedule

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/KevinStirling/scorebug.sh/internal/ui/components/scorebug"
)

type scorebugDelegate struct{}

func (d scorebugDelegate) Height() int                               { return scorebug.SB_HEIGHT }
func (d scorebugDelegate) Spacing() int                              { return 0 }
func (d scorebugDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d scorebugDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it, ok := item.(ScorebugItem)
	if !ok {
		return
	}
	style := scorebug.Border
	if index == m.Index() {
		style = scorebug.SelectedBorder
	}
	s := scorebug.Render(it.bug, style)
	if _, err := fmt.Fprint(w, s); err != nil {
		log.Error("render error", "error", err)
	}
}
