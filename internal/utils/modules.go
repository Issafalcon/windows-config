package utils

import (
	"context"
	"errors"
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
	return exec.CommandContext(ctx, pwsh, "-NoProfile", "-Command", checkCommand).Run() == nil
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
