package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	start  key.Binding
	stop   key.Binding
	reset  key.Binding
	change key.Binding
	quit   key.Binding
}

var helpTimerKeys = keymap{
	start: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start"),
	),
	stop: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "stop "),
	),
	reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
	),
	change: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "change timer value"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
