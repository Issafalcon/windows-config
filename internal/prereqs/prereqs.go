package prereqs

import (
	"fmt"
	"os/exec"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/issafalcon/windows-config-tui/internal/theme"
	"github.com/issafalcon/windows-config-tui/internal/utils"
)

type PrereqsPassedMsg struct{}

type checklistMsg struct {
	git  bool
	pwsh bool
}

type Model struct {
	width, height int
	git           bool
	pwsh          bool
	checked       bool
	help          help.Model
}

type keyMap struct{}

var quitKey = key.NewBinding(
	key.WithKeys("q", "ctrl+c"),
	key.WithHelp("q", "quit"),
)

func (keyMap) ShortHelp() []key.Binding { return []key.Binding{quitKey} }
func (keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{quitKey}}
}

func New() Model {
	return Model{help: help.New()}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		_, gitErr := exec.LookPath("git")
		_, pwshErr := utils.FindPwsh()
		return checklistMsg{git: gitErr == nil, pwsh: pwshErr == nil}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		if key.Matches(msg, quitKey) {
			return m, tea.Quit
		}
	case checklistMsg:
		m.git, m.pwsh, m.checked = msg.git, msg.pwsh, true
		if m.git && m.pwsh {
			return m, func() tea.Msg { return PrereqsPassedMsg{} }
		}
	}
	return m, nil
}

// View returns checklist content for the root app to embed in AppBorder.
func (m Model) View() string {
	row := func(name string, ok bool) string {
		status := theme.DimText.Render("checking…")
		if m.checked {
			if ok {
				status = theme.SuccessText.Render("✓ installed")
			} else {
				status = theme.ErrorText.Render("✗ missing")
			}
		}
		return fmt.Sprintf("  %-6s %s", name, status)
	}

	return strings.Join([]string{
		theme.Title.Render("Prerequisites"),
		theme.DimText.Render("Git for Windows + PowerShell 7 (pwsh)"),
		"",
		row("git", m.git),
		row("pwsh", m.pwsh),
		"",
		theme.HelpStyle.Render(m.help.View(keyMap{})),
	}, "\n")
}
