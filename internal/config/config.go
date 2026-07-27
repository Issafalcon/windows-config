package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ModulesDir string `yaml:"modules_dir"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "windows-config-tui", "config.yaml"), nil
}

func Load() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func Save(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func SetModulesDir(modulesDir string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	cfg.ModulesDir = modulesDir
	return Save(cfg)
}

func ResolveModulesDir() string {
	if dir := os.Getenv("WINDOWS_CONFIG_MODULES_DIR"); dir != "" {
		return dir
	}
	if cfg, err := Load(); err == nil && cfg.ModulesDir != "" {
		return cfg.ModulesDir
	}
	for _, start := range searchRoots() {
		if dir := findModulesDir(start); dir != "" {
			return dir
		}
	}
	return ""
}

func searchRoots() []string {
	roots := make([]string, 0, 2)
	if cwd, err := os.Getwd(); err == nil {
		roots = append(roots, cwd)
	}
	if executable, err := os.Executable(); err == nil {
		roots = append(roots, filepath.Dir(executable))
	}
	return roots
}

func findModulesDir(start string) string {
	for dir := filepath.Clean(start); ; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, "modules")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
	}
}
