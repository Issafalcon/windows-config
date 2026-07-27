package prereqs

import (
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
	git     bool
	pwsh    bool
	checked bool
	help    help.Model
}

type keyMap struct{}

var quitKey = key.NewBinding(
	key.WithKeys("q", "ctrl+c"),
	key.WithHelp("q", "quit"),
)

func (keyMap) ShortHelp() []key.Binding {
	return []key.Binding{quitKey}
}

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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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

func (m Model) View() tea.View {
	status := func(ok bool) string {
		if !m.checked {
			return theme.DimText.Render("checking")
		}
		if ok {
			return theme.SuccessText.Render("✓ installed")
		}
		return theme.ErrorText.Render("✗ missing")
	}
	body := strings.Join([]string{
		theme.Title.Render("Prerequisites"),
		"",
		"  git   " + status(m.git),
		"  pwsh  " + status(m.pwsh),
		"",
		theme.HelpStyle.Render(m.help.View(keyMap{})),
	}, "\n")
	return tea.NewView(theme.PanelFocused.Render(body))
}
