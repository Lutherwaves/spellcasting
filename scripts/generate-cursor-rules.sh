#!/usr/bin/env bash
# Generate .cursor/rules/*.mdc from skills/<name>/SKILL.md.
# The mdc files mirror the SKILL.md body with Cursor-flavored frontmatter.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGET="${1:-$REPO_ROOT}"
RULES_DIR="$TARGET/.cursor/rules"

mkdir -p "$RULES_DIR"

for skill in "$REPO_ROOT"/skills/*/; do
  name="$(basename "$skill")"
  skill_md="$skill/SKILL.md"
  [ -f "$skill_md" ] || continue

  # Extract description from the SKILL.md frontmatter.
  desc=$(awk '/^description:/{sub(/^description: */, ""); print; exit}' "$skill_md")

  mdc="$RULES_DIR/$name.mdc"
  {
    echo "---"
    echo "description: $desc"
    echo "alwaysApply: false"
    echo "---"
    echo
    # Body: skip the YAML frontmatter (everything between the first two `---`).
    awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2{print}' "$skill_md"
  } > "$mdc"

  echo "generated $name.mdc"
done
