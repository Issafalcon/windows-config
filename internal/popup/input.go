package popup

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/issafalcon/windows-config-tui/internal/theme"
)

type InputSubmitMsg struct{ Value string }
type InputCancelMsg struct{}
type InputModel struct {
	Model
	prompt string
	input  textinput.Model
}

func NewInputDialog(title, prompt string) InputModel {
	in := textinput.New()
	in.Placeholder = "Type here..."
	in.CharLimit = 256
	in.Focus()
	return InputModel{Model: NewPopup(title, "", 50, 10).Show(), prompt: prompt, input: in}
}
func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok {
		switch k.String() {
		case "enter":
			v := m.input.Value()
			return m, func() tea.Msg { return InputSubmitMsg{v} }
		case "esc":
			return m, func() tea.Msg { return InputCancelMsg{} }
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
func (m InputModel) Render(w, h int) string {
	if !m.visible {
		return ""
	}
	return renderPopup(m.title, theme.NormalText.Render(m.prompt)+"\n\n"+m.input.View()+"\n\n"+theme.HelpStyle.Render("enter submit · esc cancel"), m.width, m.height, w, h)
}
