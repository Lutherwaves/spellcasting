#!/usr/bin/env bash
# Render templates/service/, substitute placeholders, run go vet + go build against pinned magic.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MAGIC_VERSION="${MAGIC_VERSION:-v0.17.3}"
TMP="$(mktemp -d -t spellcasting-smoke-XXXXXX)"
trap "rm -rf '$TMP'" EXIT

echo "=== Smoke: rendering templates/service/ to $TMP ==="
cp -r "$REPO_ROOT/templates/service/." "$TMP/"
cd "$TMP"
mv go.mod.tmpl go.mod
mv README.md.tmpl README.md

# Substitute placeholders with smoke values.
find . -type f \( -name '*.go' -o -name '*.yaml' -o -name '*.yml' -o -name 'Makefile' -o -name 'Dockerfile' -o -name 'go.mod' -o -name 'README.md' \) \
  -exec sed -i.bak \
    -e "s/SERVICENAME/smoketest/g" \
    -e "s/SERVICEPORT/8080/g" \
    -e "s/SERVICESCHEMA/smoketest/g" {} \;
find . -name '*.bak' -delete

# Pin magic if MAGIC_VERSION differs from go.mod.tmpl default.
go mod edit -require="github.com/tink3rlabs/magic@${MAGIC_VERSION}"
go mod tidy

echo "=== Smoke: go vet ==="
go vet ./...

echo "=== Smoke: go build ==="
go build ./...

echo "=== Smoke: OK ==="
