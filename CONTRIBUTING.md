# Contributing to spellcasting

## Workflow

Trunk-based. Short-lived topic branches off `main`. Naming: `feat/<desc>`, `fix/<desc>`, `chore/<desc>`, `docs/<desc>`, `ci/<desc>`, `refactor/<desc>`, `test/<desc>`.

## Conventional Commits

Required on every commit. Allowed types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `ci`.

| Type | Release impact |
|---|---|
| `feat:` | Minor release |
| `fix:` | Patch release |
| `chore:`, `docs:`, `refactor:`, `test:`, `ci:` | No release |

Releases are cut automatically by [`go-semantic-release`](https://github.com/go-semantic-release/semantic-release) on every push to `main` that contains a `feat:` or `fix:` since the last tag.

## Branch protection

`main` requires:
- 1 PR approval
- All three CI checks green: **Conventional Commits Check**, **Lint**, **Template Smoke**
- Conversation resolution
- No force-pushes, no deletions

## Adding a new cast

Each cast is two things:
1. A **skill** at `skills/<cast-name>/SKILL.md` (markdown body with the cast procedure, prerequisites, capability boundary). Spellcasting is skills-only — there are no slash commands. The agent picks the right cast from each skill's description.
2. **Templates** (where applicable) under `templates/<layer>/` that the skill tells the agent to read and adapt.

Update `using-sorcery` to list the new cast in its catalog. Update `magic-capabilities` if the cast references new magic symbols.

## Capability-boundary check on every PR

The PR template's checklist asks: *"No project-specific common-library references introduced in templates or skills?"* Reviewers verify by grepping the diff for any of: project-specific package names that wrap magic, custom middleware that mirrors magic's functionality, or per-resource ACL authorizers beyond `RequireRole`.

## Magic version bumps

When magic cuts a new tag:

1. Bump `MAGIC_VERSION` in `tests/template-smoke.sh`.
2. Bump the version pin in `templates/service/go.mod.tmpl`.
3. Update the `Pinned to **vX.Y.Z**` line in `skills/magic-capabilities/SKILL.md`.
4. Re-fetch each `pkg.go.dev` URL in `magic-capabilities`; any 404 → symbol renamed or removed; reconcile.
5. Update spellcasting's README magic-version compatibility table.

Commit as `chore: bump magic to vX.Y.Z`. Releases this as a patch.

## Local testing

```bash
# Validate plugin manifests
jq empty .claude-plugin/plugin.json .claude-plugin/marketplace.json

# Validate skill frontmatter
for f in skills/*/SKILL.md; do
  python3 -c "import yaml; yaml.safe_load(open('$f').read().split('---')[1])"
done

# Run the template smoke (renders templates/service/ and smoke-builds)
bash tests/template-smoke.sh
```

## License

MIT, same as the project.
