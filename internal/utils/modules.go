package utils

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/issafalcon/windows-config-tui/internal/config"
)

const checkTimeout = 10 * time.Second

func GetModulesDir() string {
	return config.ResolveModulesDir()
}

func ModuleScriptPath(name, script string) string {
	return filepath.Join(GetModulesDir(), name, script)
}

func ModuleScriptExists(name, script string) bool {
	info, err := os.Stat(ModuleScriptPath(name, script))
	return err == nil && !info.IsDir()
}

func GetInstalledModules() ([]string, error) {
	path, err := trackingPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return uniqueLines(string(data)), nil
}

func SetModuleInstalled(name string) error {
	modules, err := GetInstalledModules()
	if err != nil {
		return err
	}
	for _, module := range modules {
		if module == name {
			return nil
		}
	}
	return writeInstalledModules(append(modules, name))
}

func SetModuleUninstalled(name string) error {
	modules, err := GetInstalledModules()
	if err != nil {
		return err
	}
	filtered := modules[:0]
	for _, module := range modules {
		if module != name {
			filtered = append(filtered, module)
		}
	}
	return writeInstalledModules(filtered)
}

// ModuleSatisfied reports whether a module should be skipped by the install queue.
//
//  1. Name listed in ~/.windowsConfigModules
//  2. check_command succeeds under pwsh (see wrapCheckCommand)
//
// Literals "true"/"false"/"" skip (2) — tracking-file only.
func ModuleSatisfied(name, checkCommand string) bool {
	installed, err := GetInstalledModules()
	if err == nil && contains(installed, name) {
		return true
	}
	checkCommand = strings.TrimSpace(checkCommand)
	if checkCommand == "" || strings.EqualFold(checkCommand, "true") || strings.EqualFold(checkCommand, "false") {
		return false
	}

	pwsh, err := FindPwsh()
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, pwsh, "-NoProfile", "-Command", wrapCheckCommand(checkCommand))
	return cmd.Run() == nil
}

// wrapCheckCommand makes PowerShell probes fail when they "succeed" with no result.
// `Get-Command foo -ErrorAction SilentlyContinue` exits 0 even when missing; wrapping
// treats a null/false result as failure while still honouring native exit codes.
func wrapCheckCommand(check string) string {
	return fmt.Sprintf(
		`$ErrorActionPreference = 'Continue'; $__result = $(%s); if ($null -ne $LASTEXITCODE -and $LASTEXITCODE -ne 0) { exit $LASTEXITCODE }; if ($null -eq $__result -or $__result -eq $false) { exit 1 }; exit 0`,
		check,
	)
}

func FindPwsh() (string, error) {
	for _, command := range []string{"pwsh", "powershell"} {
		if path, err := exec.LookPath(command); err == nil {
			return path, nil
		}
	}
	return "", exec.ErrNotFound
}

func trackingPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".windowsConfigModules"), nil
}

func writeInstalledModules(modules []string) error {
	path, err := trackingPath()
	if err != nil {
		return err
	}
	data := strings.Join(uniqueLines(strings.Join(modules, "\n")), "\n")
	if data != "" {
		data += "\n"
	}
	return os.WriteFile(path, []byte(data), 0o600)
}

func uniqueLines(data string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if _, exists := seen[line]; !exists {
			seen[line] = struct{}{}
			result = append(result, line)
		}
	}
	return result
}

func contains(modules []string, name string) bool {
	for _, module := range modules {
		if module == name {
			return true
		}
	}
	return false
}
