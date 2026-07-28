//go:build windows

package elevate

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	JobDir        string
	helperProcess *os.Process
}
type request struct {
	ID         string   `json:"id"`
	ScriptPath string   `json:"scriptPath"`
	Args       []string `json:"args"`
}

func (c *Client) EnsureStarted() error {
	if c.JobDir != "" {
		return nil
	}
	dir, err := os.MkdirTemp("", "windows-config-elevated-")
	if err != nil {
		return err
	}
	helper, err := helperPath()
	if err != nil {
		_ = os.RemoveAll(dir)
		return err
	}
	c.JobDir = dir
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }
	pwsh, err := exec.LookPath("pwsh")
	if err != nil {
		pwsh, err = exec.LookPath("powershell")
	}
	if err != nil {
		c.JobDir = ""
		_ = os.RemoveAll(dir)
		return fmt.Errorf("pwsh/powershell not found: %w", err)
	}
	// Single ArgumentList string is more reliable under -Verb RunAs than a
	// PowerShell array. Bypass + Unblock covers Mark-of-the-Web on zip extracts.
	argList := fmt.Sprintf(`-NoProfile -ExecutionPolicy Bypass -File "%s" -JobDir "%s"`, helper, dir)
	command := fmt.Sprintf(
		"Unblock-File -LiteralPath %s -ErrorAction SilentlyContinue; Start-Process -FilePath %s -Verb RunAs -ArgumentList %s",
		quote(helper), quote(pwsh), quote(argList),
	)
	if err := exec.Command(pwsh, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", command).Run(); err != nil {
		c.JobDir = ""
		_ = os.RemoveAll(dir)
		return fmt.Errorf("start elevated helper: %w", err)
	}
	// Wait until the elevated helper writes its ready marker (UAC may take a moment).
	deadline := time.Now().Add(2 * time.Minute)
	ready := filepath.Join(dir, "ready")
	for time.Now().Before(deadline) {
		if _, err := os.Stat(ready); err == nil {
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}
	c.JobDir = ""
	_ = os.RemoveAll(dir)
	return fmt.Errorf("elevated helper did not start (UAC cancelled or timed out); see %%LOCALAPPDATA%%\\windows-config-tui\\elevated-helper.log")
}

func (c *Client) RunScript(ctx context.Context, scriptPath string, args []string, onLine func(string, bool)) error {
	if err := c.EnsureStarted(); err != nil {
		return err
	}
	id := newID()
	data, err := json.Marshal(request{ID: id, ScriptPath: scriptPath, Args: args})
	if err != nil {
		return err
	}
	req := filepath.Join(c.JobDir, id+".req.json")
	if err := os.WriteFile(req, data, 0o600); err != nil {
		return err
	}
	logPath, donePath := filepath.Join(c.JobDir, id+".log"), filepath.Join(c.JobDir, id+".done")
	offset := int64(0)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if f, err := os.Open(logPath); err == nil {
				_, _ = f.Seek(offset, 0)
				buf := make([]byte, 4096)
				for {
					n, e := f.Read(buf)
					if n > 0 {
						offset += int64(n)
						for _, line := range strings.Split(string(buf[:n]), "\n") {
							if line != "" {
								onLine(strings.TrimSuffix(line, "\r"), false)
							}
						}
					}
					if e != nil {
						break
					}
				}
				_ = f.Close()
			}
			if data, err := os.ReadFile(donePath); err == nil {
				code, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
				if parseErr != nil {
					return parseErr
				}
				if code != 0 {
					return fmt.Errorf("elevated script exited with code %d", code)
				}
				return nil
			}
		}
	}
}

func (c *Client) Shutdown() {
	if c.JobDir == "" {
		return
	}
	jobDir := c.JobDir
	_ = os.WriteFile(filepath.Join(jobDir, "shutdown"), nil, 0o600)
	time.Sleep(300 * time.Millisecond)
	_ = os.RemoveAll(jobDir)
	c.JobDir = ""
}
func helperPath() (string, error) {
	if p := os.Getenv("WINDOWS_CONFIG_HELPER"); p != "" {
		return p, nil
	}
	starts := []string{}
	if cwd, e := os.Getwd(); e == nil {
		starts = append(starts, cwd)
	}
	if exe, e := os.Executable(); e == nil {
		starts = append(starts, filepath.Dir(exe))
	}
	for _, start := range starts {
		for dir := start; ; dir = filepath.Dir(dir) {
			p := filepath.Join(dir, "tools", "ElevatedHelper.ps1")
			if _, e := os.Stat(p); e == nil {
				return p, nil
			}
			if filepath.Dir(dir) == dir {
				break
			}
		}
	}
	return "", fmt.Errorf("tools/ElevatedHelper.ps1 not found; set WINDOWS_CONFIG_HELPER")
}
func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(b)
}
