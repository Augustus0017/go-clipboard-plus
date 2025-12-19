package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/BaseMax/go-clipboard-plus/pkg/clipboard"
	"github.com/BaseMax/go-clipboard-plus/pkg/history"
	"github.com/BaseMax/go-clipboard-plus/pkg/hooks"
	"github.com/BaseMax/go-clipboard-plus/pkg/transform"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "copy", "c":
		handleCopy()
	case "paste", "p":
		handlePaste()
	case "history", "h":
		handleHistory()
	case "transform", "t":
		handleTransform()
	case "hook":
		handleHook()
	case "clear":
		handleClear()
	case "version", "v":
		fmt.Printf("go-clipboard-plus version %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	usage := `go-clipboard-plus - Advanced clipboard manager and transformer

Usage:
  clipctl <command> [arguments]

Commands:
  copy, c [text]          Copy text to clipboard (reads from stdin if no text provided)
  paste, p                Paste text from clipboard
  history, h [options]    Show clipboard history
    -n <number>           Show only last N entries (default: 10)
    -a, --all             Show all entries
    -g, --get <index>     Get specific entry by index and copy to clipboard
  transform, t <name>     Transform clipboard content
    Available transformations:
      json        - Format JSON with indentation
      json-minify - Minify JSON
      trim        - Trim whitespace
      upper       - Convert to uppercase
      lower       - Convert to lowercase
      title       - Convert to title case
      base64      - Encode to base64
      base64d     - Decode from base64
      url         - URL encode
      urld        - URL decode
      reverse     - Reverse text
  hook <name>             Execute a custom hook script on clipboard content
    --list                List available hooks
    --dir                 Show hooks directory
  clear                   Clear clipboard history
  version, v              Show version information
  help, -h, --help        Show this help message

Examples:
  clipctl copy "Hello World"           # Copy text to clipboard
  echo "data" | clipctl copy           # Copy from stdin
  clipctl paste                        # Paste from clipboard
  clipctl history -n 5                 # Show last 5 history entries
  clipctl history -g 2                 # Copy history entry 2 to clipboard
  clipctl paste | clipctl transform json  # Format clipboard as JSON
  clipctl transform base64             # Encode clipboard to base64
  clipctl hook my-script               # Run custom hook on clipboard
`
	fmt.Print(usage)
}

func handleCopy() {
	clip, err := clipboard.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var text string
	if len(os.Args) > 2 {
		// Copy from arguments
		text = strings.Join(os.Args[2:], " ")
	} else {
		// Copy from stdin
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		text = string(data)
	}

	if err := clip.Write(text); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to clipboard: %v\n", err)
		os.Exit(1)
	}

	// Add to history
	hist, err := history.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save to history: %v\n", err)
	} else {
		if err := hist.Add(text); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to save to history: %v\n", err)
		}
	}

	fmt.Fprintf(os.Stderr, "✓ Copied %d bytes to clipboard\n", len(text))
}

func handlePaste() {
	clip, err := clipboard.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	text, err := clip.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from clipboard: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(text)
}

func handleHistory() {
	hist, err := history.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	entries := hist.List()

	// Parse options
	limit := 10
	showAll := false
	getIndex := -1

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-a", "--all":
			showAll = true
		case "-n":
			if i+1 < len(os.Args) {
				i++
				fmt.Sscanf(os.Args[i], "%d", &limit)
			}
		case "-g", "--get":
			if i+1 < len(os.Args) {
				i++
				fmt.Sscanf(os.Args[i], "%d", &getIndex)
			}
		}
	}

	// Handle get operation
	if getIndex >= 0 {
		entry, err := hist.Get(getIndex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		clip, err := clipboard.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := clip.Write(entry.Content); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to clipboard: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "✓ Copied entry %d to clipboard\n", getIndex)
		return
	}

	// Display history
	if len(entries) == 0 {
		fmt.Println("No history entries")
		return
	}

	if !showAll && len(entries) > limit {
		entries = entries[:limit]
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "INDEX\tTIME\tCONTENT")
	fmt.Fprintln(w, "-----\t----\t-------")

	for i, entry := range entries {
		content := entry.Content
		if len(content) > 60 {
			content = content[:60] + "..."
		}
		// Replace newlines with spaces for display
		content = strings.ReplaceAll(content, "\n", "\\n")

		timeStr := formatTime(entry.Timestamp)
		fmt.Fprintf(w, "%d\t%s\t%s\n", i, timeStr, content)
	}

	w.Flush()

	if !showAll && len(hist.List()) > limit {
		fmt.Fprintf(os.Stderr, "\nShowing %d of %d entries. Use -a to show all.\n", limit, len(hist.List()))
	}
}

func handleTransform() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: transformation name required")
		fmt.Fprintf(os.Stderr, "\nAvailable transformations: %s\n", strings.Join(transform.List(), ", "))
		os.Exit(1)
	}

	transformName := os.Args[2]

	clip, err := clipboard.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	text, err := clip.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from clipboard: %v\n", err)
		os.Exit(1)
	}

	transformed, err := transform.Apply(transformName, text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error transforming: %v\n", err)
		os.Exit(1)
	}

	if err := clip.Write(transformed); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to clipboard: %v\n", err)
		os.Exit(1)
	}

	// Add to history
	hist, err := history.New()
	if err == nil {
		hist.Add(transformed)
	}

	fmt.Fprintf(os.Stderr, "✓ Applied transformation: %s\n", transformName)
	fmt.Print(transformed)
}

func handleHook() {
	hm, err := hooks.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: hook name or option required")
		fmt.Fprintln(os.Stderr, "Use --list to see available hooks or --dir to see hooks directory")
		os.Exit(1)
	}

	arg := os.Args[2]

	switch arg {
	case "--list":
		hooks, err := hm.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(hooks) == 0 {
			fmt.Println("No hooks available")
			fmt.Printf("Add executable scripts to: %s\n", hm.GetHooksDir())
			return
		}

		fmt.Println("Available hooks:")
		for _, hook := range hooks {
			fmt.Printf("  - %s\n", hook.Name)
		}

	case "--dir":
		fmt.Println(hm.GetHooksDir())

	default:
		// Execute hook
		hookName := arg

		clip, err := clipboard.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		text, err := clip.Read()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from clipboard: %v\n", err)
			os.Exit(1)
		}

		output, err := hm.Execute(hookName, text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing hook: %v\n", err)
			os.Exit(1)
		}

		if err := clip.Write(output); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to clipboard: %v\n", err)
			os.Exit(1)
		}

		// Add to history
		hist, err := history.New()
		if err == nil {
			hist.Add(output)
		}

		fmt.Fprintf(os.Stderr, "✓ Executed hook: %s\n", hookName)
		fmt.Print(output)
	}
}

func handleClear() {
	hist, err := history.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := hist.Clear(); err != nil {
		fmt.Fprintf(os.Stderr, "Error clearing history: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ History cleared")
}

func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 02, 2006")
	}
}
