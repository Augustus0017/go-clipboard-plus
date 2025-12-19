#!/usr/bin/env python3
# Example hook: Count words and characters
# Place this in ~/.config/go-clipboard-plus/hooks/count-stats.py

import sys

text = sys.stdin.read()
words = len(text.split())
chars = len(text)
lines = len(text.splitlines())

print(f"Statistics:")
print(f"  Lines: {lines}")
print(f"  Words: {words}")
print(f"  Characters: {chars}")
