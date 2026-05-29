---
description: Route layer; auto-chains casting-a-service-layer if service is missing
argument-hint: <resource-name>
---

Invoke the `casting-a-route` skill via the Skill tool. Pass `$ARGUMENTS` as the resource name. If empty, ask. The skill detects whether the service layer exists and chains `casting-a-service-layer` first if needed (unless the user passes `--custom`).
