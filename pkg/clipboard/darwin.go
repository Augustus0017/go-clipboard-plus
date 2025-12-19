package clipboard

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type darwinClipboard struct{}

func newDarwinClipboard() (*darwinClipboard, error) {
	// Check if pbcopy/pbpaste are available (they should be on macOS)
	if _, err := exec.LookPath("pbcopy"); err != nil {
		return nil, fmt.Errorf("pbcopy not found: %v", err)
	}
	if _, err := exec.LookPath("pbpaste"); err != nil {
		return nil, fmt.Errorf("pbpaste not found: %v", err)
	}
	return &darwinClipboard{}, nil
}

func (c *darwinClipboard) Available() bool {
	return true
}

func (c *darwinClipboard) Read() (string, error) {
	cmd := exec.Command("pbpaste")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to read clipboard: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func (c *darwinClipboard) Write(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write clipboard: %v, stderr: %s", err, stderr.String())
	}

	return nil
}
