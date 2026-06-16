# NPM CLI Starter Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a reusable npm CLI starter template with a `hello-world` command, a `-v`/`--version` flag, and a complete formatting/linting/unit/e2e quality pipeline.

**Architecture:** A thin commander-based entry point (`src/cli.ts`) delegates to pure functions in `src/commands/`. tsup (esbuild) bundles the TypeScript ESM source into a single executable `dist/cli.js`. Unit tests import the command functions directly; e2e tests spawn the built binary via execa and assert on stdout/exit code.

**Tech Stack:** pnpm, TypeScript (ESM), commander, tsup, Vitest, execa, ESLint (flat config), Prettier, GitHub Actions.

## Global Constraints

- Package manager: **pnpm**. All scripts invoked as `pnpm <script>`.
- Module system: **ESM** — `package.json` has `"type": "module"`.
- Project/application name placeholder: the literal string **`new-application-name`** wherever a name is required.
- Example command name: **`hello-world`** (independent of the application name).
- CLI output for the command: exactly **`Hello, world!`**.
- Version flag must be lowercase **`-v`** / **`--version`** (override commander's default `-V`).
- Version is sourced from `package.json` — single source of truth, no duplicated constant. Initial version: **`0.0.1`**.
- Prettier config: `singleQuote: true`, `printWidth: 100`.
- ESLint: flat config using `@eslint/js` recommended + `typescript-eslint` recommended + `eslint-config-prettier`.
- CI: GitHub Actions, **Node 22**, single `ubuntu-latest` runner, runs `pnpm test`. No `xvfb`.
- Verification pipeline order (fail-fast): `format:check → lint → typecheck → test:unit → build → test:e2e`. The only hard constraint is `build` before `test:e2e`.
- `dist/` is the only published artifact (`files` allowlist) and is git-ignored.

---

### Task 1: Toolchain foundation + `helloWorld` command

Scaffolds the project (package.json, all config files, dependency install) and delivers the first unit-tested pure function. Setup is folded in because the unit test, lint, format, and typecheck commands all need it.

**Files:**
- Create: `package.json`
- Create: `tsconfig.json`
- Create: `eslint.config.mjs`
- Create: `.prettierrc`
- Create: `.prettierignore`
- Create: `.gitignore`
- Create: `vitest.config.ts`
- Create: `src/commands/helloWorld.ts`
- Test: `tests/unit/helloWorld.test.ts`

**Interfaces:**
- Consumes: nothing (first task).
- Produces:
  - `helloWorld(): string` exported from `src/commands/helloWorld.ts`, returning `'Hello, world!'`.
  - npm scripts: `format`, `format:check`, `lint`, `typecheck`, `test:unit`.

- [ ] **Step 1: Create `.gitignore`**

```
node_modules/
dist/
```

- [ ] **Step 2: Create `package.json`** (dependency versions are filled in by the install step below; leave the dependency blocks out for now)

```json
{
  "name": "new-application-name",
  "version": "0.0.1",
  "description": "A starter template for npm CLI applications",
  "type": "module",
  "engines": {
    "node": ">=18"
  },
  "packageManager": "pnpm@11.3.0",
  "repository": "https://github.com/your-org/new-application-name",
  "scripts": {
    "typecheck": "tsc --noEmit",
    "test:unit": "vitest run",
    "lint": "eslint .",
    "format": "prettier --write .",
    "format:check": "prettier --check ."
  }
}
```

- [ ] **Step 3: Install dependencies**

Run:
```bash
pnpm add commander
pnpm add -D typescript @types/node vitest eslint @eslint/js typescript-eslint eslint-config-prettier globals prettier
```
Expected: both commands succeed; `package.json` now has `dependencies` (commander) and `devDependencies` (the rest); `pnpm-lock.yaml` and `node_modules/` are created.

- [ ] **Step 4: Create `tsconfig.json`**

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "esModuleInterop": true,
    "strict": true,
    "skipLibCheck": true,
    "noEmit": true,
    "types": ["node"]
  },
  "include": ["src", "tests", "*.config.ts", "*.config.mjs"]
}
```

- [ ] **Step 5: Create `.prettierrc`**

```json
{
  "singleQuote": true,
  "printWidth": 100
}
```

- [ ] **Step 6: Create `.prettierignore`**

```
dist/
pnpm-lock.yaml
```

- [ ] **Step 7: Create `eslint.config.mjs`**

```js
import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import globals from 'globals';
import prettier from 'eslint-config-prettier';

export default tseslint.config(
  { ignores: ['dist/'] },
  eslint.configs.recommended,
  tseslint.configs.recommended,
  prettier,
  {
    files: ['src/**/*.ts', 'tests/**/*.ts', '*.config.ts', '*.config.mjs'],
    languageOptions: {
      globals: globals.node,
    },
  },
);
```

- [ ] **Step 8: Create `vitest.config.ts`**

```ts
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'node',
    include: ['tests/unit/**/*.test.ts'],
  },
});
```

- [ ] **Step 9: Write the failing unit test** — create `tests/unit/helloWorld.test.ts`

```ts
import { describe, it, expect } from 'vitest';
import { helloWorld } from '../../src/commands/helloWorld';

describe('helloWorld', () => {
  it('returns the greeting string', () => {
    expect(helloWorld()).toBe('Hello, world!');
  });
});
```

- [ ] **Step 10: Run the test to verify it fails**

Run: `pnpm test:unit`
Expected: FAIL — cannot resolve `../../src/commands/helloWorld` (module does not exist yet).

- [ ] **Step 11: Write the minimal implementation** — create `src/commands/helloWorld.ts`

```ts
export function helloWorld(): string {
  return 'Hello, world!';
}
```

- [ ] **Step 12: Run the test to verify it passes**

Run: `pnpm test:unit`
Expected: PASS — 1 test passing.

- [ ] **Step 13: Verify lint, format, and typecheck pass**

Run: `pnpm lint`
Expected: no errors (exit 0).

Run: `pnpm format:check`
Expected: all files formatted (exit 0). If it reports issues, run `pnpm format` then re-run `pnpm format:check`.

Run: `pnpm typecheck`
Expected: no errors (exit 0).

- [ ] **Step 14: Commit**

```bash
git add .gitignore package.json pnpm-lock.yaml tsconfig.json .prettierrc .prettierignore eslint.config.mjs vitest.config.ts src/commands/helloWorld.ts tests/unit/helloWorld.test.ts
git commit -m "feat: scaffold toolchain and hello-world command function"
```

---

### Task 2: CLI entry point, build, and e2e tests

Adds the commander entry point, the tsup build, and e2e tests that spawn the built binary. The `bin` field and build/e2e scripts are added here because the e2e deliverable needs them.

**Files:**
- Create: `src/cli.ts`
- Create: `tsup.config.ts`
- Create: `vitest.e2e.config.ts`
- Test: `tests/e2e/cli.test.ts`
- Modify: `package.json` (add `bin`, `files`, and `build` / `test:e2e` scripts)

**Interfaces:**
- Consumes: `helloWorld()` from `src/commands/helloWorld.ts` (Task 1).
- Produces:
  - Built executable at `dist/cli.js` (ESM, with `#!/usr/bin/env node` shebang).
  - `bin` mapping `hello-world` → `dist/cli.js`.
  - npm scripts: `build`, `test:e2e`.

- [ ] **Step 1: Install build and e2e dependencies**

Run:
```bash
pnpm add -D tsup execa
```
Expected: success; `tsup` and `execa` added to `devDependencies`.

- [ ] **Step 2: Create `tsup.config.ts`**

```ts
import { defineConfig } from 'tsup';

export default defineConfig({
  entry: { cli: 'src/cli.ts' },
  format: ['esm'],
  target: 'node18',
  clean: true,
});
```

(The `#!/usr/bin/env node` shebang lives at the top of `src/cli.ts`; esbuild preserves a leading shebang in the bundled output.)

- [ ] **Step 3: Create `src/cli.ts`**

```ts
#!/usr/bin/env node
import { Command } from 'commander';
import packageJson from '../package.json';
import { helloWorld } from './commands/helloWorld';

const program = new Command();

program
  .name('hello-world')
  .description('A starter template for npm CLI applications')
  .version(packageJson.version, '-v, --version', 'output the version number');

program
  .command('hello-world')
  .description('Print a greeting')
  .action(() => {
    console.log(helloWorld());
  });

program.parse();
```

- [ ] **Step 4: Add `bin`, `files`, and scripts to `package.json`**

Add the `bin` and `files` keys (top level) and the two new scripts. After this edit the relevant parts of `package.json` read:

```json
  "bin": {
    "hello-world": "dist/cli.js"
  },
  "files": ["dist"],
  "scripts": {
    "build": "tsup",
    "typecheck": "tsc --noEmit",
    "test:unit": "vitest run",
    "test:e2e": "vitest run --config vitest.e2e.config.ts",
    "lint": "eslint .",
    "format": "prettier --write .",
    "format:check": "prettier --check ."
  }
```

- [ ] **Step 5: Create `vitest.e2e.config.ts`**

```ts
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'node',
    include: ['tests/e2e/**/*.test.ts'],
    testTimeout: 30000,
  },
});
```

- [ ] **Step 6: Write the failing e2e test** — create `tests/e2e/cli.test.ts`

```ts
import { describe, it, expect } from 'vitest';
import { execa } from 'execa';
import { fileURLToPath } from 'node:url';

const cliPath = fileURLToPath(new URL('../../dist/cli.js', import.meta.url));

describe('hello-world CLI', () => {
  it('prints the greeting for the hello-world command', async () => {
    const { stdout, exitCode } = await execa('node', [cliPath, 'hello-world']);
    expect(stdout).toBe('Hello, world!');
    expect(exitCode).toBe(0);
  });

  it('prints the version with --version', async () => {
    const { stdout, exitCode } = await execa('node', [cliPath, '--version']);
    expect(stdout.trim()).toBe('0.0.1');
    expect(exitCode).toBe(0);
  });
});
```

- [ ] **Step 7: Run the e2e test to verify it fails**

Run: `pnpm test:e2e`
Expected: FAIL — `dist/cli.js` does not exist yet (execa cannot launch it).

- [ ] **Step 8: Build the CLI**

Run: `pnpm build`
Expected: success; `dist/cli.js` is created. Confirm the first line of `dist/cli.js` is `#!/usr/bin/env node`.

- [ ] **Step 9: Run the e2e test to verify it passes**

Run: `pnpm test:e2e`
Expected: PASS — 2 tests passing.

- [ ] **Step 10: Verify lint, format, and typecheck still pass**

Run: `pnpm lint`
Expected: exit 0.

Run: `pnpm format:check`
Expected: exit 0 (run `pnpm format` first if needed).

Run: `pnpm typecheck`
Expected: exit 0.

- [ ] **Step 11: Commit**

```bash
git add package.json pnpm-lock.yaml tsup.config.ts src/cli.ts vitest.e2e.config.ts tests/e2e/cli.test.ts
git commit -m "feat: add commander CLI, tsup build, and e2e tests"
```

---

### Task 3: Aggregate pipeline, CI workflow, and README

Wires the individual scripts into the single `pnpm test` gate, adds the GitHub Actions workflow, and documents the project including how to replace the name placeholder and install the local build globally.

**Files:**
- Modify: `package.json` (add aggregate `test` script)
- Create: `.github/workflows/ci.yml`
- Create: `README.md`

**Interfaces:**
- Consumes: all scripts from Tasks 1–2 (`format:check`, `lint`, `typecheck`, `test:unit`, `build`, `test:e2e`).
- Produces: aggregate `test` script; CI workflow; README.

- [ ] **Step 1: Add the aggregate `test` script to `package.json`**

Add `test` to the `scripts` block (fail-fast order). The `scripts` block becomes:

```json
  "scripts": {
    "build": "tsup",
    "typecheck": "tsc --noEmit",
    "test:unit": "vitest run",
    "test:e2e": "vitest run --config vitest.e2e.config.ts",
    "test": "pnpm format:check && pnpm lint && pnpm typecheck && pnpm test:unit && pnpm build && pnpm test:e2e",
    "lint": "eslint .",
    "format": "prettier --write .",
    "format:check": "prettier --check ."
  }
```

- [ ] **Step 2: Run the full pipeline to verify it passes**

Run: `pnpm test`
Expected: PASS — runs format:check, lint, typecheck, unit (1 test), build, e2e (2 tests), all green, exit 0.

- [ ] **Step 3: Create `.github/workflows/ci.yml`**

```yaml
name: CI

on:
  push:
    branches: ['**']
  pull_request:

permissions:
  contents: read

concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true

jobs:
  test:
    runs-on: ubuntu-latest
    timeout-minutes: 20
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Install pnpm
        uses: pnpm/action-setup@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: pnpm

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Run test suite
        run: pnpm test
```

- [ ] **Step 4: Create `README.md`**

````markdown
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
````

- [ ] **Step 5: Verify formatting of the new files**

Run: `pnpm format:check`
Expected: exit 0. If Prettier reports the new Markdown/YAML files, run `pnpm format` then re-run.

- [ ] **Step 6: Commit**

```bash
git add package.json .github/workflows/ci.yml README.md
git commit -m "feat: add aggregate test pipeline, CI workflow, and README"
```

---

## Self-Review Notes

- **Spec coverage:** Toolchain table → Task 1 (configs/deps) + Task 2 (tsup/execa). Architecture (`cli.ts` thin, `commands/` pure) → Tasks 1–2. Version handling (`-v/--version` from package.json) → Task 2 Step 3. Error handling/exit codes → exercised by e2e (Task 2 Step 6). Testing (unit + e2e) → Tasks 1–2. Verification pipeline + order → Task 3 Step 1. CI → Task 3 Step 3. Distribution/local install → README (Task 3 Step 4). Naming placeholder → README "Starting a new app" (Task 3 Step 4). All spec sections mapped.
- **Placeholder scan:** No TBD/TODO; all code and commands are complete.
- **Type consistency:** `helloWorld(): string` defined in Task 1, consumed in Task 2 with matching import path and signature. Version string `0.0.1` consistent between `package.json` (Task 1) and the e2e assertion (Task 2).
