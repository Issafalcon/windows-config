package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

func quotePS(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func RunPwshFileStreaming(ctx context.Context, scriptPath string, args []string, onLine func(string, bool)) error {
	pwsh, err := FindPwsh()
	if err != nil {
		return fmt.Errorf("find PowerShell: %w", err)
	}
	quotedArgs := make([]string, 0, len(args))
	for _, a := range args {
		quotedArgs = append(quotedArgs, quotePS(a))
	}
	// Run via -Command so *>&1 merges Write-Host / Warning / Error into the
	// success stream that Go's StdoutPipe can read (config.ps1 is mostly host output).
	cmdStr := fmt.Sprintf(
		"$InformationPreference='Continue'; $ProgressPreference='SilentlyContinue'; "+
			"& %s @(%s) *>&1 | ForEach-Object { $_.ToString() }",
		quotePS(scriptPath), strings.Join(quotedArgs, ","),
	)
	return stream(ctx, exec.CommandContext(ctx, pwsh,
		"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", cmdStr), onLine)
}

func RunCommandStreaming(ctx context.Context, command string, onLine func(string, bool)) error {
	pwsh, err := FindPwsh()
	if err != nil {
		return fmt.Errorf("find PowerShell: %w", err)
	}
	wrapped := "$InformationPreference='Continue'; $ProgressPreference='SilentlyContinue'; " +
		"& { " + command + " } *>&1 | ForEach-Object { $_.ToString() }"
	return stream(ctx, exec.CommandContext(ctx, pwsh,
		"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", wrapped), onLine)
}

func stream(ctx context.Context, cmd *exec.Cmd, onLine func(string, bool)) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	var wg sync.WaitGroup
	var scanErr error
	var mu sync.Mutex
	read := func(r io.Reader, isStderr bool) {
		defer wg.Done()
		s := bufio.NewScanner(r)
		s.Buffer(make([]byte, 4096), 1024*1024)
		for s.Scan() {
			onLine(s.Text(), isStderr)
		}
		if err := s.Err(); err != nil {
			mu.Lock()
			scanErr = err
			mu.Unlock()
		}
	}
	wg.Add(2)
	go read(stdout, false)
	go read(stderr, true)
	wg.Wait()
	if err := cmd.Wait(); err != nil {
		return err
	}
	return scanErr
}
