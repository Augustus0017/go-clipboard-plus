package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	h, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, ".config", "go-clipboard-plus", "history.json")
	if h.filePath != expectedPath {
		t.Errorf("History file path = %v, want %v", h.filePath, expectedPath)
	}
}

func TestAdd(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	h, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test adding entries
	err = h.Add("test1")
	if err != nil {
		t.Errorf("Add() error = %v", err)
	}

	err = h.Add("test2")
	if err != nil {
		t.Errorf("Add() error = %v", err)
	}

	entries := h.List()
	if len(entries) != 2 {
		t.Errorf("List() length = %v, want 2", len(entries))
	}

	// Check order (most recent first)
	if entries[0].Content != "test2" {
		t.Errorf("First entry = %v, want test2", entries[0].Content)
	}
}

func TestAddEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	h, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Adding empty content should be ignored
	err = h.Add("")
	if err != nil {
		t.Errorf("Add() error = %v", err)
	}

	entries := h.List()
	if len(entries) != 0 {
		t.Errorf("List() length = %v, want 0", len(entries))
	}
}

func TestAddDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	h, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Add same content twice
	h.Add("test")
	h.Add("test")

	entries := h.List()
	// Should only have one entry (duplicates are not added)
	if len(entries) != 1 {
		t.Errorf("List() length = %v, want 1 (duplicates should not be added)", len(entries))
	}
}

func TestGet(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	h, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	h.Add("test1")
	h.Add("test2")

	// Test valid index
	entry, err := h.Get(0)
	if err != nil {
		t.Errorf("Get(0) error = %v", err)
	}
	if entry.Content != "test2" {
		t.Errorf("Get(0).Content = %v, want test2", entry.Content)
	}

	// Test invalid index
	_, err = h.Get(10)
	if err == nil {
		t.Error("Get(10) should return error for out of range index")
	}
}

func TestClear(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	h, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	h.Add("test1")
	h.Add("test2")

	err = h.Clear()
	if err != nil {
		t.Errorf("Clear() error = %v", err)
	}

	entries := h.List()
	if len(entries) != 0 {
		t.Errorf("List() length after clear = %v, want 0", len(entries))
	}
}

func TestPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	// Create history and add entry
	h1, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	h1.Add("test persistent")

	// Create new history instance (should load from disk)
	h2, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	entries := h2.List()
	if len(entries) != 1 {
		t.Errorf("Loaded history length = %v, want 1", len(entries))
	}

	if entries[0].Content != "test persistent" {
		t.Errorf("Loaded entry content = %v, want 'test persistent'", entries[0].Content)
	}
}

func TestTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	h, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	before := time.Now()
	h.Add("test")
	after := time.Now()

	entries := h.List()
	if len(entries) != 1 {
		t.Fatalf("List() length = %v, want 1", len(entries))
	}

	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("Timestamp %v is not between %v and %v", ts, before, after)
	}
}
