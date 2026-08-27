# Formatting is applied inside the pipeline, not checked

The first step of `go run ./tools/check` runs `golangci-lint fmt ./...`, which **rewrites
source files**. The pipeline is therefore not read-only, and there is no read-only variant.

This is a deliberate departure from the sibling starter kits — `npm-command` runs
`prettier --check` and fails; `android-gradle` separates a mutating `precommit` from a
read-only `ci`. It is recorded here so nobody "fixes" it back.

## Why

Having formatting fail the pipeline means a developer runs `check`, watches it fail on
whitespace, runs a separate format command, and runs `check` again. Auto-applying removes
that loop. Because formatting runs before linting, `lint` never reports a gofumpt
violation — they have already been fixed.

## Consequences

"Did the pipeline pass" and "is my working tree clean" become separate questions: a
passing `check` can leave unstaged changes. This is acceptable because the pipeline is
local-only — there is no CI (per the `npm-command` starter kit's most recent commit, which
removed its workflow) that would need a non-mutating variant. Should CI be added later,
skipping step 1 yields one, and the README says so.
