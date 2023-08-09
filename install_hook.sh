#!/bin/bash

HOOK_DIR=".git/hooks"
HOOK_FILE="${HOOK_DIR}/pre-commit"

# Check if we're in a git repository
if [ ! -d "${HOOK_DIR}" ]; then
    echo "Error: Not a git repository or no .git directory found"
    exit 1
fi

# Define the hook content
HOOK_CONTENT='#!/bin/sh

# Redirect output to stderr.
exec 1>&2

# Check for unformatted Go files.
unformatted=$(git diff --cached --name-only --diff-filter=ACM | grep ".go$" | xargs gofmt -l)
if [ ! -z "$unformatted" ]; then
    echo >&2 "Go files must be formatted with gofmt. Unformatted files:"
    echo >&2 "$unformatted"
    exit 1
fi
'

# Write the hook content to the pre-commit file
echo "${HOOK_CONTENT}" > "${HOOK_FILE}"
chmod +x "${HOOK_FILE}"

echo "pre-commit hook for go fmt installed successfully!"
