package utils

import (
	"strings"
	"testing"
)

func TestWrapCheckCommandFailsOnNullResult(t *testing.T) {
	wrapped := wrapCheckCommand("Get-Command nosuchthing -ErrorAction SilentlyContinue")
	if !strings.Contains(wrapped, "exit 1") {
		t.Fatalf("wrap missing failure path: %s", wrapped)
	}
	if !strings.Contains(wrapped, "Get-Command nosuchthing") {
		t.Fatalf("wrap dropped original check: %s", wrapped)
	}
}

func TestModuleSatisfiedIgnoresFalseLiteral(t *testing.T) {
	if ModuleSatisfied("vim", "false") {
		t.Fatal("false check_command must not count as satisfied")
	}
	if ModuleSatisfied("omnisharp", "") {
		t.Fatal("empty check_command must not count as satisfied")
	}
}
