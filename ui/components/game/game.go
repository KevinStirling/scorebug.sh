package game

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	GameContent     string
	ContainerWidth  int
	ContainerHeight int
	container       string
	enabled         bool
}

func NewModel() Model {
	return Model{
		GameContent: "test",
		enabled:     true,
	}
}

func (m Model) Init() tea.Cmd {
	// could use this to get more game details upon selection
	return tea.Batch()
}

func (m Model) View() string {
	if !m.enabled {
		return ""
	}
	return m.buildContainer()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.container = m.buildContainer()
	}

	return m, nil
}

func (m Model) buildContainer() string {
	darkerField := newField(m.ContainerHeight, m.ContainerWidth, lipgloss.BrightBlack)
	// lighterField := newField(17, 43, lipgloss.Magenta)
	gameContent := lipgloss.NewStyle().Height(1).Width(4).Foreground(lipgloss.Magenta).Render(m.GameContent)
	bg := lipgloss.NewLayer(darkerField)
	fg := lipgloss.NewLayer(gameContent)
	comp := lipgloss.NewCompositor(bg, fg)

	container := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Green).
		AlignHorizontal(lipgloss.Left).
		MarginTop(2).
		Render(comp.Render())
	return container
}

// newField fills a rectangular area with a given character in a given color.
func newField(rows, cols int, color color.Color) string {
	fieldSetyle := lipgloss.NewStyle().Foreground(color).AlignHorizontal(lipgloss.Left)
	fieldBuilder := strings.Builder{}
	for i := range rows {
		for range cols {
			fieldBuilder.WriteString("#")
		}
		if i < rows-1 {
			fieldBuilder.WriteString("\n")
		}
	}
	return fieldSetyle.Render(fieldBuilder.String())
}
