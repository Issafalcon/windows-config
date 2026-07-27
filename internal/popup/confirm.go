package popup

import (
	tea "charm.land/bubbletea/v2"
	"github.com/issafalcon/windows-config-tui/internal/theme"
	"strings"
)

type ConfirmAction string

const (
	ActionInstall   ConfirmAction = "install"
	ActionUninstall ConfirmAction = "uninstall"
)

type ConfirmYesMsg struct {
	ModuleName string
	Action     ConfirmAction
}
type ConfirmNoMsg struct{}
type ConfirmReviewMsg struct {
	ModuleName string
	Action     ConfirmAction
}
type ConfirmModel struct {
	Model
	module    string
	action    ConfirmAction
	items     []string
	hasScript bool
	cursor    int
}

func NewConfirmDialog(module string, items []string, hasScript bool) ConfirmModel {
	return ConfirmModel{Model: NewPopup("Install "+module+"?", "", 54, 16).Show(), module: module, action: ActionInstall, items: items, hasScript: hasScript}
}
func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	maxCursor := 1
	if m.hasScript {
		maxCursor = 2
	}
	switch k.String() {
	case "left", "h":
		if m.cursor > 0 {
			m.cursor--
		}
	case "right", "l", "tab":
		if m.cursor < maxCursor {
			m.cursor++
		}
	case "esc":
		return m, func() tea.Msg { return ConfirmNoMsg{} }
	case "r":
		if m.hasScript {
			return m, func() tea.Msg { return ConfirmReviewMsg{m.module, m.action} }
		}
	case "enter":
		if m.cursor == 0 {
			return m, func() tea.Msg { return ConfirmYesMsg{m.module, m.action} }
		}
		if m.cursor == 2 && m.hasScript {
			return m, func() tea.Msg { return ConfirmReviewMsg{m.module, m.action} }
		}
		return m, func() tea.Msg { return ConfirmNoMsg{} }
	}
	return m, nil
}
func (m ConfirmModel) Render(w, h int) string {
	if !m.visible {
		return ""
	}
	buttons := "[Yes]  No"
	if m.cursor == 1 {
		buttons = "Yes  [No]"
	}
	if m.hasScript {
		buttons += "  Review"
	}
	return renderPopup(m.title, theme.NormalText.Render(strings.Join(m.items, "\n"))+"\n\n"+theme.HelpStyle.Render(buttons+"  ←/→ select · enter confirm"), m.width, m.height, w, h)
}
