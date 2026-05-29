#!/usr/bin/env bash
# Apply branch protection to main. Idempotent — safe to re-run.
set -euo pipefail

REPO="${REPO:-Lutherwaves/spellcasting}"

gh api -X PUT "repos/${REPO}/branches/main/protection" \
  --input - <<'JSON'
{
  "required_status_checks": {
    "strict": true,
    "contexts": ["Conventional Commits Check", "Lint", "Template Smoke"]
  },
  "enforce_admins": false,
  "required_pull_request_reviews": {
    "required_approving_review_count": 1,
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": false
  },
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_conversation_resolution": true
}
JSON

echo "Branch protection applied to ${REPO}@main"
