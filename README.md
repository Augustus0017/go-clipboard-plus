# go-clipboard-plus

An advanced clipboard manager and transformer for developers with cross-platform support, history tracking, text transformations, and custom scripting hooks.

## Features

- 🖥️ **Cross-platform**: Works on Linux, macOS, and Windows
- 📋 **Clipboard Operations**: Read from and write to system clipboard
- 📚 **History Tracking**: Automatically tracks clipboard history (up to 100 entries)
- 🔄 **Text Transformations**: Built-in transformations for common tasks
  - Format/minify JSON
  - Trim whitespace
  - Case conversions (upper, lower, title)
  - Base64 encode/decode
  - URL encode/decode
  - Reverse text
- 🔧 **Custom Hooks**: Execute custom scripts on clipboard content
- ⌨️ **Clean CLI**: Simple and intuitive command-line interface

## Installation

### Prerequisites

**Linux**: Install `xclip` or `xsel`
```bash
# Ubuntu/Debian
sudo apt-get install xclip

# Fedora
sudo dnf install xclip

# Arch Linux
sudo pacman -S xclip
```

**macOS**: No additional dependencies (uses built-in `pbcopy`/`pbpaste`)

**Windows**: No additional dependencies (uses PowerShell)

### Build from source

```bash
git clone https://github.com/BaseMax/go-clipboard-plus.git
cd go-clipboard-plus
go build -o clipctl ./cmd/clipctl
```

Then move the `clipctl` binary to a directory in your PATH:
```bash
sudo mv clipctl /usr/local/bin/
```

## Usage

### Basic Commands

#### Copy to clipboard
```bash
# Copy text directly
clipctl copy "Hello World"

# Copy from stdin
echo "Hello World" | clipctl copy
cat file.txt | clipctl copy
```

#### Paste from clipboard
```bash
# Paste to stdout
clipctl paste

# Pipe to other commands
clipctl paste | grep "pattern"
```

#### View clipboard history
```bash
# Show last 10 entries (default)
clipctl history

# Show last N entries
clipctl history -n 5

# Show all entries
clipctl history -a

# Copy a specific history entry to clipboard
clipctl history -g 2
```

#### Clear history
```bash
clipctl clear
```

### Transformations

Apply transformations to clipboard content:

```bash
# Format JSON
clipctl transform json

# Minify JSON
clipctl transform json-minify

# Trim whitespace
clipctl transform trim

# Convert to uppercase
clipctl transform upper

# Convert to lowercase
clipctl transform lower

# Title case
clipctl transform title

# Base64 encode
clipctl transform base64

# Base64 decode
clipctl transform base64d

# URL encode
clipctl transform url

# URL decode
clipctl transform urld

# Reverse text
clipctl transform reverse
```

### Custom Hooks

Create custom script hooks to process clipboard content:

```bash
# Show hooks directory
clipctl hook --dir

# List available hooks
clipctl hook --list

# Execute a hook
clipctl hook my-script
```

#### Creating a Hook

1. Create an executable script in the hooks directory:
```bash
mkdir -p ~/.config/go-clipboard-plus/hooks
```

2. Create a script (e.g., `uppercase.sh`):
```bash
#!/bin/bash
tr '[:lower:]' '[:upper:]'
```

3. Make it executable:
```bash
chmod +x ~/.config/go-clipboard-plus/hooks/uppercase.sh
```

4. Use it:
```bash
clipctl hook uppercase.sh
```

Hooks receive clipboard content via stdin and should output the processed content to stdout.

### Advanced Examples

```bash
# Copy file contents and format as JSON
cat data.json | clipctl copy
clipctl transform json

# Chain operations with paste
clipctl paste | jq . | clipctl copy

# Get history entry and encode to base64
clipctl history -g 0
clipctl transform base64

# Process clipboard with custom hook
clipctl paste | clipctl hook my-processor
```

## Architecture

```
go-clipboard-plus/
├── cmd/
│   └── clipctl/          # CLI application
│       └── main.go
├── pkg/
│   ├── clipboard/        # Cross-platform clipboard access
│   │   ├── clipboard.go  # Interface
│   │   ├── linux.go      # Linux implementation
│   │   ├── darwin.go     # macOS implementation
│   │   └── windows.go    # Windows implementation
│   ├── history/          # History tracking
│   │   └── history.go
│   ├── transform/        # Text transformations
│   │   └── transform.go
│   └── hooks/            # Custom script hooks
│       └── hooks.go
└── README.md
```

## Configuration

Configuration files are stored in:
- **Linux/macOS**: `~/.config/go-clipboard-plus/`
- **Windows**: `%USERPROFILE%\.config\go-clipboard-plus\`

### Files
- `history.json`: Clipboard history storage
- `hooks/`: Directory for custom hook scripts

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./pkg/transform -v
go test ./pkg/history -v

# Run tests with coverage
go test -cover ./...
```

### Building

```bash
# Build for current platform
go build -o clipctl ./cmd/clipctl

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o clipctl-linux ./cmd/clipctl
GOOS=darwin GOARCH=amd64 go build -o clipctl-mac ./cmd/clipctl
GOOS=windows GOARCH=amd64 go build -o clipctl.exe ./cmd/clipctl
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.

## Author

**Max Base**

- GitHub: [@BaseMax](https://github.com/BaseMax)

## Acknowledgments

- Uses OS-specific clipboard utilities for reliable clipboard access
- Inspired by various clipboard manager tools
