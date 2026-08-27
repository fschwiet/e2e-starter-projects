# Pin golangci-lint via a dedicated tool modfile

The defining constraint for this starter kit is that `go` is the only prerequisite: clone
the branch, run the pipeline, install nothing else. golangci-lint is therefore pinned as a
Go tool dependency in a committed `golangci-lint.mod` (plus `.sum`) and invoked as
`go tool -modfile=golangci-lint.mod golangci-lint …`, which compiles it on first use and
caches it thereafter.

## Considered options

- **Global install** (winget/scoop/install script) is the path golangci-lint most
  endorses, and it starts up faster. Rejected because it grows the prerequisite list and
  lets linter versions drift between machines, producing "works on my machine" lint
  failures. The sibling `npm-command` starter kit pins every tool in its lockfile and
  requires only Node and pnpm; this gives Go the same property.
- **A `tool` directive in the application's `go.mod`** is simplest to write, but
  golangci-lint explicitly discourages it: its dependency tree is enormous (200+ modules)
  and would entangle the application's, so a `go get -u` on the app could silently produce
  an untested linter build. The dedicated modfile is golangci-lint's own documented
  recommendation for exactly this reason.

## Consequences

The first pipeline run compiles golangci-lint from source and takes minutes; this needs
calling out in the README so users do not assume a hang. The invocation is wordier than a
bare `golangci-lint`, which the pipeline runner hides.

One Windows gotcha, since it produces confusing errors: PowerShell splits an unquoted
`-modfile=golangci-lint.mod` argument at the extension, yielding either
`file does not have .mod extension` or `'go mod init' accepts at most one argument`. Quote
the flag. The pipeline runner shells out via `os/exec` and is unaffected.
