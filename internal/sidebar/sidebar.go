package sidebar

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/issafalcon/windows-config-tui/internal/theme"
)

type ModuleSelectedMsg struct{ Name string }
type CursorChangedMsg struct{ Name string }

type Model struct {
	items, filtered       []ModuleItem
	category              string
	cursor, width, height int
	search                bool
	input                 textinput.Model
	selected              string
	focused               bool
}

func NewModel(items []ModuleItem, width, height int) Model {
	input := textinput.New()
	input.Placeholder = "Type to filter modules..."
	input.CharLimit = 50
	m := Model{items: items, filtered: items, width: width, height: height, input: input, focused: true}
	if len(items) > 0 {
		m.selected = items[0].Name
	}
	return m
}
func (m Model) Init() tea.Cmd { return nil }
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	if m.search {
		switch key.String() {
		case "esc":
			m.search = false
			m.input.Blur()
			return m, nil
		case "enter":
			return m, m.selectCmd()
		case "up", "ctrl+p":
			m.move(-1)
			return m, m.cursorCmd()
		case "down", "ctrl+n":
			m.move(1)
			return m, m.cursorCmd()
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		m.refilter()
		return m, tea.Batch(cmd, m.cursorCmd())
	}
	switch key.String() {
	case "s", "/":
		m.search = true
		return m, m.input.Focus()
	case "esc":
		m.input.SetValue("")
		m.refilter()
		return m, m.cursorCmd()
	case "j", "down":
		m.move(1)
		return m, m.cursorCmd()
	case "k", "up":
		m.move(-1)
		return m, m.cursorCmd()
	case "enter":
		return m, m.selectCmd()
	}
	return m, nil
}
func (m *Model) move(delta int) {
	if len(m.filtered) == 0 {
		return
	}
	m.cursor = (m.cursor + delta + len(m.filtered)) % len(m.filtered)
	m.selected = m.filtered[m.cursor].Name
}
func (m Model) cursorCmd() tea.Cmd {
	return func() tea.Msg { return CursorChangedMsg{Name: m.selected} }
}
func (m Model) selectCmd() tea.Cmd {
	if m.selected == "" {
		return nil
	}
	return func() tea.Msg { return ModuleSelectedMsg{Name: m.selected} }
}
func (m *Model) refilter() {
	m.filtered = FilterModules(FilterByCategory(m.items, m.category), m.input.Value())
	m.cursor = 0
	m.selected = ""
	if len(m.filtered) > 0 {
		m.selected = m.filtered[0].Name
	}
}
func (m Model) Selected() string                   { return m.selected }
func (m Model) IsSearching() bool                  { return m.search }
func (m *Model) ActivateSearch() tea.Cmd           { m.search = true; m.focused = true; return m.input.Focus() }
func (m *Model) SetFocused(f bool)                 { m.focused = f }
func (m *Model) SetSize(w, h int)                  { m.width = w; m.height = h }
func (m *Model) SetYOffset(int)                    {}
func (m Model) CategoryFilter() string             { return m.category }
func (m *Model) SetCategoryFilter(category string) { m.category = category; m.refilter() }
func (m *Model) SetItems(items []ModuleItem)       { m.items = items; m.refilter() }
func (m *Model) SetInstalled(name string, installed bool) {
	for i := range m.items {
		if m.items[i].Name == name {
			m.items[i].Installed = installed
		}
	}
	m.refilter()
}
func (m Model) View() string {
	category := "All"
	if m.category != "" {
		category = m.category
	}
	var b strings.Builder
	b.WriteString(theme.DimText.Render("Category: "+category+"  (c to change)") + "\n")
	if m.search {
		b.WriteString(m.input.View() + "\n")
	} else {
		b.WriteString(theme.DimText.Render("/ search") + "\n")
	}
	if len(m.filtered) == 0 {
		return b.String() + theme.DimText.Render("No modules match your search.")
	}
	max := m.height / 4
	if max < 1 {
		max = 1
	}
	start := m.cursor - max/2
	if start < 0 {
		start = 0
	}
	end := start + max
	if end > len(m.filtered) {
		end = len(m.filtered)
	}
	for i := start; i < end; i++ {
		item := m.filtered[i]
		icon := item.Icon
		if icon == "" {
			icon = theme.GetModuleIcon(item.Name)
		}
		status := theme.IconNotInstalled + " not installed"
		if item.Installed {
			status = theme.IconInstalled + " installed"
		}
		line := fmt.Sprintf("%s %s\n  %s\n  %s", icon, item.Name, item.Description, status)
		if i == m.cursor {
			line = theme.SidebarItemActive.Render(line)
		} else {
			line = theme.SidebarItem.Render(line)
		}
		b.WriteString(line + "\n")
	}
	return theme.Clip(b.String(), m.width, m.height)
}
