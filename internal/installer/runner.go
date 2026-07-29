package installer

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/issafalcon/windows-config-tui/internal/elevate"
	"github.com/issafalcon/windows-config-tui/internal/module"
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

// RunInstallStreaming runs install.ps1 then config.ps1, elevating only the
// scripts that mod.ScriptNeedsAdmin reports (symlinks / machine changes).
func RunInstallStreaming(p *tea.Program, mod module.Module, elevateClient *elevate.Client) tea.Cmd {
	return runModuleScripts(p, mod, elevateClient, []string{"install.ps1", "config.ps1"}, true)
}

// RunUninstallStreaming runs uninstall.ps1 (if present) then clears install tracking.
func RunUninstallStreaming(p *tea.Program, mod module.Module, elevateClient *elevate.Client) tea.Cmd {
	return runModuleScripts(p, mod, elevateClient, []string{"uninstall.ps1"}, false)
}

func runModuleScripts(
	p *tea.Program,
	mod module.Module,
	elevateClient *elevate.Client,
	scripts []string,
	markInstalled bool,
) tea.Cmd {
	return func() tea.Msg {
		moduleName := mod.Name
		p.Send(InstallStartMsg{moduleName})
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
		defer cancel()

		run := func(script string) error {
			path := utils.ModuleScriptPath(moduleName, script)
			if !utils.ModuleScriptExists(moduleName, script) {
				return nil
			}
			elevated := mod.ScriptNeedsAdmin(script)
			label := "▸ Running " + script
			if elevated {
				label += " (elevated — approve UAC if prompted)"
			}
			p.Send(InstallOutputMsg{ModuleName: moduleName, Line: label})
			line := func(line string, stderr bool) {
				p.Send(InstallOutputMsg{ModuleName: moduleName, Line: line, IsStderr: stderr})
			}
			if elevated {
				if elevateClient == nil {
					return fmt.Errorf("elevated helper unavailable")
				}
				return elevateClient.RunScript(ctx, path, nil, line)
			}
			return utils.RunPwshFileStreaming(ctx, path, nil, line)
		}

		for _, script := range scripts {
			if err := run(script); err != nil {
				return InstallCompleteMsg{moduleName, false, fmt.Errorf("%s: %w", script, err)}
			}
		}
		if markInstalled {
			if err := utils.SetModuleInstalled(moduleName); err != nil {
				return InstallCompleteMsg{moduleName, false, fmt.Errorf("tracking install: %w", err)}
			}
		} else {
			if err := utils.SetModuleUninstalled(moduleName); err != nil {
				return InstallCompleteMsg{moduleName, false, fmt.Errorf("tracking uninstall: %w", err)}
			}
		}
		return InstallCompleteMsg{ModuleName: moduleName, Success: true}
	}
}
