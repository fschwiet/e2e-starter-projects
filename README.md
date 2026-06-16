# new-application-name

A starter template for building npm CLI applications with TypeScript, ESM, commander,
tsup, ESLint, Prettier, Vitest unit tests, and execa-driven end-to-end tests.

## Prerequisites

- [Node.js](https://nodejs.org/) v18+
- [pnpm](https://pnpm.io/) — install with `npm install -g pnpm`

## Install

```
pnpm install
```

## Verification Pipeline

Run these commands in order to verify a change is correct (fail-fast order):

| Step | Command             | What it checks                          |
| ---- | ------------------- | --------------------------------------- |
| 1    | `pnpm format:check` | Prettier formatting                     |
| 2    | `pnpm lint`         | ESLint rules                            |
| 3    | `pnpm typecheck`    | TypeScript types (no emit)              |
| 4    | `pnpm test:unit`    | Unit tests (Vitest)                     |
| 5    | `pnpm build`        | Production build (tsup → `dist/cli.js`) |
| 6    | `pnpm test:e2e`     | End-to-end tests against the built CLI  |

Run the full pipeline in one command:

```
pnpm test
```

## Install the local build globally

No npm publishing required — install your local build globally for development:

```
pnpm build
pnpm link --global
```

Then run the command from anywhere:

```
hello-world hello-world
```

`pnpm link --global` creates a symlink, so rebuilds are picked up without re-linking.
Alternatively, `pnpm add -g .` installs a snapshot of the current build.

## Starting a new app

This template uses the placeholder name `new-application-name`. Replace it in these
locations:

- `package.json` → `name`
- `package.json` → `description`
- `package.json` → `repository`
- `README.md` → the title and any references

The example command name (`hello-world`) is independent of the application name. Replace
it when you add your own commands (the `bin` key in `package.json`, `src/cli.ts`, and the
files under `src/commands/`).
