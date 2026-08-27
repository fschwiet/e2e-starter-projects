# Working in this repo

A Go CLI starter kit. `go` is the only prerequisite — golangci-lint is pinned in
`golangci-lint.mod` and compiled on demand. Do not install tooling.

## Verify a change

Run this before considering any change done:

```
go run ./tools/check
```

Steps run in order and stop at the first failure:

| # | Step | What it checks |
| - | ---- | -------------- |
| 1 | `format` | applies gofumpt + goimports — **rewrites source files** |
| 2 | `lint` | golangci-lint: errcheck, govet, ineffassign, staticcheck, unused |
| 3 | `unit` | `go test ./...` |
| 4 | `build` | compiles the binary into `bin/` |
| 5 | `e2e` | runs the built binary as a subprocess |

Run one step with `go run ./tools/check <step>`. `e2e` rebuilds first, so it never tests
a stale binary.

**`go run ./tools/check` mutates the working tree** — step 1 applies formatting. Unstaged
changes after a run are expected, not a failure.

`tools/check` is the entry point. Do not invent an ad-hoc pipeline.

The first run compiles golangci-lint from source and takes minutes. It is not hung.

## Where code goes

- `internal/commands/` — logic, as pure functions. No cobra, no I/O. Unit-tested here,
  in `*_test.go` files beside the code.
- `internal/cli/` — cobra wiring only. Each command calls one function and writes the
  result.
- `cmd/new-application-name/` — the `main` package. Thin.
- `test/e2e/` — tests the built binary. Tagged `//go:build e2e`.
- `tools/check/` — the pipeline runner. Standard library only; it must add nothing to
  `go.mod`.

Keep logic out of `internal/cli`. That split is what makes commands testable without
constructing a cobra command, and it is the pattern this starter kit exists to show.

## Gotchas

- `go test ./...` does not run the e2e tests — they are behind the `e2e` build tag. Use
  `go run ./tools/check e2e`.
- Invoking golangci-lint directly from PowerShell needs the flag quoted, or the `.mod`
  extension is split off the argument:
  `go tool "-modfile=golangci-lint.mod" golangci-lint run ./...`
- If errcheck flags an unchecked error, handle the error. Do not add an exclusion to
  `.golangci.yml`.

## Design decisions

`docs/adr/` records why the non-obvious choices were made — the pinned tool modfile, the
Go pipeline runner, and the mutating format step. Read before changing any of them.
