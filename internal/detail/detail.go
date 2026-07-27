package detail

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/issafalcon/windows-config-tui/internal/theme"
	"strings"
)

type Tab int

const (
	TabOverview Tab = iota
	TabOutput
	TabConfig
)

type DepStatus struct {
	Name, Method        string
	Installed, Checking bool
}
type ConfigOption struct {
	Name, Description, Default string
	Choices                    []string
	Selected                   string
}
type Model struct {
	active                           Tab
	width, height                    int
	name, description, website, repo string
	deps                             []DepStatus
	lines                            []string
	installing                       bool
	installingModule                 string
	config                           []ConfigOption
}

func NewModel(w, h int) Model { return Model{width: w, height: h} }
func (m Model) Init() tea.Cmd { return nil }
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == "tab" {
		m.active = (m.active + 1) % 3
	}
	return m, nil
}
func (m Model) View() string {
	tabs := []string{"Overview", "Output", "Configuration"}
	for i := range tabs {
		if Tab(i) == m.active {
			tabs[i] = theme.ActiveTab.Render(tabs[i])
		} else {
			tabs[i] = theme.InactiveTab.Render(tabs[i])
		}
	}
	var body string
	switch m.active {
	case TabOverview:
		body = m.overview()
	case TabOutput:
		body = m.output()
	case TabConfig:
		body = m.configuration()
	}
	return theme.Clip(strings.Join([]string{lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...), strings.Repeat("─", m.width), body}, "\n"), m.width, m.height)
}
func (m Model) overview() string {
	if m.name == "" {
		return theme.DimText.Render("Select a module from the sidebar")
	}
	s := theme.Title.Render(theme.GetModuleIcon(m.name)+" "+m.name) + "\n" + theme.NormalText.Render(m.description)
	if m.website != "" {
		s += "\n\n" + theme.Subtitle.Render("Website: ") + m.website
	}
	if m.repo != "" {
		s += "\n" + theme.Subtitle.Render("Repo: ") + m.repo
	}
	if len(m.deps) > 0 {
		s += "\n\n" + theme.Subtitle.Render("Dependencies:")
		for _, d := range m.deps {
			status := theme.IconNotInstalled
			if d.Installed {
				status = theme.IconInstalled
			}
			s += "\n  " + status + " " + d.Name + " " + theme.DimText.Render("("+d.Method+")")
		}
	}
	return s
}
func (m Model) output() string {
	if len(m.lines) == 0 {
		return theme.DimText.Render("No installation output yet.\nSelect a module and press i to install.")
	}
	prefix := ""
	if m.installing {
		prefix = theme.Subtitle.Render("Installing: "+m.installingModule) + "\n"
	}
	return prefix + strings.Join(m.lines, "\n")
}
func (m Model) configuration() string {
	if m.name == "" {
		return theme.DimText.Render("Select a module from the sidebar")
	}
	if len(m.config) == 0 {
		return theme.DimText.Render("No configuration options for " + m.name + ".")
	}
	var b strings.Builder
	for _, o := range m.config {
		b.WriteString(o.Name + ": " + o.Selected + "\n")
	}
	return b.String()
}
func (m *Model) SetSize(w, h int)              { m.width = w; m.height = h }
func (m *Model) SetActiveTab(t Tab)            { m.active = t }
func (m *Model) OverviewModel() *OverviewModel { return (*OverviewModel)(m) }
func (m *Model) OutputModel() *OutputModel     { return (*OutputModel)(m) }
func (m *Model) ConfigModel() *ConfigModel     { return (*ConfigModel)(m) }

type OverviewModel Model

func (m *OverviewModel) SetModule(name, description, website, repo string, deps []DepStatus) {
	(*Model)(m).name = name
	(*Model)(m).description = description
	(*Model)(m).website = website
	(*Model)(m).repo = repo
	(*Model)(m).deps = deps
}

type OutputModel Model

func (m *OutputModel) AppendLine(line string) { (*Model)(m).lines = append((*Model)(m).lines, line) }
func (m *OutputModel) Clear()                 { (*Model)(m).lines = nil }
func (m *OutputModel) SetInstalling(name string, installing bool) {
	(*Model)(m).installing = installing
	(*Model)(m).installingModule = name
}

type ConfigModel Model

func (m *ConfigModel) SetModule(name string, opts []ConfigOption) {
	(*Model)(m).name = name
	(*Model)(m).config = opts
}
