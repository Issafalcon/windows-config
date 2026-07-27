package popup

import (
	tea "charm.land/bubbletea/v2"
	"github.com/issafalcon/windows-config-tui/internal/theme"
	"strings"
)

type ScriptDismissMsg struct{}
type ScriptModel struct {
	Model
	lines  []string
	offset int
}

func NewScriptViewer(title, content string) ScriptModel {
	return ScriptModel{Model: NewPopup(title, "", 72, 22).Show(), lines: strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")}
}
func (m ScriptModel) Update(msg tea.Msg) (ScriptModel, tea.Cmd) {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	visible := max(3, m.height-3)
	limit := max(0, len(m.lines)-visible)
	switch k.String() {
	case "esc", "q":
		return m, func() tea.Msg { return ScriptDismissMsg{} }
	case "up", "k":
		m.offset = max(0, m.offset-1)
	case "down", "j":
		m.offset = min(limit, m.offset+1)
	case "pgup":
		m.offset = max(0, m.offset-visible)
	case "pgdown":
		m.offset = min(limit, m.offset+visible)
	}
	return m, nil
}
func (m ScriptModel) Render(w, h int) string {
	if !m.visible {
		return ""
	}
	end := min(len(m.lines), m.offset+max(3, m.height-3))
	return renderPopup(m.title, theme.DimText.Render(strings.Join(m.lines[m.offset:end], "\n"))+"\n\n"+theme.HelpStyle.Render("j/k scroll · esc back"), m.width, m.height, w, h)
}
