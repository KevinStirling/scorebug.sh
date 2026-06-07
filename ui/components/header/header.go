package header

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/KevinStirling/scorebug.sh/ui/components/theme"
)

type Model struct {
	title     string
	tabs      []string
	ActiveTab int
}

func New() Model {
	return Model{
		title:     "\n scorebug.sh ",
		tabs:      []string{"live", "scheduled", "final"},
		ActiveTab: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Model:
		m.ActiveTab = msg.ActiveTab
	}

	return m, nil

}

func (m Model) Render() string {
	parts := make([]string, len(m.tabs))
	for i, t := range m.tabs {
		if i == m.ActiveTab {
			parts[i] = theme.AccentText.Render(t)
		} else {
			parts[i] = theme.AccentText.Render(t[:1]) + theme.SecondaryText.Render(t[1:])
		}
	}
	return theme.PrimaryText.Render("\n scorebug.sh ") +
		strings.Join(parts, theme.SecondaryText.Render(" • "))
}
