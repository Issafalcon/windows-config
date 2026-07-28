package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if got := ExpandHome("~"); got != home {
		t.Fatalf("ExpandHome(~) = %q, want %q", got, home)
	}
	if got := ExpandHome("~/repos/x"); got != filepath.Join(home, "repos", "x") {
		t.Fatalf("ExpandHome(~/repos/x) = %q", got)
	}
	if runtime.GOOS == "windows" {
		if got := ExpandHome(`~\repos\x`); got != filepath.Join(home, "repos", "x") {
			t.Fatalf(`ExpandHome(~\repos\x) = %q`, got)
		}
	}
}

func TestNormalizeModulesDirUsesModulesSubdir(t *testing.T) {
	root := t.TempDir()
	modules := filepath.Join(root, "modules")
	if err := os.Mkdir(modules, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := NormalizeModulesDir(root); got != modules {
		t.Fatalf("NormalizeModulesDir(repo root) = %q, want %q", got, modules)
	}
	if got := NormalizeModulesDir(modules); got != modules {
		t.Fatalf("NormalizeModulesDir(modules) = %q, want same", got)
	}
}
