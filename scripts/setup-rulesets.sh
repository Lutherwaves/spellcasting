#!/usr/bin/env bash
# Apply repository ruleset to main using GitHub's newer Rulesets API.
# Rulesets supersede classic branch protection — keep both scripts during
# transition, but prefer this one for new repos.
#
# Idempotent: looks up the ruleset by name and PUTs an update if it exists,
# POSTs a new one otherwise. Safe to re-run.
#
# Usage: REPO=Lutherwaves/spellcasting bash scripts/setup-rulesets.sh
set -euo pipefail

REPO="${REPO:-Lutherwaves/spellcasting}"
RULESET_NAME="main-protection"

read -r -d '' PAYLOAD <<JSON || true
{
  "name": "${RULESET_NAME}",
  "target": "branch",
  "enforcement": "active",
  "bypass_actors": [
    { "actor_id": 5, "actor_type": "RepositoryRole", "bypass_mode": "always" }
  ],
  "conditions": {
    "ref_name": {
      "include": ["~DEFAULT_BRANCH"],
      "exclude": []
    }
  },
  "rules": [
    { "type": "deletion" },
    { "type": "non_fast_forward" },
    {
      "type": "pull_request",
      "parameters": {
        "required_approving_review_count": 1,
        "dismiss_stale_reviews_on_push": true,
        "require_code_owner_review": false,
        "require_last_push_approval": false,
        "required_review_thread_resolution": true
      }
    },
    {
      "type": "required_status_checks",
      "parameters": {
        "strict_required_status_checks_policy": true,
        "required_status_checks": [
          { "context": "Conventional Commits Check" },
          { "context": "Lint" },
          { "context": "Template Smoke" }
        ]
      }
    },
    { "type": "required_linear_history" }
  ]
}
JSON

# Look up existing ruleset id by name.
existing_id=$(
  gh api "repos/${REPO}/rulesets" --jq \
    ".[] | select(.name == \"${RULESET_NAME}\") | .id" 2>/dev/null \
  | head -1
)

if [ -n "${existing_id:-}" ]; then
  echo "Updating existing ruleset ${existing_id} on ${REPO}"
  printf '%s' "$PAYLOAD" | gh api -X PUT "repos/${REPO}/rulesets/${existing_id}" --input -
else
  echo "Creating new ruleset on ${REPO}"
  printf '%s' "$PAYLOAD" | gh api -X POST "repos/${REPO}/rulesets" --input -
fi

echo
echo "Ruleset '${RULESET_NAME}' applied to ${REPO} (target: default branch)"
echo "Bypass: repo admins (matches enforce_admins:false on classic)"
