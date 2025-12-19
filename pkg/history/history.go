package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const maxHistorySize = 100

// Entry represents a single clipboard history entry
type Entry struct {
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// History manages clipboard history
type History struct {
	entries  []Entry
	filePath string
}

// New creates a new History instance
func New() (*History, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "go-clipboard-plus")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	h := &History{
		filePath: filepath.Join(configDir, "history.json"),
		entries:  []Entry{},
	}

	// Load existing history
	if err := h.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load history: %v", err)
	}

	return h, nil
}

// Add adds a new entry to the history
func (h *History) Add(content string) error {
	// Don't add empty content
	if content == "" {
		return nil
	}

	// Don't add duplicate of the most recent entry
	if len(h.entries) > 0 && h.entries[0].Content == content {
		return nil
	}

	entry := Entry{
		Content:   content,
		Timestamp: time.Now(),
	}

	// Add to the beginning
	h.entries = append([]Entry{entry}, h.entries...)

	// Limit history size
	if len(h.entries) > maxHistorySize {
		h.entries = h.entries[:maxHistorySize]
	}

	return h.save()
}

// List returns all history entries
func (h *History) List() []Entry {
	return h.entries
}

// Get returns a specific entry by index
func (h *History) Get(index int) (Entry, error) {
	if index < 0 || index >= len(h.entries) {
		return Entry{}, fmt.Errorf("index out of range")
	}
	return h.entries[index], nil
}

// Clear removes all history entries
func (h *History) Clear() error {
	h.entries = []Entry{}
	return h.save()
}

// load loads history from disk
func (h *History) load() error {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &h.entries)
}

// save saves history to disk
func (h *History) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}

	return os.WriteFile(h.filePath, data, 0644)
}
