#!/usr/bin/env bash
# Install spellcasting skills into ~/.agents/skills/ for Codex CLI.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEST="${AGENTS_DIR:-$HOME/.agents}/skills"

mkdir -p "$DEST"

for skill in "$REPO_ROOT"/skills/*/; do
  name="$(basename "$skill")"
  target="$DEST/$name"
  if [ -L "$target" ] || [ -d "$target" ]; then
    rm -rf "$target"
  fi
  ln -s "$skill" "$target"
  echo "linked $name → $target"
done

echo
echo "✓ Installed $(ls -1 "$REPO_ROOT/skills" | wc -l) skills to $DEST"
echo "  Codex will auto-discover them on next launch."
echo "  Tool mapping reference: $REPO_ROOT/references/codex-tools.md"
