#!/usr/bin/env python3
"""
Hook to prevent Claude from reading restricted files.
Reads tool input from stdin and exits with code 2 to block the read if the file is restricted.
"""

import json
import sys
import os

BLOCKED_FILES = [
    ".env",
]


def is_blocked(file_path: str) -> bool:
    filename = os.path.basename(file_path)
    return filename in BLOCKED_FILES


def main():
    input_data = json.load(sys.stdin)
    file_path = input_data.get("tool_input", {}).get("file_path", "")

    if is_blocked(file_path):
        print(f"Access denied: reading '{file_path}' is not permitted.", file=sys.stderr)
        sys.exit(2)


if __name__ == "__main__":
    main()
