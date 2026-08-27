# The verification pipeline runner is a Go program

Go has no built-in equivalent of npm scripts, so the aggregate `check` command needs a
host. It is a plain Go program at `tools/check`, run as `go run ./tools/check`, using only
the standard library.

This follows directly from [ADR-0001](./0001-pin-golangci-lint-via-tool-modfile.md): once
`go` is the only prerequisite, every candidate runner except a Go program partially gives
that property back.

## Considered options

- **Make** — not installed on the Windows machines this repo targets, and not a reasonable
  prerequisite there.
- **go-task (`Taskfile.yml`)** — declarative and closest in spirit to the `npm-command`
  starter kit's `package.json` scripts, with good discoverability via `task --list`. But
  the runner itself then needs bootstrapping: pinning it the way ADR-0001 pins the linter
  means typing `go tool -modfile=tools/task.mod task check`, which is worse than what it
  replaces, and installing it globally reintroduces the prerequisite we just removed.
- **mage** — Go-based like the chosen option, with target discovery for free, but adds a
  dependency and a magefile build-tag convention to save roughly forty lines.
- **Paired `check.ps1` / `check.sh`** — matches the `android-gradle` starter kit's script
  style, but means two files to keep in sync and makes the starter kit Windows-first
  rather than cross-platform.

## Consequences

Step sequencing, output streaming and single-step dispatch are hand-written rather than
provided. In exchange the pipeline is cross-platform for free, adds nothing to `go.mod`,
is debuggable in the same language as the project, and can be read top to bottom in one
file. `go run ./tools/check` is more typing than `task check`.
