package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/issafalcon/windows-config-tui/internal/config"
	"github.com/issafalcon/windows-config-tui/internal/detail"
	"github.com/issafalcon/windows-config-tui/internal/elevate"
	"github.com/issafalcon/windows-config-tui/internal/installer"
	"github.com/issafalcon/windows-config-tui/internal/module"
	"github.com/issafalcon/windows-config-tui/internal/popup"
	"github.com/issafalcon/windows-config-tui/internal/prereqs"
	"github.com/issafalcon/windows-config-tui/internal/sidebar"
	"github.com/issafalcon/windows-config-tui/internal/theme"
	"github.com/issafalcon/windows-config-tui/internal/utils"
)

type AppState int

const (
	StatePrereqCheck AppState = iota
	StateModulesSetup
	StateDashboard
	StateInstalling
)

type FocusArea int

const (
	FocusSidebar FocusArea = iota
	FocusDetail
)

type ProgramReadyMsg struct{ Program *tea.Program }
type Model struct {
	state                                                      AppState
	focus                                                      FocusArea
	width, height                                              int
	ready                                                      bool
	prereq                                                     prereqs.Model
	sidebar                                                    sidebar.Model
	detail                                                     detail.Model
	help                                                       popup.HelpModel
	confirm                                                    popup.ConfirmModel
	script                                                     popup.ScriptModel
	category                                                   popup.CategoryModel
	input                                                      popup.InputModel
	showHelp, showConfirm, showScript, showCategory, showInput bool
	selected                                                   string
	program                                                    *tea.Program
	elevate                                                    elevate.Client
	queue                                                      []string
}

func NewModel() Model {
	return Model{state: StatePrereqCheck, focus: FocusSidebar, prereq: prereqs.New(), sidebar: sidebar.NewModel(items(), 30, 30), detail: detail.NewModel(60, 30), help: popup.NewHelpPopup(nil)}
}
func items() []sidebar.ModuleItem {
	installed, _ := utils.GetInstalledModules()
	set := map[string]bool{}
	for _, n := range installed {
		set[n] = true
	}
	all := module.DefaultRegistry.All()
	out := make([]sidebar.ModuleItem, 0, len(all))
	for _, m := range all {
		out = append(out, sidebar.ModuleItem{Name: m.Name, Icon: m.Icon, Description: m.Description, Category: m.Category, Installed: set[m.Name]})
	}
	return out
}
func (m Model) Init() tea.Cmd { return m.prereq.Init() }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case ProgramReadyMsg:
		m.program = v.Program
		return m, nil
	case tea.WindowSizeMsg:
		m.width, m.height = v.Width, v.Height
		m.ready = true
		m.resize()
		return m, nil
	case prereqs.PrereqsPassedMsg:
		if utils.GetModulesDir() == "" {
			m.state = StateModulesSetup
			m.input = popup.NewInputDialog("Modules directory", "Path to your modules/ folder:")
			m.showInput = true
			return m, nil
		}
		m.reload()
		m.enterDashboard()
		return m, nil
	case popup.InputSubmitMsg:
		if m.state == StateModulesSetup {
			path := expandHome(v.Value)
			if path == "" {
				return m, nil
			}
			if err := os.MkdirAll(path, 0o755); err != nil {
				return m, nil
			}
			if err := config.SetModulesDir(path); err != nil {
				return m, nil
			}
			m.showInput = false
			m.reload()
			m.enterDashboard()
		}
		return m, nil
	case popup.InputCancelMsg:
		return m, nil
	case sidebar.CursorChangedMsg:
		m.selected = v.Name
		m.updateDetail(v.Name)
		return m, nil
	case sidebar.ModuleSelectedMsg:
		m.selected = v.Name
		m.updateDetail(v.Name)
		m.focus = FocusDetail
		m.sidebar.SetFocused(false)
		return m, nil
	case popup.ConfirmYesMsg:
		m.showConfirm = false
		m.detail.SetActiveTab(detail.TabOutput)
		m.detail.OutputModel().Clear()
		queue, err := planInstallQueue(v.ModuleName)
		if err != nil {
			m.detail.OutputModel().AppendLine("✗ Could not plan install: " + err.Error())
			return m, nil
		}
		if len(queue) == 0 {
			m.detail.OutputModel().AppendLine("✓ already installed")
			return m, nil
		}
		m.state = StateInstalling
		m.queue = queue[1:]
		return m.begin(queue[0])
	case popup.ConfirmNoMsg:
		m.showConfirm = false
		return m, nil
	case popup.ConfirmReviewMsg:
		m.showScript = true
		m.script = popup.NewScriptViewer("Scripts — "+v.ModuleName, reviewScripts(v.ModuleName))
		return m, nil
	case popup.ScriptDismissMsg:
		m.showScript = false
		return m, nil
	case popup.CategorySelectedMsg:
		m.showCategory = false
		m.sidebar.SetCategoryFilter(v.Category)
		m.selected = m.sidebar.Selected()
		m.updateDetail(m.selected)
		return m, nil
	case popup.CategoryCancelMsg:
		m.showCategory = false
		return m, nil
	case installer.InstallStartMsg:
		m.detail.OutputModel().SetInstalling(v.ModuleName, true)
		return m, nil
	case installer.InstallOutputMsg:
		m.detail.OutputModel().AppendLine(v.Line)
		return m, nil
	case installer.InstallCompleteMsg:
		m.detail.OutputModel().SetInstalling(v.ModuleName, false)
		if !v.Success {
			m.state = StateDashboard
			m.queue = nil
			m.detail.OutputModel().AppendLine("✗ " + v.ModuleName + " failed: " + v.Error.Error())
			return m, nil
		}
		m.sidebar.SetInstalled(v.ModuleName, true)
		m.detail.OutputModel().AppendLine("✓ " + v.ModuleName + " installed successfully!")
		if len(m.queue) > 0 {
			next := m.queue[0]
			m.queue = m.queue[1:]
			return m.begin(next)
		}
		m.state = StateDashboard
		return m, nil
	case tea.KeyPressMsg:
		if m.showHelp {
			if v.String() == "?" || v.String() == "esc" || v.String() == "q" {
				m.showHelp = false
			}
			return m, nil
		}
		if m.showScript {
			var c tea.Cmd
			m.script, c = m.script.Update(msg)
			return m, c
		}
		if m.showCategory {
			var c tea.Cmd
			m.category, c = m.category.Update(msg)
			return m, c
		}
		if m.showConfirm {
			var c tea.Cmd
			m.confirm, c = m.confirm.Update(msg)
			return m, c
		}
		if m.showInput {
			var c tea.Cmd
			m.input, c = m.input.Update(msg)
			return m, c
		}
		if v.String() == "ctrl+c" || v.String() == "q" {
			m.elevate.Shutdown()
			return m, tea.Quit
		}
		if v.String() == "?" {
			m.showHelp = true
			return m, nil
		}
	}
	if m.state == StatePrereqCheck {
		updated, c := m.prereq.Update(msg)
		m.prereq = updated.(prereqs.Model)
		return m, c
	}
	if m.state == StateDashboard || m.state == StateInstalling {
		return m.updateDashboard(msg)
	}
	return m, nil
}
func (m *Model) enterDashboard() {
	m.state = StateDashboard
	m.selected = m.sidebar.Selected()
	m.updateDetail(m.selected)
}
func (m Model) updateDashboard(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok {
		if m.sidebar.IsSearching() {
			var c tea.Cmd
			m.sidebar, c = m.sidebar.Update(msg)
			return m, c
		}
		switch k.String() {
		case "/", "s":
			m.focus = FocusSidebar
			m.sidebar.SetFocused(true)
			return m, m.sidebar.ActivateSearch()
		case "shift+tab":
			if m.focus == FocusSidebar {
				m.focus = FocusDetail
			} else {
				m.focus = FocusSidebar
			}
			m.sidebar.SetFocused(m.focus == FocusSidebar)
			return m, nil
		case "c":
			m.category = popup.NewCategoryPicker(module.DefaultRegistry.Categories(), m.sidebar.CategoryFilter())
			m.showCategory = true
			return m, nil
		case "i", "enter":
			if m.selected != "" {
				mod, ok := module.DefaultRegistry.Get(m.selected)
				if ok {
					queue, err := planInstallQueue(m.selected)
					list := []string{mod.Name + " — " + mod.Description}
					if err != nil {
						list = append(list, "✗ "+err.Error())
					} else {
						for _, name := range queue {
							if name != m.selected {
								list = append(list, "▸ dependency: "+name)
							}
						}
					}
					has := utils.ModuleScriptExists(m.selected, "install.ps1") || utils.ModuleScriptExists(m.selected, "config.ps1")
					m.confirm = popup.NewConfirmDialog(m.selected, list, has)
					m.showConfirm = true
					return m, nil
				}
			}
		}
	}
	var c tea.Cmd
	if m.focus == FocusSidebar {
		m.sidebar, c = m.sidebar.Update(msg)
	} else {
		m.detail, c = m.detail.Update(msg)
	}
	return m, c
}
func (m Model) begin(name string) (tea.Model, tea.Cmd) {
	mod, ok := module.DefaultRegistry.Get(name)
	if !ok {
		return m, func() tea.Msg {
			return installer.InstallCompleteMsg{ModuleName: name, Error: fmt.Errorf("unknown module")}
		}
	}
	m.detail.OutputModel().AppendLine("▸ Installing " + name + "...")
	return m, installer.RunInstallStreaming(m.program, name, mod.RequiresAdmin, &m.elevate)
}
func (m *Model) reload() {
	_, _ = module.LoadFromDir(utils.GetModulesDir())
	m.sidebar.SetItems(items())
}
func (m *Model) updateDetail(name string) {
	mod, ok := module.DefaultRegistry.Get(name)
	if !ok {
		return
	}
	m.detail.OverviewModel().SetModule(mod.Name, mod.Description, mod.Website, mod.Repo, nil)
	m.detail.ConfigModel().SetModule(mod.Name, nil)
}
func planInstallQueue(target string) ([]string, error) {
	order, err := module.DefaultRegistry.GetInstallOrder(target)
	if err != nil {
		return nil, err
	}
	var todo []string
	for _, mod := range order {
		if !utils.ModuleSatisfied(mod.Name, mod.CheckCommand) {
			todo = append(todo, mod.Name)
		}
	}
	return todo, nil
}
func reviewScripts(name string) string {
	var parts []string
	for _, script := range []string{"install.ps1", "config.ps1"} {
		if data, err := os.ReadFile(utils.ModuleScriptPath(name, script)); err == nil {
			parts = append(parts, "# "+script+"\n"+string(data))
		}
	}
	if len(parts) == 0 {
		return "(no install.ps1 or config.ps1)"
	}
	return strings.Join(parts, "\n\n")
}
func (m *Model) resize() {
	w, h := int(float64(m.width)*.8)-2, int(float64(m.height)*.8)-2
	if w < 60 {
		w = 60
	}
	if h < 20 {
		h = 20
	}
	side := w * 3 / 10
	m.sidebar.SetSize(side-2, h-7)
	m.detail.SetSize(w-side-3, h-7)
}
func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("Initializing...")
	}
	w, h := int(float64(m.width)*.8)-2, int(float64(m.height)*.8)-2
	if w < 60 {
		w = 60
	}
	if h < 20 {
		h = 20
	}
	var content string
	switch m.state {
	case StatePrereqCheck:
		content = m.prereq.View().Content
	case StateModulesSetup:
		content = theme.Title.Render("Modules path setup") + "\n\n" + theme.NormalText.Render("Choose the folder containing your modules.")
	default:
		content = m.dashboard(w, h)
	}
	view := theme.AppBorder.Width(w).Height(h).Render(content)
	if m.showHelp {
		view = m.help.Render(w+2, h+2)
	}
	if m.showConfirm {
		view = m.confirm.Render(w+2, h+2)
	}
	if m.showScript {
		view = m.script.Render(w+2, h+2)
	}
	if m.showCategory {
		view = m.category.Render(w+2, h+2)
	}
	if m.showInput {
		view = m.input.Render(w+2, h+2)
	}
	out := tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view))
	out.AltScreen = true
	return out
}
func (m Model) dashboard(w, h int) string {
	side := w * 3 / 10
	detailW := w - side - 3
	body := lipgloss.JoinHorizontal(lipgloss.Top, theme.PanelFocused.Width(side).Height(h-6).Render(m.sidebar.View()), " ", theme.PanelFocused.Width(detailW).Height(h-6).Render(m.detail.View()))
	return theme.Title.Render("⚡ Windows Config") + "\n\n" + body + "\n\n" + theme.HelpStyle.Render("/: search • i/enter: install • c: category • tab: focus • ?: help • q: quit")
}
func expandHome(path string) string {
	path = strings.TrimSpace(path)
	if path == "~" || strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}
	return path
}
