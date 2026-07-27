package theme

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	ColorPink            = lipgloss.Color("#FF79C6")
	ColorPurple          = lipgloss.Color("#BD93F9")
	ColorCyan            = lipgloss.Color("#8BE9FD")
	ColorGreen           = lipgloss.Color("#50FA7B")
	ColorRed             = lipgloss.Color("#FF5555")
	ColorYellow          = lipgloss.Color("#F1FA8C")
	ColorSurface         = lipgloss.Color("#44475A")
	ColorForeground      = lipgloss.Color("#F8F8F2")
	ColorForegroundDim   = lipgloss.Color("#6272A4")
	ColorForegroundMuted = lipgloss.Color("#B0B0B0")

	AppBorder            = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorPurple).Padding(0, 1)
	PopupStyle           = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(ColorPink).Padding(1, 2)
	Title                = lipgloss.NewStyle().Bold(true).Foreground(ColorPink)
	Subtitle             = lipgloss.NewStyle().Foreground(ColorPurple)
	NormalText           = lipgloss.NewStyle().Foreground(ColorForeground)
	DimText              = lipgloss.NewStyle().Foreground(ColorSurface)
	SuccessText          = lipgloss.NewStyle().Foreground(ColorGreen)
	ErrorText            = lipgloss.NewStyle().Foreground(ColorRed)
	WarningText          = lipgloss.NewStyle().Foreground(ColorYellow)
	ActiveTab            = lipgloss.NewStyle().Bold(true).Foreground(ColorPink).Background(ColorSurface).Padding(0, 1)
	InactiveTab          = lipgloss.NewStyle().Foreground(ColorPurple).Padding(0, 1)
	SidebarItem          = lipgloss.NewStyle().Foreground(ColorForeground).PaddingLeft(1)
	SidebarItemActive    = lipgloss.NewStyle().Bold(true).Foreground(ColorCyan).Background(ColorSurface).PaddingLeft(1)
	SidebarItemInstalled = lipgloss.NewStyle().Foreground(ColorGreen).PaddingLeft(1)
	HelpStyle            = lipgloss.NewStyle().Foreground(ColorSurface)
	KeyStyle             = lipgloss.NewStyle().Bold(true).Foreground(ColorCyan)
	DescStyle            = lipgloss.NewStyle().Foreground(ColorForeground)
	PanelFocused         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorCyan).Padding(0, 1)
	PanelUnfocused       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorSurface).Padding(0, 1)
)

const (
	IconInstalled    = "✓"
	IconNotInstalled = "○"
	IconGear         = "⚙"
)

var (
	StatusInstalled    = lipgloss.NewStyle().Foreground(ColorGreen).SetString(IconInstalled)
	StatusNotInstalled = lipgloss.NewStyle().Foreground(ColorForegroundDim).SetString(IconNotInstalled)
)

func GetModuleIcon(name string) string {
	icons := map[string]string{"git": "󰊢", "neovim": "", "vim": "", "node": "󰎙", "python3": "󰌠", "powershell": ""}
	if icon := icons[name]; icon != "" {
		return icon
	}
	return "󰏖"
}

func Clip(s string, w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) > h {
		lines = lines[:h]
	}
	for i, line := range lines {
		runes := []rune(line)
		if len(runes) > w {
			lines[i] = string(runes[:w])
		}
	}
	return strings.Join(lines, "\n")
}
