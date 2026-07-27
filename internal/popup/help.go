package popup

import (
	"github.com/issafalcon/windows-config-tui/internal/theme"
	"strings"
)

type HelpBinding struct{ Key, Description string }
type HelpModel struct {
	Model
	bindings []HelpBinding
}

func NewHelpPopup(bindings []HelpBinding) HelpModel {
	if bindings == nil {
		bindings = []HelpBinding{{"j/k", "navigate"}, {"enter", "select/install"}, {"i", "install"}, {"c", "category"}, {"/", "search"}, {"tab", "switch panel/tab"}, {"?", "help"}, {"q", "quit"}}
	}
	return HelpModel{Model: NewPopup("Keyboard Shortcuts", "", 55, 20).Show(), bindings: bindings}
}
func (m HelpModel) Render(w, h int) string {
	if !m.visible {
		return ""
	}
	var lines []string
	for _, b := range m.bindings {
		lines = append(lines, theme.KeyStyle.Render(b.Key)+"  "+theme.DescStyle.Render(b.Description))
	}
	return renderPopup(m.title, strings.Join(lines, "\n"), m.width, m.height, w, h)
}
