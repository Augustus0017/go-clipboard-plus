package clipboard

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type linuxClipboard struct {
	copyCmd  string
	pasteCmd string
}

func newLinuxClipboard() (*linuxClipboard, error) {
	c := &linuxClipboard{}

	// Try to find xclip first, then xsel
	if _, err := exec.LookPath("xclip"); err == nil {
		c.copyCmd = "xclip"
		c.pasteCmd = "xclip"
		return c, nil
	}

	if _, err := exec.LookPath("xsel"); err == nil {
		c.copyCmd = "xsel"
		c.pasteCmd = "xsel"
		return c, nil
	}

	return nil, fmt.Errorf("no clipboard utility found (install xclip or xsel)")
}

func (c *linuxClipboard) Available() bool {
	return c.copyCmd != "" && c.pasteCmd != ""
}

func (c *linuxClipboard) Read() (string, error) {
	var cmd *exec.Cmd

	if c.pasteCmd == "xclip" {
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
	} else {
		cmd = exec.Command("xsel", "--clipboard", "--output")
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to read clipboard: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func (c *linuxClipboard) Write(text string) error {
	var cmd *exec.Cmd

	if c.copyCmd == "xclip" {
		cmd = exec.Command("xclip", "-selection", "clipboard")
	} else {
		cmd = exec.Command("xsel", "--clipboard", "--input")
	}

	cmd.Stdin = strings.NewReader(text)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write clipboard: %v, stderr: %s", err, stderr.String())
	}

	return nil
}
