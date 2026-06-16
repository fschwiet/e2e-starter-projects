# NPM CLI Starter — Design

## Overview

A starting-point template for building Node.js CLI applications distributed via npm.
It ships a single example command (`hello-world`) and a `-v`/`--version` flag, plus a
complete quality pipeline: formatting, linting, unit tests, and end-to-end tests.

The toolchain mirrors the existing `e2e-starter-projects` Electron baseline, with the
Electron/Playwright layer replaced by CLI-appropriate equivalents.

## Toolchain

| Concern           | Choice                                                       |
| ----------------- | ------------------------------------------------------------ |
| Package manager   | pnpm                                                         |
| Language / modules | TypeScript, ESM (`"type": "module"`)                        |
| CLI parsing       | commander                                                    |
| Build             | tsup (esbuild) → single `dist/cli.js` with shebang           |
| Unit tests        | Vitest (test functions directly)                             |
| E2E tests         | Vitest + execa (spawn the built binary)                      |
| Lint              | ESLint flat config (typescript-eslint, eslint-config-prettier) |
| Format            | Prettier (`singleQuote: true`, `printWidth: 100`)            |
| CI                | GitHub Actions, Node 22, runs full `pnpm test`               |
| Publishing        | Publishable structure (`bin`, `files`, `engines`); no release workflow |

## Architecture

Core principle: **`cli.ts` is a thin entry point; logic lives in testable functions.**
Each command parses nothing and contains no business logic — it delegates to a pure
function in `commands/`. Future commands follow the same split, which is what makes this
a reusable template.

```
src/
  cli.ts                  # shebang + commander setup; registers commands, parses argv
  commands/
    helloWorld.ts         # exports helloWorld(): string — pure, returns "Hello, world!"
tests/
  unit/
    helloWorld.test.ts    # imports helloWorld() directly, asserts return value
  e2e/
    cli.test.ts           # execa-spawns dist/cli.js, asserts stdout / exit code
```

### Components

- **`src/cli.ts`** — Begins with a `#!/usr/bin/env node` shebang. Configures commander:
  registers the `hello-world` command and `-v`/`--version`. The `hello-world` handler
  calls `helloWorld()` and writes the result to stdout. No business logic lives here.
- **`src/commands/helloWorld.ts`** — Exports `helloWorld(): string` returning
  `"Hello, world!"`. Pure, no I/O, trivially unit-testable.

### Data flow

```
argv → commander → command handler → command function → stdout
```

## Version handling

- The version is read from `package.json` and registered with commander via
  `.version(...)`, so `package.json` is the single source of truth — no duplicated
  version constant.
- tsup is configured so the built `dist/cli.js` can resolve the package version at
  runtime (e.g. importing the version from `package.json`).

## Error handling & exit codes

- **`hello-world` success:** prints `Hello, world!` to stdout, exit code `0`.
- **`-v` / `--version`:** commander prints the version, exit code `0`.
- **Unknown command or flag:** commander prints an error and usage text to stderr and
  exits with a non-zero code (built-in behavior).

## Testing

- **Unit (`pnpm test:unit`):** Vitest runs `tests/unit/**/*.test.ts`. The hello-world
  test imports `helloWorld()` and asserts it returns `"Hello, world!"`.
- **E2E (`pnpm test:e2e`):** A separate Vitest config targets `tests/e2e/**/*.test.ts`.
  The build runs first; each test uses execa to spawn the built binary
  (`node dist/cli.js hello-world`, `node dist/cli.js --version`) and asserts on stdout
  and exit code — exercising the actual shipped artifact.

## Verification pipeline

`pnpm test` runs the full gate in order (mirrors the reference project):

| Step | Command             | What it checks                          |
| ---- | ------------------- | --------------------------------------- |
| 1    | `pnpm typecheck`    | TypeScript types (no emit)              |
| 2    | `pnpm test:unit`    | Unit tests (Vitest)                     |
| 3    | `pnpm lint`         | ESLint rules                            |
| 4    | `pnpm build`        | Production build (tsup)                 |
| 5    | `pnpm test:e2e`     | End-to-end tests against the built CLI  |
| 6    | `pnpm format:check` | Prettier formatting                     |

## CI

GitHub Actions workflow (`.github/workflows/ci.yml`), triggered on every push and pull
request:

1. Checkout
2. Install pnpm (`pnpm/action-setup`)
3. Setup Node.js 22 with pnpm cache
4. `pnpm install --frozen-lockfile`
5. `pnpm test`

Single `ubuntu-latest` runner to match the reference baseline. No `xvfb` is needed (no
display). Expanding to an OS/Node matrix is a straightforward later change.

## Distribution / local install

The package is structured to be publishable to npm (`bin`, `files` allowlist, `engines`,
`repository`) but no release automation is included. To install the local build globally
for development:

1. `pnpm build`
2. `pnpm link --global` (symlink — rebuilds are picked up without re-linking), or
   `pnpm add -g .` (installs a snapshot of the current build)
3. Run the command (`hello-world`) from anywhere.

## Project naming placeholder

Wherever a project/application name is required, the template uses the literal
placeholder `new-application-name`. The README includes a "Starting a new app" section
listing every location this placeholder must be replaced:

- `package.json` → `name`
- `package.json` → `description`
- `package.json` → `repository` (and related URL fields, if present)
- Any reference to the name in `README.md` (e.g. the title)

(The example command name remains `hello-world` and is independent of the application
name; replacing it is left to the developer when they add real commands.)

## README contents

- Overview
- Prerequisites (Node.js, pnpm)
- Install (`pnpm install`)
- Verification pipeline table
- Local global install workflow (`pnpm build` → `pnpm link --global`)
- "Starting a new app" — where to replace `new-application-name`

## Out of scope (YAGNI)

- Publishing/release automation to the npm registry.
- Multiple example commands or subcommand groups.
- Cross-platform / multi-Node CI matrix.
- Configuration files, plugins, or interactive prompts.
