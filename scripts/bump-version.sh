#!/usr/bin/env bash
# Sync version across all manifest files declared in .version-bump.json.
# Usage: bash scripts/bump-version.sh <new-version>
set -euo pipefail
NEW="${1:?usage: bump-version.sh <new-version>}"

jq -r '.files[] | "\(.path)\t\(.field)"' .version-bump.json | while IFS=$'\t' read -r path field; do
  if [[ "$field" == "plugins.0.version" ]]; then
    tmp=$(mktemp)
    jq --arg v "$NEW" '.plugins[0].version = $v' "$path" > "$tmp" && mv "$tmp" "$path"
  else
    tmp=$(mktemp)
    jq --arg v "$NEW" ".${field} = \$v" "$path" > "$tmp" && mv "$tmp" "$path"
  fi
  echo "bumped $path .$field -> $NEW"
done
