package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/issafalcon/windows-config-tui/internal/app"
)

func main() {
	m := app.NewModel()
	p := tea.NewProgram(m)
	go p.Send(app.ProgramReadyMsg{Program: p})
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "windows-config-tui:", err)
		os.Exit(1)
	}
}
