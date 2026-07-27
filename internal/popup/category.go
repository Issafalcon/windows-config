package popup

import (
	tea "charm.land/bubbletea/v2"
	"github.com/issafalcon/windows-config-tui/internal/theme"
	"strings"
)

type CategorySelectedMsg struct{ Category string }
type CategoryCancelMsg struct{}
type CategoryModel struct {
	Model
	options []string
	cursor  int
}

func NewCategoryPicker(categories []string, current string) CategoryModel {
	opts := append([]string{"All"}, categories...)
	c := 0
	for i, v := range opts {
		if v == current {
			c = i
		}
	}
	return CategoryModel{Model: NewPopup("Filter by category", "", 36, min(12, len(opts)+4)).Show(), options: opts, cursor: c}
}
func (m CategoryModel) Update(msg tea.Msg) (CategoryModel, tea.Cmd) {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch k.String() {
	case "up", "k":
		m.cursor = max(0, m.cursor-1)
	case "down", "j":
		m.cursor = min(len(m.options)-1, m.cursor+1)
	case "esc", "q":
		return m, func() tea.Msg { return CategoryCancelMsg{} }
	case "enter":
		c := m.options[m.cursor]
		if c == "All" {
			c = ""
		}
		return m, func() tea.Msg { return CategorySelectedMsg{c} }
	}
	return m, nil
}
func (m CategoryModel) Render(w, h int) string {
	if !m.visible {
		return ""
	}
	var lines []string
	for i, o := range m.options {
		if i == m.cursor {
			o = "▸ " + o
		} else {
			o = "  " + o
		}
		lines = append(lines, o)
	}
	return renderPopup(m.title, theme.NormalText.Render(strings.Join(lines, "\n")), m.width, m.height, w, h)
}
