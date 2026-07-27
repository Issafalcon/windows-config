package popup

import (
	lipgloss "charm.land/lipgloss/v2"
	"github.com/issafalcon/windows-config-tui/internal/theme"
	"strings"
)

type Model struct {
	title, content string
	width, height  int
	visible        bool
}

func NewPopup(title, content string, width, height int) Model {
	return Model{title: title, content: content, width: width, height: height}
}
func (m Model) Show() Model     { m.visible = true; return m }
func (m Model) Hide() Model     { m.visible = false; return m }
func (m Model) IsVisible() bool { return m.visible }
func renderPopup(title, content string, width, height, screenWidth, screenHeight int) string {
	separator := theme.DimText.Render(strings.Repeat("─", max(1, width-4)))
	body := lipgloss.JoinVertical(lipgloss.Left, theme.Title.Render(title), separator, "", content)
	return lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, theme.PopupStyle.Width(width).Height(height).Render(body))
}
