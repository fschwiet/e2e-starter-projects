# Go CLI Starter Kit — Design

## Overview

A minimal "hello, world" Go CLI application on the `golang-cli` branch, with a local
verification pipeline covering formatting, linting, unit tests and end-to-end tests.

This is structurally parallel to the `npm-command` starter kit — thin entry point,
logic in testable functions, one aggregate `check` command — with Go-native tooling
substituted throughout.

The defining constraint: **`go` is the only prerequisite.** Every other tool is pinned
in-repo and compiled on demand. A developer with a fresh Go install can clone the branch
and run the full pipeline without installing anything else.

## Terminology

This repo uses four names for the same concept across its branches (`kits` in the repo
name, "projects"/"environments" in the `main` README, "template" in `npm-command`,
"pack" in `android-gradle`). This starter standardizes on:

- **starter kit** — the minimal project contained in one branch
- **environment** — what distinguishes one starter kit from another (Go CLI, npm CLI,
  Electron, Android)

No repo-wide glossary is introduced here. Because starter kits live on branches that
never merge, a shared glossary belongs on `main` and is out of scope for this work.

## Toolchain

| Concern | Choice |
| --- | --- |
| Language | Go 1.26 (`go 1.26.7` directive in `go.mod`) |
| CLI parsing | `spf13/cobra` v1.10.2 |
| Lint | golangci-lint v2, pinned as a tool dependency in a dedicated modfile |
| Format | gofumpt + goimports, run as golangci-lint v2 formatters |
| Unit tests | stdlib `testing`, files alongside the code |
| E2E tests | stdlib `testing`, `os/exec` against the built binary, `e2e` build tag |
| Pipeline runner | a Go program at `tools/check` |
| CI | none |

Rationale for the non-obvious choices is recorded in `docs/adr/`.

## Repository layout

```
go.mod                                  module github.com/your-org/new-application-name
go.sum
golangci-lint.mod                       pinned golangci-lint tool dependency
golangci-lint.sum
.golangci.yml                           formatters, linters, build tags
.gitignore
README.md
AGENTS.md                               pipeline documentation for agents
CLAUDE.md                               one line: @AGENTS.md
cmd/new-application-name/main.go        thin entry point
internal/cli/root.go                    cobra root command, version
internal/cli/helloworld.go              hello-world subcommand wiring
internal/commands/helloworld.go         Greet() string — pure logic
internal/commands/helloworld_test.go    unit test
test/e2e/cli_test.go                    e2e tests (//go:build e2e)
tools/check/main.go                     verification pipeline runner
docs/adr/                               architecture decision records
docs/honist-v/specs/                    design specs
```

### Layout rationale

**`cmd/` holds the `main` package, not the cobra commands.** `cobra-cli` scaffolds
`main.go` at the repo root with commands in a `cmd/` package, which overloads a directory
name that conventionally holds `main` packages in Go. Putting the binary at
`cmd/new-application-name/` and the cobra wiring in `internal/cli/` keeps both
conventions intact.

**Unit tests live next to the code they test,** following Go convention, rather than in a
mirrored `tests/unit/` tree as in the `npm-command` starter. E2E tests are the exception
and get their own directory, because they test the binary rather than any package.

## Components

- **`cmd/new-application-name/main.go`** — Calls `cli.Execute()`. If it returns an error,
  exits with status 1. Contains no other logic.

- **`internal/cli/root.go`** — Constructs the cobra root command (`Use`, `Short`,
  `Version`), registers subcommands, exposes `Execute() error`. Declares
  `var version = "0.0.1"`.

- **`internal/cli/helloworld.go`** — Constructs the `hello-world` cobra subcommand. Its
  `RunE` calls `commands.Greet()` and writes the result to `cmd.OutOrStdout()`, returning
  any write error. No business logic.

- **`internal/commands/helloworld.go`** — Exports `Greet() string` returning
  `"Hello, world!"`. Pure, no I/O, no cobra dependency.

- **`tools/check/main.go`** — The verification pipeline runner. See below.

The load-bearing boundary is between `internal/cli` and `internal/commands`: **cobra
wiring never contains logic.** `Greet()` is testable without constructing a cobra
command; `internal/cli` only translates between cobra and that function. This split is
the pattern the starter kit exists to demonstrate, and every command added later should
follow it.

## Data flow

```
argv → cobra → subcommand RunE → commands.Greet() → cmd.OutOrStdout()
```

## Version handling

`internal/cli/root.go` declares `var version = "0.0.1"`, wired to cobra's `Version`
field so `--version` works. It is a `var` rather than a `const` specifically so it can be
overridden at build time:

```
go build -ldflags "-X github.com/your-org/new-application-name/internal/cli.version=1.2.3" ./cmd/...
```

The pipeline does not set it. Go has no equivalent of `package.json` as a single source
of truth for a version string, and computing one from `git describe` would add release
tooling this starter kit has no use for yet. The README documents the `-ldflags` seam as
the growth path.

## Tool pinning

golangci-lint is pinned as a Go tool dependency in a **dedicated modfile** rather than in
the application's `go.mod`, per golangci-lint's own documented guidance — its dependency
tree is very large and would otherwise entangle the application's.

`golangci-lint.mod` (committed, along with `golangci-lint.sum`) contains:

- `module github.com/your-org/new-application-name/golangci-lint`
- a `go` directive matching the toolchain (`go 1.26.0`)
- `tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint`
- the resulting `require` block (200+ indirect dependencies, generated by `go get -tool`)

Pinned version at time of writing: **v2.13.1**.

Invocation:

```
go tool -modfile=golangci-lint.mod golangci-lint run ./...
go tool -modfile=golangci-lint.mod golangci-lint fmt ./...
```

This was verified working against Go 1.26.7: the invocation correctly resolves and
type-checks the *application* module (not the tool module) despite `-modfile` substituting
the main module for tool resolution.

### Implementation notes

These are quirks confirmed empirically that an implementer will otherwise rediscover the
hard way:

1. **PowerShell strips the `.mod` extension** from an unquoted `-modfile=golangci-lint.mod`
   argument, splitting it into `-modfile=golangci-lint` and `.mod`. This surfaces as
   `file does not have .mod extension` from `go get`, or as
   `'go mod init' accepts at most one argument` from `go mod init` (the stray `.mod`
   becomes a second positional argument). The flag must be quoted in PowerShell:
   `go tool "-modfile=golangci-lint.mod" golangci-lint run ./...`. The pipeline runner is a
   Go program using `os/exec` and is unaffected, but README examples a user might paste
   into PowerShell must show the quoted form.

2. **Setup commands** (run once when building the starter kit, not by template users):

   ```
   go mod init -modfile=golangci-lint.mod github.com/your-org/new-application-name/golangci-lint
   go get -tool -modfile=golangci-lint.mod github.com/golangci/golangci-lint/v2@v2.13.1
   ```

   Verified working on Go 1.26.7. `go mod init -modfile=...` is supported — quote the flag
   if running it from PowerShell (see note 1).

3. **Both `golangci-lint.mod` and `golangci-lint.sum` must be committed.** The `require`
   graph in the modfile is what selects versions; the `.sum` records integrity hashes.
   Without the `.sum`, builds must re-resolve hashes and lose verification against
   tampering, so both belong in version control.

4. **The first run compiles golangci-lint from source** and is slow (order of minutes).
   Subsequent runs hit Go's build cache and are fast. The README should set this
   expectation so a first-time user does not assume the pipeline has hung.

## Lint and format configuration

`.golangci.yml`:

```yaml
version: "2"
run:
  build-tags:
    - e2e
formatters:
  enable:
    - gofumpt
    - goimports
```

- **Formatters:** `gofumpt` (a stricter superset of gofmt — more opinionated, which
  removes style debate from a starter kit) plus `goimports` (import grouping and pruning,
  which gofumpt does not do).

- **Linters:** golangci-lint's default set — `errcheck`, `govet`, `ineffassign`,
  `staticcheck`, `unused`. Deliberately not expanded. A starter kit shipping thirty
  enabled linters makes the user's first act deleting linters; adding them as a project
  matures is the easier direction.

- **`run.build-tags: [e2e]`** is required, not optional. Verified: without it,
  golangci-lint does not see `test/e2e/` at all (the tagged files are excluded from the
  build), so e2e test code would silently rot unlinted. With it, the tagged files are
  linted normally.

## Verification pipeline

`tools/check` is a Go program invoked as:

```
go run ./tools/check              # all steps, in order, fail-fast
go run ./tools/check <step>       # a single step
```

| # | Step name | Command | Mutates files |
| --- | --- | --- | --- |
| 1 | `format` | `go tool -modfile=golangci-lint.mod golangci-lint fmt ./...` | **yes** |
| 2 | `lint` | `go tool -modfile=golangci-lint.mod golangci-lint run ./...` | no |
| 3 | `unit` | `go test ./...` | no |
| 4 | `build` | `go build -o bin/new-application-name[.exe] ./cmd/new-application-name` | no |
| 5 | `e2e` | `go test -tags e2e ./test/e2e/...` | no |

The only hard ordering constraint is that `build` precedes `e2e`. The rest is ordered
cheapest-first for fast feedback.

**Step 1 mutates the working tree.** This is a deliberate departure from the sibling
starter kits, which check formatting and fail. Auto-applying formatting removes the
friction of running a separate format command and re-running the pipeline. The tradeoff
accepted: "did the pipeline pass" and "is my working tree clean" become separate
questions. Because formatting is applied before step 2, `lint` never reports a gofumpt
violation — they have already been fixed.

There is no read-only variant of the pipeline, because there is no CI to need one.

Should one be wanted later, formatting would still be *verified* rather than silently
skipped. Verified against golangci-lint v2.13.1: `golangci-lint run` reports configured
formatters' violations as ordinary issues (`File is not properly formatted (gofumpt)`,
exit 1) even though it does not rewrite files. So running steps 2–5 and omitting step 1
is a complete read-only gate. `golangci-lint fmt --diff ./...` is the other option — it
prints a unified diff, exits 1 when changes are needed, and leaves files untouched.

The runner deliberately offers no single "everything except format" command; invoking the
individual steps covers it, and adding a mode nothing currently uses would be speculative.

### Runner behaviour

- Steps run in the order above. On the first failing step, the runner prints which step
  failed, propagates the step's stdout/stderr, and exits non-zero immediately. Remaining
  steps do not run.
- Each step's output streams to the runner's own stdout/stderr rather than being buffered,
  so long-running steps show progress.
- The runner prints a short banner per step (e.g. `==> lint`) so failures are attributable
  when reading scrollback.
- Invoked with an argument, it runs only that named step, with one exception: **`e2e`
  runs `build` first.** Otherwise `go run ./tools/check e2e` would test whatever stale
  binary happens to be in `bin/`, and could report a pass for source that no longer
  compiles or that produces different output. Detecting only the *missing*-binary case
  does not catch this. `build` is cheap and incremental, so the cost is negligible.
- An unrecognized step name is an error listing the valid names, exit non-zero.
- The binary name gets a `.exe` suffix when `runtime.GOOS == "windows"`.
- The runner uses only the standard library (`os`, `os/exec`, `runtime`), so it adds no
  dependency to `go.mod`.

## Error handling

- **`hello-world` success:** prints `Hello, world!` to stdout, exit 0.
- **`--version`:** cobra prints the version, exit 0.
- **Unknown command or flag:** cobra prints to stderr and exits 1. The two cases differ,
  verified against cobra v1.10.2 — stdout is empty in both:

  | Input | stderr |
  | --- | --- |
  | unknown command | `Error: unknown command "x" for "app"` + `Run 'app --help' for usage.` — **no usage block** |
  | unknown flag | `Error: unknown flag: --x` + the **full usage block** |

  No custom cobra output configuration is needed; the defaults already send both to
  stderr. The e2e test asserts on the unknown-*command* case only (see Testing).
- **`main.go`:** on error from `cli.Execute()`, exits 1. cobra has already printed the
  message, so `main` does not print it again.
- **E2E with no built binary:** the test fails with an actionable message naming the fix
  (`run "go run ./tools/check build" first`) rather than an opaque exec error.
- **`errcheck` and cobra output:** cobra's idiomatic
  `fmt.Fprintln(cmd.OutOrStdout(), ...)` returns an error that is conventionally ignored,
  which `errcheck` flags. The subcommand uses `RunE` and returns that error rather than
  adding an errcheck exclusion — a starter kit should not teach "configure the linter
  away" as the first move.

## Testing

### Unit tests

`internal/commands/helloworld_test.go` asserts `Greet()` returns `"Hello, world!"`. Run
via `go test ./...`, which is fast and passes on a clean checkout because the e2e build
tag excludes the e2e package.

`-race` is deliberately not enabled. On `windows/amd64` it requires `CGO_ENABLED=1` and a
working C toolchain, which would break the "`go` is the only prerequisite" property that
shapes this entire design; on `windows/arm64` the race detector is not supported at all. A
hello-world CLI has no concurrency for it to find. The README notes how to add it.

### E2E tests

`test/e2e/cli_test.go` carries `//go:build e2e`. Without the tag, a developer typing the
natural `go test ./...` would pick up e2e tests and get a confusing failure because the
binary has not been built.

Verified against Go 1.26.7: with every file in `test/e2e/` excluded by the tag, both
`go test ./...` and `go build ./...` skip the package silently and exit 0 — a package with
no buildable files is not an error under a `./...` pattern. `go test -tags e2e
./test/e2e/...` then runs it normally.

Tests locate the binary at `../../bin/new-application-name[.exe]` — Go runs tests with the
working directory set to the package directory, so this resolves deterministically. No
environment variable indirection.

Two tests, chosen as **exemplars rather than for coverage**. Their purpose is to show a
template user where and how to add tests:

1. `hello-world` → stdout is exactly `Hello, world!\n`, exit code 0.
2. `bogus` (an unknown command) → exit code 1, stderr *contains* `unknown command`,
   stdout empty.

Test 2 asserts a substring rather than cobra's exact stderr text, and uses the unknown-
*command* case rather than unknown-*flag*, because the command case produces a short help
hint with no usage block — a stable assertion that will not break when subcommands or
flags are added. Asserting the full usage block (which the unknown-flag case emits) would
make the test fail every time the template user adds a command.

Together they cover the two shapes an e2e assertion takes: the success path via stdout,
and the failure path via stderr and exit code.

There is deliberately **no `--version` e2e test.** Asserting the version string creates a
test that fails on every version bump, which teaches the wrong lesson in a template.

## Documentation

**`README.md`** — overview; prerequisites (Go 1.26+, and *only* Go); the verification
pipeline table; a note that step 1 rewrites files; a note that the first run compiles
golangci-lint and is slow; how to build and run the CLI locally; the `-ldflags` version
seam; how to add `-race`; and a "Starting a new app" section listing every location the
`new-application-name` placeholder appears.

**`AGENTS.md`** — the pipeline commands and when to use each, written for coding agents,
following the `android-gradle` starter kit's precedent. Explicitly states that
`go run ./tools/check` mutates files. Without this, an agent unaware of `tools/check`
will invent an inferior ad-hoc pipeline.

**`CLAUDE.md`** — a single line, `@AGENTS.md`, matching `android-gradle`.

## Placeholder naming

The placeholder application name is `new-application-name`, consistent with the other
starter kits. It appears in:

- `go.mod` — the `module` path
- `golangci-lint.mod` — the `module` path
- `cmd/new-application-name/` — the directory name
- `cmd/new-application-name/main.go` — the **import** of `.../internal/cli`
- `internal/cli/helloworld.go` — the **import** of `.../internal/commands`
- `internal/cli/root.go` — the root command's `Use` field
- `tools/check/main.go` — the `-o bin/...` output name
- `test/e2e/cli_test.go` — the binary path being executed
- `README.md` — the title, references, and the `-ldflags` example path

The import paths matter: renaming the module without updating them leaves the starter kit
uncompilable, so the README's "Starting a new app" list must include them. A
project-wide find-and-replace of `new-application-name` covers every entry above, and the
README should recommend exactly that rather than a manual file-by-file walk.

`hello-world` remains the example subcommand name and is independent of the application
name, matching `npm-command`. Replacing it is left to the developer when they add real
commands.

`.gitignore` covers `bin/` plus standard Go entries.

## Out of scope (YAGNI)

- CI workflow — the `npm-command` starter kit's most recent commit removed its GitHub
  Actions workflow, establishing local-only verification as the current convention.
- Release automation, `git describe` versioning, or goreleaser.
- A read-only variant of the pipeline.
- `-race` in the default pipeline.
- Additional example commands, subcommand groups, config files, or interactive prompts.
- A repo-wide `CONTEXT.md` glossary — belongs on `main`, not on a starter kit branch.

## Decision records

- `docs/adr/0001-pin-golangci-lint-via-tool-modfile.md`
- `docs/adr/0002-go-program-as-pipeline-runner.md`
- `docs/adr/0003-format-mutates-inside-the-pipeline.md`
