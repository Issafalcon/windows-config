//go:build !windows

package elevate

import (
	"context"
	"fmt"
)

type Client struct{ JobDir string }

func (*Client) EnsureStarted() error {
	return fmt.Errorf("Windows elevation is unavailable on this platform")
}

func (c *Client) RunScript(ctx context.Context, scriptPath string, args []string, onLine func(string, bool)) error {
	return c.EnsureStarted()
}

func (*Client) Shutdown() {}
