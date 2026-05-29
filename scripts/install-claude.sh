#!/usr/bin/env bash
cat <<'EOF'
Claude Code installs spellcasting via its plugin marketplace, not this script.

Run in Claude Code:
  /plugin marketplace add https://github.com/Lutherwaves/spellcasting
  /plugin install spellcasting@magic
  /reload-plugins

After install, the SessionStart hook auto-loads using-sorcery.
EOF
