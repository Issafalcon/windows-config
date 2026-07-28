package module

import "testing"

func TestScriptNeedsAdmin(t *testing.T) {
	m := Module{RequiresAdmin: true}
	if !m.ScriptNeedsAdmin("config.ps1") || !m.ScriptNeedsAdmin("install.ps1") {
		t.Fatal("RequiresAdmin should elevate both scripts")
	}

	m = Module{AdminScripts: []string{"config.ps1"}}
	if !m.ScriptNeedsAdmin("config.ps1") {
		t.Fatal("config.ps1 should elevate")
	}
	if m.ScriptNeedsAdmin("install.ps1") {
		t.Fatal("install.ps1 should stay unelevated when only config is listed")
	}
}
