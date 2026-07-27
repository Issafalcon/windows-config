package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

func RunPwshFileStreaming(ctx context.Context, scriptPath string, args []string, onLine func(string, bool)) error {
	pwsh, err := FindPwsh()
	if err != nil {
		return fmt.Errorf("find PowerShell: %w", err)
	}
	return stream(ctx, exec.CommandContext(ctx, pwsh, append([]string{"-NoProfile", "-File", scriptPath}, args...)...), onLine)
}

func RunCommandStreaming(ctx context.Context, command string, onLine func(string, bool)) error {
	pwsh, err := FindPwsh()
	if err != nil {
		return fmt.Errorf("find PowerShell: %w", err)
	}
	return stream(ctx, exec.CommandContext(ctx, pwsh, "-NoProfile", "-Command", command), onLine)
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
