#!/bin/sh
# Managed by commitkit – do not edit by hand unless you know what you are doing.
# https://github.com/destyk/commitkit

set -e

if command -v commitkit >/dev/null 2>&1; then
	exec commitkit check --file "$1"
fi

# Fallback: look for commitkit next to the hook (useful during development).
HOOK_DIR=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd -- "$HOOK_DIR/../.." && pwd)

if [ -x "$REPO_ROOT/commitkit" ]; then
	exec "$REPO_ROOT/commitkit" check --file "$1"
fi

if [ -x "$REPO_ROOT/bin/commitkit" ]; then
	exec "$REPO_ROOT/bin/commitkit" check --file "$1"
fi

echo "commitkit: executable not found in PATH or repository root" >&2
echo "Install commitkit or add it to PATH." >&2
exit 1