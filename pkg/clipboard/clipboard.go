package clipboard

import (
	"fmt"
	"runtime"
)

// Clipboard provides cross-platform clipboard access
type Clipboard interface {
	Read() (string, error)
	Write(text string) error
	Available() bool
}

// New returns a platform-specific clipboard implementation
func New() (Clipboard, error) {
	switch runtime.GOOS {
	case "linux":
		return newLinuxClipboard()
	case "darwin":
		return newDarwinClipboard()
	case "windows":
		return newWindowsClipboard()
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
