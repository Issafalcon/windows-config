package module

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadFromDir(modulesDir string) (*Registry, error) {
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, fmt.Errorf("read modules directory: %w", err)
	}

	registry := NewRegistry()
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(modulesDir, entry.Name(), "module.yaml")
		data, err := os.ReadFile(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}

		var m Module
		if err := yaml.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		m.Name = entry.Name()
		if err := registry.Register(m); err != nil {
			return nil, fmt.Errorf("load %s: %w", path, err)
		}
	}
	DefaultRegistry = registry
	return registry, nil
}
