package module

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
	CheckCommand  string   `yaml:"check_command"`
	EstimatedTime string   `yaml:"estimated_time"`
	EstimatedSize string   `yaml:"estimated_size"`
	ExternalDep   string   `yaml:"external_dep,omitempty"`
}
