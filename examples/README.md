# go-clipboard-plus Examples

This directory contains example scripts and usage patterns for go-clipboard-plus.

## Hook Scripts

The `hooks/` directory contains example hook scripts that can be used with the `clipctl hook` command.

### Installation

Copy the hook scripts to your hooks directory:

```bash
# Create hooks directory
mkdir -p ~/.config/go-clipboard-plus/hooks

# Copy example hooks
cp examples/hooks/*.sh ~/.config/go-clipboard-plus/hooks/
cp examples/hooks/*.py ~/.config/go-clipboard-plus/hooks/

# Make them executable
chmod +x ~/.config/go-clipboard-plus/hooks/*
```

### Available Example Hooks

#### uppercase.sh
Converts all text to uppercase.

```bash
echo "hello world" | clipctl copy
clipctl hook uppercase.sh
clipctl paste  # Output: HELLO WORLD
```

#### line-numbers.sh
Adds line numbers to each line of text.

```bash
cat myfile.txt | clipctl copy
clipctl hook line-numbers.sh
```

#### count-stats.py
Displays statistics about the text (lines, words, characters).

```bash
clipctl copy "Hello world, this is a test"
clipctl hook count-stats.py
```

#### sort-lines.sh
Sorts lines alphabetically.

```bash
printf "zebra\napple\nbanana" | clipctl copy
clipctl hook sort-lines.sh
clipctl paste  # Output: apple, banana, zebra
```

## Usage Patterns

### Pattern 1: Quick Text Transformation

```bash
# Copy JSON and format it
echo '{"name":"test","value":123}' | clipctl copy
clipctl transform json

# View result
clipctl paste
```

### Pattern 2: Working with Files

```bash
# Copy file contents
cat data.json | clipctl copy

# Format and save back
clipctl transform json > formatted.json
```

### Pattern 3: Using History

```bash
# Copy multiple items
echo "First item" | clipctl copy
echo "Second item" | clipctl copy
echo "Third item" | clipctl copy

# View history
clipctl history

# Retrieve old item
clipctl history -g 2
clipctl paste  # Output: First item
```

### Pattern 4: Chaining Transformations

```bash
# Multiple transformations
echo "hello world" | clipctl copy
clipctl transform upper
clipctl paste | clipctl copy  # Re-copy
clipctl transform base64
clipctl paste  # Output: base64 encoded uppercase text
```

### Pattern 5: Integration with Other Tools

```bash
# With jq
clipctl paste | jq '.data[] | select(.active == true)' | clipctl copy

# With sed
clipctl paste | sed 's/old/new/g' | clipctl copy

# With curl
curl -s https://api.example.com/data | clipctl copy
clipctl transform json
```

### Pattern 6: Creating Custom Workflows

Create a script combining multiple operations:

```bash
#!/bin/bash
# workflow.sh - Process API response

# Get data from clipboard (assumed to be JSON)
clipctl paste > /tmp/data.json

# Process with jq
jq '.items[] | .name' /tmp/data.json | clipctl copy

echo "Extracted names to clipboard"
```

### Pattern 7: Temporary Data Storage

```bash
# Save current clipboard to history
clipctl paste | clipctl copy

# Do other work...
echo "temporary data" | clipctl copy

# Restore from history
clipctl history -g 1
```

## Creating Custom Hooks

### Bash Hook Template

```bash
#!/bin/bash
# Description of what your hook does

# Read from stdin
input=$(cat)

# Process the input
# ... your logic here ...

# Output to stdout
echo "$processed_output"
```

### Python Hook Template

```python
#!/usr/bin/env python3
# Description of what your hook does

import sys

# Read from stdin
text = sys.stdin.read()

# Process the text
# ... your logic here ...
processed = text.upper()  # example

# Output to stdout
print(processed, end='')
```

### Advanced Hook Example: Markdown to HTML

```bash
#!/bin/bash
# markdown-to-html.sh
# Convert markdown to HTML using pandoc

pandoc -f markdown -t html
```

Usage:
```bash
cat README.md | clipctl copy
clipctl hook markdown-to-html.sh > output.html
```

## Tips and Tricks

### 1. Quick Alias Setup

Add to your `.bashrc` or `.zshrc`:

```bash
alias cb='clipctl'
alias cbp='clipctl paste'
alias cbc='clipctl copy'
alias cbh='clipctl history'
```

### 2. Keyboard Shortcuts

If using a desktop environment, bind keyboard shortcuts to:
- `clipctl history` - View history
- `clipctl paste` - Paste to terminal

### 3. Use with tmux

```bash
# Copy tmux buffer to system clipboard
tmux save-buffer - | clipctl copy

# Copy system clipboard to tmux buffer
clipctl paste | tmux load-buffer -
```

### 4. Git Integration

```bash
# Copy last commit message
git log -1 --pretty=%B | clipctl copy

# Copy current branch name
git branch --show-current | clipctl copy

# Copy diff of unstaged changes
git diff | clipctl copy
```

### 5. Data Encryption Hook

Create a hook for encrypting sensitive data:

```bash
#!/bin/bash
# encrypt.sh
# Requires: openssl

# Encrypt using AES-256
openssl enc -aes-256-cbc -salt -pbkdf2
```

And for decryption:

```bash
#!/bin/bash
# decrypt.sh

# Decrypt
openssl enc -d -aes-256-cbc -pbkdf2
```

## Contributing Examples

Have a useful hook or workflow? Please contribute by:

1. Adding your script to `examples/hooks/`
2. Documenting it in this README
3. Submitting a pull request

Make sure your examples are:
- Well-commented
- Documented with usage examples
- Safe to run (no destructive operations without warnings)
