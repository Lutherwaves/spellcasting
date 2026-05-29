# PR

## What changed

<!-- one sentence -->

## Cast affected

- [ ] new skill: <name>
- [ ] modified skill: <name>
- [ ] new template: <path>
- [ ] modified template: <path>
- [ ] platform adapter (Codex / Cursor / install scripts)
- [ ] CI / release / branch protection
- [ ] docs only

## Checklist

- [ ] Conventional Commits (`feat:` / `fix:` / `chore:` / `docs:` / `refactor:` / `test:` / `ci:`)
- [ ] **Capability-boundary check** — no project-specific common-library references, custom rate-limit/audit/cache middleware, or per-resource ACL authorizers introduced in skills or templates
- [ ] `bash tests/template-smoke.sh` runs locally (if templates changed)
- [ ] Skill frontmatter YAML parses (if skill changed)
- [ ] `magic-capabilities` updated if a new magic symbol is referenced
- [ ] `using-sorcery` cast catalog updated if a new cast was added
