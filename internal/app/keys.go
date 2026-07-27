package app

import "charm.land/bubbles/v2/key"

// KeyMap centralizes shortcuts so the dashboard and its help surface agree.
type KeyMap struct {
	Install, Reinstall, Search, FilterCategory, Help, Quit key.Binding
	SwitchTab                                              key.Binding
}

var DefaultKeyMap = KeyMap{
	Install: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "install module")),
	Reinstall: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "re-run install"),
	),
	Search: key.NewBinding(key.WithKeys("s", "/"), key.WithHelp("s", "search modules")),
	FilterCategory: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "filter by category"),
	),
	Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	SwitchTab: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch tab")),
}
