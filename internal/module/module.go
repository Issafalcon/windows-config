package module

import (
	"strings"
)
// Module describes one installable Windows configuration module.
type Module struct {
	Name          string   `yaml:"name"`
	Icon          string   `yaml:"icon"`
	Description   string   `yaml:"description"`
	Category      string   `yaml:"category"`
	Website       string   `yaml:"website"`
	Repo          string   `yaml:"repo"`
	Dependencies  []string `yaml:"dependencies"`
	RequiresAdmin bool     `yaml:"requires_admin"`
	// AdminScripts lists script basenames (e.g. config.ps1) to run elevated.
	// When non-empty, only those scripts use the elevated helper; install.ps1
	// can stay unelevated (scoop-friendly). When empty, RequiresAdmin elevates both.
	AdminScripts  []string `yaml:"admin_scripts"`
	CheckCommand  string   `yaml:"check_command"`
	EstimatedTime string   `yaml:"estimated_time"`
	EstimatedSize string   `yaml:"estimated_size"`
	ExternalDep   string   `yaml:"external_dep,omitempty"`
}

// ScriptNeedsAdmin reports whether scriptBasename should run via the elevated helper.
func (m Module) ScriptNeedsAdmin(scriptBasename string) bool {
	if len(m.AdminScripts) > 0 {
		for _, s := range m.AdminScripts {
			if strings.EqualFold(s, scriptBasename) {
				return true
			}
		}
		return false
	}
	return m.RequiresAdmin
}
