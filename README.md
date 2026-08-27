# new-application-name

A starter kit for Go CLI applications: cobra, golangci-lint, gofumpt, unit tests and
end-to-end tests, with a local verification pipeline that needs nothing installed but Go.

## Prerequisites

- [Go](https://go.dev/dl/) 1.26 or newer.

That is the entire list. golangci-lint is pinned in `golangci-lint.mod` and compiled on
demand by the Go toolchain — there is nothing else to install.

## Verification pipeline

Run everything:

```
go run ./tools/check
```

Run a single step:

```
go run ./tools/check lint
```

| # | Step | What it does | Rewrites files |
| - | ---- | ------------ | -------------- |
| 1 | `format` | applies gofumpt + goimports | **yes** |
| 2 | `lint` | golangci-lint: errcheck, govet, ineffassign, staticcheck, unused | no |
| 3 | `unit` | `go test ./...` | no |
| 4 | `build` | compiles to `bin/new-application-name` | no |
| 5 | `e2e` | runs the built binary as a subprocess | no |

Steps run in order and stop at the first failure.

**Step 1 rewrites your source files.** Formatting is applied rather than reported, so you
never have to run a formatter and then re-run the pipeline. The trade-off is that a
passing run can leave you with unstaged changes — "did it pass" and "is my tree clean"
are separate questions here.

**The first run is slow** — a minute or more — because the Go toolchain compiles
golangci-lint from source. Later runs hit the build cache and are quick. It is not hung.

`go run ./tools/check e2e` builds before testing, so it can never pass against a stale
binary.

`go test ./...` never runs the e2e tests; they are behind the `e2e` build tag.

## Running the CLI

```
go run ./cmd/new-application-name hello-world
# Hello, world!
```

Or build first:

```
go run ./tools/check build
./bin/new-application-name hello-world
```

## Project structure

```
cmd/new-application-name/   the main package — thin, just an entry point
internal/cli/               cobra wiring only, no logic
internal/commands/          the logic, as pure functions — unit-tested here
test/e2e/                   tests the built binary (build tag: e2e)
tools/check/                the verification pipeline runner
```

The important rule: **keep logic out of `internal/cli`.** Commands delegate to
`internal/commands`, which is what makes the logic testable without constructing a cobra
command. Follow that split when adding commands.

## Starting a new app

Find-and-replace `new-application-name` across the whole repo, then rename the
`cmd/new-application-name/` directory to match. The name appears in:

- `go.mod` — the module path
- `golangci-lint.mod` — the module path
- `cmd/new-application-name/` — the directory name
- `cmd/new-application-name/main.go` — the import of `internal/cli`
- `internal/cli/root.go` — the root command's `Use` field
- `internal/cli/helloworld.go` — the import of `internal/commands`
- `tools/check/main.go` — the `binaryName` constant
- `test/e2e/cli_test.go` — the binary name it executes
- `README.md` — the title and examples

The import paths matter: renaming the module without updating them leaves the project
uncompilable. A repo-wide replace covers all of the above at once.

`hello-world` is the example command name and is independent of the application name.
Replace it when you add real commands.

## Versioning

`--version` reads `version` in `internal/cli/root.go`, which defaults to `0.0.1`. It is a
`var`, not a `const`, so a release build can override it without touching the source:

```
go build -ldflags "-X github.com/your-org/new-application-name/internal/cli.version=1.2.3" ./cmd/new-application-name
```

The pipeline does not set it. There is deliberately no e2e test asserting the version —
that test would fail on every version bump.

## Adding the race detector

Not enabled by default: on `windows/amd64` it needs `CGO_ENABLED=1` and a C toolchain,
which would break the "only Go required" promise, and it is unsupported on
`windows/arm64`. To enable it, change the `unit` step in `tools/check/main.go` to:

```go
{name: "unit", args: []string{"go", "test", "-race", "./..."}},
```

## Running golangci-lint directly

```
go tool -modfile=golangci-lint.mod golangci-lint run ./...
```

In PowerShell the flag must be quoted, or the `.mod` extension gets split off the
argument:

```powershell
go tool "-modfile=golangci-lint.mod" golangci-lint run ./...
```
