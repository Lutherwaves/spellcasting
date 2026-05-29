# Cursor Tool Mapping

Cursor does not expose a named tool surface the way Claude Code does. Skills in Cursor work as **rules** (`.cursor/rules/*.mdc`), not as invocable tools. This file explains how Claude Code skill concepts map to the Cursor workflow.

## How spellcasting works in Cursor

| Skill references | Cursor equivalent |
|-----------------|-------------------|
| `Read` (read a file) | Cursor attaches files to context automatically via `@<filename>` or glob rules |
| `Write` / `Edit` (modify files) | Cursor generates diffs in chat; you accept or reject them |
| `Bash` (run commands) | Cursor Terminal tab — run commands manually or via the integrated terminal |
| `Skill` (invoke a skill) | Reference the rule in chat with `@<rule-name>` (e.g. `@using-sorcery`) |
| `Agent` / `Task` (dispatch subagent) | Not natively supported — describe the subtask in chat and proceed sequentially |
| `TodoWrite` (task tracking) | Maintain a checklist manually in chat or in a scratch file |
| `WebFetch` (fetch a URL) | Use `@Web` in Cursor chat to pull live context from a URL |
| `Grep` / `Glob` (search) | Use `@Codebase` in chat to search across the project |

## Installing skills as Cursor rules

Skills are converted to `.cursor/rules/*.mdc` files by `scripts/install-cursor.sh`. Run it from your project root:

```bash
bash /path/to/spellcasting/scripts/install-cursor.sh
```

This generates one `.mdc` file per skill under `.cursor/rules/`.

## Triggering rules in Cursor

- **On demand**: type `@<rule-name>` in the Cursor chat (e.g. `@using-sorcery`, `@divining-intent`).
- **Automatic**: set `alwaysApply: true` in the frontmatter to attach the rule to every chat. Avoid this for most spellcasting skills — they are cast-specific and should only activate when relevant.
- **By file pattern**: set `globs: "**/*.go"` in the frontmatter to attach the rule whenever Go files are in context.

## Key difference from Claude Code

In Claude Code, `Skill` is a first-class tool call — the agent reads the skill file and follows its instructions programmatically. In Cursor, the rule body is injected as additional context into the model's prompt when triggered. The instructions are the same; the delivery mechanism differs.
