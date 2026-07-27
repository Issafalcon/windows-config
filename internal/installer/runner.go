package installer

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/issafalcon/windows-config-tui/internal/elevate"
	"github.com/issafalcon/windows-config-tui/internal/utils"
)

type InstallStartMsg struct{ ModuleName string }
type InstallOutputMsg struct {
	ModuleName, Line string
	IsStderr         bool
}
type InstallCompleteMsg struct {
	ModuleName string
	Success    bool
	Error      error
}

func RunInstallStreaming(p *tea.Program, moduleName string, requiresAdmin bool, elevateClient *elevate.Client) tea.Cmd {
	return func() tea.Msg {
		p.Send(InstallStartMsg{moduleName})
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		run := func(name string) error {
			path := utils.ModuleScriptPath(moduleName, name)
			if !utils.ModuleScriptExists(moduleName, name) {
				return nil
			}
			p.Send(InstallOutputMsg{ModuleName: moduleName, Line: "▸ Running " + name})
			line := func(line string, stderr bool) {
				p.Send(InstallOutputMsg{ModuleName: moduleName, Line: line, IsStderr: stderr})
			}
			if requiresAdmin {
				return elevateClient.RunScript(ctx, path, nil, line)
			}
			return utils.RunPwshFileStreaming(ctx, path, nil, line)
		}
		if err := run("install.ps1"); err != nil {
			return InstallCompleteMsg{moduleName, false, fmt.Errorf("install.ps1: %w", err)}
		}
		if err := run("config.ps1"); err != nil {
			return InstallCompleteMsg{moduleName, false, fmt.Errorf("config.ps1: %w", err)}
		}
		if err := utils.SetModuleInstalled(moduleName); err != nil {
			return InstallCompleteMsg{moduleName, false, fmt.Errorf("tracking install: %w", err)}
		}
		return InstallCompleteMsg{ModuleName: moduleName, Success: true}
	}
}
