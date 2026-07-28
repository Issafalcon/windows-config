package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	cfg.ModulesDir = NormalizeModulesDir(modulesDir)
	return Save(cfg)
}

func ResolveModulesDir() string {
	if dir := os.Getenv("WINDOWS_CONFIG_MODULES_DIR"); dir != "" {
		return NormalizeModulesDir(dir)
	}
	if cfg, err := Load(); err == nil && cfg.ModulesDir != "" {
		if dir := NormalizeModulesDir(cfg.ModulesDir); dir != "" {
			return dir
		}
	}
	for _, start := range searchRoots() {
		if dir := findModulesDir(start); dir != "" {
			return dir
		}
	}
	return ""
}

// ExpandHome expands a leading ~ / ~/ / ~\ to the user home directory.
func ExpandHome(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	path = os.ExpandEnv(path)
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
		return path
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

// NormalizeModulesDir expands home shortcuts and, if the path is a repo root
// that contains modules/, returns that modules/ directory.
func NormalizeModulesDir(path string) string {
	path = ExpandHome(path)
	if path == "" {
		return ""
	}
	path = filepath.Clean(path)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	if base := filepath.Base(path); !strings.EqualFold(base, "modules") {
		candidate := filepath.Join(path, "modules")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	return path
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
