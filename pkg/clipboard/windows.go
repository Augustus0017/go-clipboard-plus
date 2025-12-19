package clipboard

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type windowsClipboard struct{}

func newWindowsClipboard() (*windowsClipboard, error) {
	// Check if PowerShell is available
	if _, err := exec.LookPath("powershell.exe"); err != nil {
		return nil, fmt.Errorf("powershell.exe not found: %v", err)
	}
	return &windowsClipboard{}, nil
}

func (c *windowsClipboard) Available() bool {
	return true
}

func (c *windowsClipboard) Read() (string, error) {
	cmd := exec.Command("powershell.exe", "-Command", "Get-Clipboard")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to read clipboard: %v, stderr: %s", err, stderr.String())
	}

	// PowerShell adds carriage returns, normalize to just newlines
	result := stdout.String()
	result = strings.ReplaceAll(result, "\r\n", "\n")
	return strings.TrimRight(result, "\n"), nil
}

func (c *windowsClipboard) Write(text string) error {
	cmd := exec.Command("powershell.exe", "-Command", "$input | Set-Clipboard")
	cmd.Stdin = strings.NewReader(text)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write clipboard: %v, stderr: %s", err, stderr.String())
	}

	return nil
}
