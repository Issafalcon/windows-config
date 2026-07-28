package app

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
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
	all := module.DefaultRegistry.All()
	out := make([]sidebar.ModuleItem, 0, len(all))
	for _, m := range all {
		out = append(out, sidebar.ModuleItem{
			Name:        m.Name,
			Icon:        m.Icon,
			Description: m.Description,
			Category:    m.Category,
			Installed:   utils.ModuleSatisfied(m.Name, m.CheckCommand),
		})
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
			path := config.ExpandHome(v.Value)
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

		var queue []string
		var err error
		if v.Action == popup.ActionReinstall {
			queue = []string{v.ModuleName}
			m.detail.OutputModel().AppendLine(
				fmt.Sprintf("Force re-run: %s (deps skipped)", v.ModuleName))
		} else {
			queue, err = planInstallQueue(v.ModuleName)
			if err != nil {
				m.detail.OutputModel().AppendLine("✗ Could not plan install: " + err.Error())
				return m, nil
			}
			if len(queue) == 0 {
				m.detail.OutputModel().AppendLine(
					fmt.Sprintf("✓ %s and its dependencies are already satisfied", v.ModuleName))
				m.detail.OutputModel().AppendLine(
					"  Tip: press r to force re-run this module's install scripts")
				m.sidebar.SetInstalled(v.ModuleName, true)
				return m, nil
			}
		}

		m.state = StateInstalling
		m.queue = queue[1:]
		m.detail.OutputModel().AppendLine(
			fmt.Sprintf("Install plan (%d): %s", len(queue), strings.Join(queue, " → ")))
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
		if m.selected != "" {
			m.updateDetail(m.selected)
		}
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
		// Search owns ordinary keystrokes, including q, while it is active.
		if m.state == StateDashboard && m.sidebar.IsSearching() {
			if v.String() == "ctrl+c" {
				m.elevate.Shutdown()
				return m, tea.Quit
			}
			var c tea.Cmd
			m.sidebar, c = m.sidebar.Update(msg)
			return m, c
		}
		if key.Matches(v, DefaultKeyMap.Quit) {
			m.elevate.Shutdown()
			return m, tea.Quit
		}
		if key.Matches(v, DefaultKeyMap.Help) {
			m.showHelp = true
			return m, nil
		}
	}
	if m.state == StatePrereqCheck {
		var c tea.Cmd
		m.prereq, c = m.prereq.Update(msg)
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
		switch {
		case key.Matches(k, DefaultKeyMap.Search):
			m.focus = FocusSidebar
			m.sidebar.SetFocused(true)
			return m, m.sidebar.ActivateSearch()
		case k.String() == "shift+tab":
			if m.focus == FocusSidebar {
				m.focus = FocusDetail
			} else {
				m.focus = FocusSidebar
			}
			m.sidebar.SetFocused(m.focus == FocusSidebar)
			return m, nil
		case key.Matches(k, DefaultKeyMap.FilterCategory):
			m.category = popup.NewCategoryPicker(module.DefaultRegistry.Categories(), m.sidebar.CategoryFilter())
			m.showCategory = true
			return m, nil
		case key.Matches(k, DefaultKeyMap.Install):
			if m.selected != "" {
				mod, ok := module.DefaultRegistry.Get(m.selected)
				if ok {
					queue, err := planInstallQueue(m.selected)
					list := []string{mod.Name + " — " + mod.Description}
					if err != nil {
						list = append(list, "✗ "+err.Error())
					} else if len(queue) == 0 {
						list = append(list, "  (already satisfied — press r to re-run)")
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
		case key.Matches(k, DefaultKeyMap.Reinstall):
			if m.selected != "" {
				has := utils.ModuleScriptExists(m.selected, "install.ps1") ||
					utils.ModuleScriptExists(m.selected, "config.ps1")
				if !has {
					return m, nil
				}
				m.confirm = popup.NewReinstallDialog(m.selected, has)
				m.showConfirm = true
				return m, nil
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
	if len(mod.AdminScripts) > 0 {
		m.detail.OutputModel().AppendLine(
			"  elevated scripts: " + strings.Join(mod.AdminScripts, ", "))
	} else if mod.RequiresAdmin {
		m.detail.OutputModel().AppendLine("  (entire module runs elevated)")
	}
	return m, installer.RunInstallStreaming(m.program, mod, &m.elevate)
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
	deps := make([]detail.DepStatus, 0, len(mod.Dependencies))
	for _, depName := range mod.Dependencies {
		check := ""
		if dep, ok := module.DefaultRegistry.Get(depName); ok {
			check = dep.CheckCommand
		}
		deps = append(deps, detail.DepStatus{
			Name:      depName,
			Method:    "module",
			Installed: utils.ModuleSatisfied(depName, check),
		})
	}
	m.detail.OverviewModel().SetModule(mod.Name, mod.Description, mod.Website, mod.Repo, deps)
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
	m.sidebar.SetFocused(m.focus == FocusSidebar)
	w, h := m.contentDimensions()
	sidebarWidth := int(float64(w) * 0.3)
	detailWidth := w - sidebarWidth - 1
	panelHeight := h - 5
	if panelHeight < 10 {
		panelHeight = 10
	}
	m.sidebar.SetSize(sidebarWidth-2, panelHeight-2)
	m.detail.SetSize(detailWidth-2, panelHeight-2)
	m.sidebar.SetYOffset((m.height-(h+2))/2 + 4)
}

func (m Model) contentDimensions() (int, int) {
	w, h := int(float64(m.width)*.8)-2, int(float64(m.height)*.8)-2
	if w < 60 {
		w = 60
	}
	if h < 20 {
		h = 20
	}
	return w, h
}
func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("Initializing...")
	}
	contentWidth, contentHeight := m.contentDimensions()

	var content string
	switch m.state {
	case StatePrereqCheck:
		content = m.prereq.View()
	case StateModulesSetup:
		content = theme.Title.Render("Modules path setup") + "\n\n" +
			theme.NormalText.Render("Point the TUI at a folder that contains your modules") + "\n" +
			theme.DimText.Render("(each subfolder has module.yaml + optional install.ps1 / config.ps1).") + "\n\n" +
			theme.DimText.Render("Enter a path in the dialog (use ~ for your home directory).")
	case StateDashboard, StateInstalling:
		content = m.viewDashboard(contentWidth, contentHeight)
	}

	content = lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		MaxWidth(contentWidth).
		MaxHeight(contentHeight).
		Render(content)

	window := theme.AppBorder.
		Width(contentWidth).
		Height(contentHeight).
		Render(content)

	finalView := window
	if m.showHelp {
		finalView = m.help.Render(contentWidth+2, contentHeight+2)
	}
	if m.showConfirm {
		finalView = m.confirm.Render(contentWidth+2, contentHeight+2)
	}
	if m.showScript {
		finalView = m.script.Render(contentWidth+2, contentHeight+2)
	}
	if m.showCategory {
		finalView = m.category.Render(contentWidth+2, contentHeight+2)
	}
	if m.showInput {
		finalView = m.input.Render(contentWidth+2, contentHeight+2)
	}

	view := tea.NewView(
		lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, finalView),
	)
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}

func (m Model) viewDashboard(width, height int) string {
	focusLabel := "modules"
	if m.focus == FocusDetail {
		focusLabel = "detail"
	}
	title := theme.Title.Render("⚡ Windows Config") +
		theme.DimText.Render("  · focus: ") +
		theme.Subtitle.Render(focusLabel)

	help := theme.HelpStyle.Render(
		"q: quit • ?: help • j/k: navigate • shift+tab: switch panel • tab: switch tab • i: install • r: re-run • s: search • c: category",
	)

	panelHeight := height - 5
	if panelHeight < 10 {
		panelHeight = 10
	}

	sidebarWidth := int(float64(width) * 0.3)
	detailWidth := width - sidebarWidth - 1

	sidebarContent := m.sidebar.View()
	detailContent := m.detail.View()
	if m.focus != FocusSidebar {
		sidebarContent = lipgloss.NewStyle().Faint(true).Render(sidebarContent)
	}
	if m.focus != FocusDetail {
		detailContent = lipgloss.NewStyle().Faint(true).Render(detailContent)
	}

	sidebarView := framePanel(sidebarContent, m.focus == FocusSidebar, sidebarWidth, panelHeight)
	detailView := framePanel(detailContent, m.focus == FocusDetail, detailWidth, panelHeight)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, " ", detailView)
	body = theme.Clip(body, width, panelHeight)

	return fmt.Sprintf("%s\n\n%s\n\n%s", title, body, help)
}

// framePanel draws a fixed-size panel border. Width/Height are TOTAL outer size.
func framePanel(content string, focused bool, width, height int) string {
	innerW := width - 2
	innerH := height - 2
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}

	borderColor := theme.ColorSurface
	if focused {
		borderColor = theme.ColorCyan
	}

	clipped := theme.Clip(content, innerW, innerH)
	framed := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(innerW).
		Render(clipped)

	return theme.Clip(framed, width, height)
}
