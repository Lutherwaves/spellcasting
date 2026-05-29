#!/usr/bin/env bash
# Dispatches to hooks/<name>. Same name kept for Windows-compat parity with superpowers.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
HOOK="$1"
exec "${SCRIPT_DIR}/${HOOK}"
