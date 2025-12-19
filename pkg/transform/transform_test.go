package transform

import (
	"strings"
	"testing"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid json",
			input:   `{"name":"test","age":30}`,
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.Contains(result, "\n") {
				t.Errorf("FormatJSON() should contain newlines for formatted output")
			}
		})
	}
}

func TestMinifyJSON(t *testing.T) {
	input := `{
		"name": "test",
		"age": 30
	}`
	result, err := MinifyJSON(input)
	if err != nil {
		t.Errorf("MinifyJSON() error = %v", err)
	}
	if strings.Contains(result, "\n") || strings.Contains(result, "  ") {
		t.Errorf("MinifyJSON() should not contain newlines or extra spaces")
	}
}

func TestTrim(t *testing.T) {
	input := "  hello world  \n"
	expected := "hello world"
	result, err := Trim(input)
	if err != nil {
		t.Errorf("Trim() error = %v", err)
	}
	if result != expected {
		t.Errorf("Trim() = %v, want %v", result, expected)
	}
}

func TestUpper(t *testing.T) {
	input := "hello world"
	expected := "HELLO WORLD"
	result, err := Upper(input)
	if err != nil {
		t.Errorf("Upper() error = %v", err)
	}
	if result != expected {
		t.Errorf("Upper() = %v, want %v", result, expected)
	}
}

func TestLower(t *testing.T) {
	input := "HELLO WORLD"
	expected := "hello world"
	result, err := Lower(input)
	if err != nil {
		t.Errorf("Lower() error = %v", err)
	}
	if result != expected {
		t.Errorf("Lower() = %v, want %v", result, expected)
	}
}

func TestBase64(t *testing.T) {
	input := "hello world"

	// Test encode
	encoded, err := Base64Encode(input)
	if err != nil {
		t.Errorf("Base64Encode() error = %v", err)
	}

	// Test decode
	decoded, err := Base64Decode(encoded)
	if err != nil {
		t.Errorf("Base64Decode() error = %v", err)
	}

	if decoded != input {
		t.Errorf("Base64 round trip failed: got %v, want %v", decoded, input)
	}
}

func TestURLEncodeDecode(t *testing.T) {
	input := "hello world & test=value"

	// Test encode
	encoded, err := URLEncode(input)
	if err != nil {
		t.Errorf("URLEncode() error = %v", err)
	}

	// Test decode
	decoded, err := URLDecode(encoded)
	if err != nil {
		t.Errorf("URLDecode() error = %v", err)
	}

	if decoded != input {
		t.Errorf("URL encode/decode round trip failed: got %v, want %v", decoded, input)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"abc", "cba"},
		{"", ""},
		{"a", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Reverse(tt.input)
			if err != nil {
				t.Errorf("Reverse() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Reverse(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name           string
		transformation string
		input          string
		wantErr        bool
	}{
		{"trim", "trim", "  test  ", false},
		{"upper", "upper", "test", false},
		{"unknown", "unknown-transform", "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Apply(tt.transformation, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestList(t *testing.T) {
	transformations := List()
	if len(transformations) == 0 {
		t.Error("List() should return at least one transformation")
	}

	// Check that common transformations are present
	found := false
	for _, name := range transformations {
		if name == "json" || name == "trim" || name == "base64" {
			found = true
			break
		}
	}
	if !found {
		t.Error("List() should contain common transformations like json, trim, or base64")
	}
}
