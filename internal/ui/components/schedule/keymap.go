package schedule

import "charm.land/bubbles/v2/key"

type ScheduleKeyMap struct {
	Up              key.Binding
	Down            key.Binding
	Select          key.Binding
	FilterLive      key.Binding
	FilterScheduled key.Binding
	FilterFinal     key.Binding
	Help            key.Binding
}

var keys = ScheduleKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", "return"),
		key.WithHelp("enter", "select game"),
	),
	FilterLive: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "show live games"),
	),
	FilterScheduled: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "show scheudled games"),
	),
	FilterFinal: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "show final games"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle full help"),
	),
}

func (k ScheduleKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Help}
}

func (k ScheduleKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Help},
		{k.FilterLive, k.FilterScheduled, k.FilterFinal},
	}
}
