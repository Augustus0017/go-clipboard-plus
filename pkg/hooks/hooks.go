package hooks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Hook represents a script hook
type Hook struct {
	Name   string
	Script string
}

// Manager manages script hooks
type Manager struct {
	hooksDir string
}

// New creates a new hook manager
func New() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	hooksDir := filepath.Join(homeDir, ".config", "go-clipboard-plus", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create hooks directory: %v", err)
	}

	return &Manager{
		hooksDir: hooksDir,
	}, nil
}

// List returns all available hooks
func (m *Manager) List() ([]Hook, error) {
	entries, err := os.ReadDir(m.hooksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read hooks directory: %v", err)
	}

	var hooks []Hook
	for _, entry := range entries {
		if entry.IsDir() || !isExecutable(entry) {
			continue
		}

		hooks = append(hooks, Hook{
			Name:   entry.Name(),
			Script: filepath.Join(m.hooksDir, entry.Name()),
		})
	}

	return hooks, nil
}

// Execute runs a hook with the given input
func (m *Manager) Execute(name string, input string) (string, error) {
	hookPath := filepath.Join(m.hooksDir, name)

	// Check if hook exists and is executable
	info, err := os.Stat(hookPath)
	if err != nil {
		return "", fmt.Errorf("hook not found: %s", name)
	}

	if info.IsDir() {
		return "", fmt.Errorf("hook is a directory: %s", name)
	}

	if !isExecutableFile(info) {
		return "", fmt.Errorf("hook is not executable: %s", name)
	}

	// Execute the hook
	cmd := exec.Command(hookPath)
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("hook execution failed: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// GetHooksDir returns the hooks directory path
func (m *Manager) GetHooksDir() string {
	return m.hooksDir
}

// isExecutable checks if a directory entry is executable
func isExecutable(entry os.DirEntry) bool {
	info, err := entry.Info()
	if err != nil {
		return false
	}
	return isExecutableFile(info)
}

// isExecutableFile checks if a file is executable
func isExecutableFile(info os.FileInfo) bool {
	return info.Mode()&0111 != 0
}
