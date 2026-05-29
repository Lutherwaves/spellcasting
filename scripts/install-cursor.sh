#!/usr/bin/env bash
# Generate .cursor/rules/ files from skills/ and install into a target project.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGET="${1:-$PWD}"

if [ ! -d "$TARGET" ]; then
  echo "error: target directory '$TARGET' does not exist" >&2
  exit 1
fi

bash "$REPO_ROOT/scripts/generate-cursor-rules.sh" "$TARGET"
echo
echo "✓ Installed Cursor rules to $TARGET/.cursor/rules/"
echo "  Trigger a rule in Cursor chat with @<rule-name> (e.g. @using-sorcery)."
