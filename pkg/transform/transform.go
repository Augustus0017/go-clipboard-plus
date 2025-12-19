package transform

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Transformer is a function that transforms text
type Transformer func(string) (string, error)

// Available transformers
var Transformers = map[string]Transformer{
	"json":        FormatJSON,
	"json-minify": MinifyJSON,
	"trim":        Trim,
	"upper":       Upper,
	"lower":       Lower,
	"title":       Title,
	"base64":      Base64Encode,
	"base64d":     Base64Decode,
	"url":         URLEncode,
	"urld":        URLDecode,
	"reverse":     Reverse,
}

// FormatJSON formats JSON with indentation
func FormatJSON(text string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return "", fmt.Errorf("invalid JSON: %v", err)
	}

	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %v", err)
	}

	return string(formatted), nil
}

// MinifyJSON minifies JSON by removing whitespace
func MinifyJSON(text string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return "", fmt.Errorf("invalid JSON: %v", err)
	}

	minified, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to minify JSON: %v", err)
	}

	return string(minified), nil
}

// Trim removes leading and trailing whitespace
func Trim(text string) (string, error) {
	return strings.TrimSpace(text), nil
}

// Upper converts text to uppercase
func Upper(text string) (string, error) {
	return strings.ToUpper(text), nil
}

// Lower converts text to lowercase
func Lower(text string) (string, error) {
	return strings.ToLower(text), nil
}

// Title converts text to title case
func Title(text string) (string, error) {
	caser := cases.Title(language.English)
	return caser.String(text), nil
}

// Base64Encode encodes text to base64
func Base64Encode(text string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(text)), nil
}

// Base64Decode decodes base64 text
func Base64Decode(text string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}
	return string(decoded), nil
}

// URLEncode encodes text for URLs
func URLEncode(text string) (string, error) {
	return url.QueryEscape(text), nil
}

// URLDecode decodes URL-encoded text
func URLDecode(text string) (string, error) {
	decoded, err := url.QueryUnescape(text)
	if err != nil {
		return "", fmt.Errorf("failed to decode URL: %v", err)
	}
	return decoded, nil
}

// Reverse reverses the text
func Reverse(text string) (string, error) {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

// Apply applies a transformation by name
func Apply(name string, text string) (string, error) {
	transformer, ok := Transformers[name]
	if !ok {
		return "", fmt.Errorf("unknown transformation: %s", name)
	}
	return transformer(text)
}

// List returns all available transformation names
func List() []string {
	names := make([]string, 0, len(Transformers))
	for name := range Transformers {
		names = append(names, name)
	}
	return names
}
