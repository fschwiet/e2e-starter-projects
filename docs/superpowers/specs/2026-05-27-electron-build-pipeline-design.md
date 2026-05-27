# Electron App Build & Verification Pipeline

**Date:** 2026-05-27  
**Status:** Approved

---

## Overview

Set up a minimal, well-configured Electron desktop app using TypeScript, with a complete local build and verification pipeline. No UI framework — plain HTML/CSS/TS for the renderer. The project serves as a clean, correct starting point for future feature development.

**Stack:**

- Electron (main + renderer processes)
- TypeScript (both processes, separate configs)
- Vite (renderer bundler)
- `tsc` (main process compiler)
- ESLint + `@typescript-eslint`
- Prettier
- Vitest (unit tests)
- Playwright (E2E tests)
- pnpm (package manager)

---

## Project Structure

```
jsonloserbaby/
├── src/
│   ├── main/
│   │   ├── main.ts          # Entry point: creates BrowserWindow, loads renderer
│   │   └── preload.ts       # contextBridge preload (contextIsolation: true)
│   └── renderer/
│       ├── index.html
│       ├── index.ts
│       └── style.css
├── dist/                    # Build output (gitignored)
│   ├── main/
│   └── renderer/
├── tests/
│   ├── unit/
│   │   └── main.test.ts     # Smoke test: verifies Vitest setup
│   └── e2e/
│       └── app.test.ts      # Smoke test: launches app, checks window opens
├── tsconfig.main.json       # tsc config for main + preload
├── tsconfig.renderer.json   # tsc config for renderer (used by Vite)
├── vite.config.ts
├── vitest.config.ts
├── playwright.config.ts
├── .eslintrc.cjs
├── .prettierrc
├── .prettierignore
├── .gitignore
├── README.md
└── package.json
```

---

## Verification Pipeline Order

The README documents these commands in this order — fast/cheap checks first, slow/expensive last:

| Step | Command             | Purpose                                    |
| ---- | ------------------- | ------------------------------------------ |
| 1    | `pnpm typecheck`    | Type-check both processes without emitting |
| 2    | `pnpm test:unit`    | Run unit tests (Vitest)                    |
| 3    | `pnpm lint`         | ESLint across src + tests                  |
| 4    | `pnpm build`        | Production build (tsc + Vite)              |
| 5    | `pnpm test:e2e`     | E2E tests against built app (Playwright)   |
| 6    | `pnpm format:check` | Verify Prettier formatting                 |

The `pnpm test` script runs all of the above in sequence.

---

## npm Scripts

```json
{
  "scripts": {
    "dev:renderer": "vite",
    "dev:main": "tsc -p tsconfig.main.json",
    "build": "tsc -p tsconfig.main.json && vite build",
    "typecheck": "tsc -p tsconfig.main.json --noEmit && tsc -p tsconfig.renderer.json --noEmit",
    "test:unit": "vitest run",
    "test:e2e": "playwright test",
    "test": "pnpm typecheck && pnpm test:unit && pnpm lint && pnpm build && pnpm test:e2e && pnpm format:check",
    "lint": "eslint src tests",
    "format": "prettier --write .",
    "format:check": "prettier --check ."
  }
}
```

**Development workflow** (no single `dev` script — kept simple and explicit):

1. `pnpm dev:main` — compile main process once
2. `pnpm dev:renderer` — start Vite dev server (`http://localhost:5173`)
3. `npx electron .` — launch Electron

`main.ts` detects dev vs. production via the `VITE_DEV_SERVER_URL` environment variable (set by Vite during `dev:renderer`): if present, load from that URL; otherwise load from `dist/renderer/index.html`.

---

## TypeScript Configuration

### `tsconfig.main.json` (main process + preload)

- `module: "CommonJS"` — required by Electron's Node.js main process
- `target: "ES2020"`
- `outDir: "dist/main"`
- `strict: true`
- Includes: `src/main/**/*`

### `tsconfig.renderer.json` (renderer, consumed by Vite)

- `module: "ESNext"` — Vite handles bundling
- `moduleResolution: "bundler"` — matches Vite's resolution strategy
- `target: "ES2020"`
- `strict: true`
- No `outDir` — Vite owns output
- Includes: `src/renderer/**/*`

No root `tsconfig.json` — the two named configs are the authoritative, non-overlapping definitions. This prevents accidental mixing of compilation contexts.

---

## ESLint Configuration (`.eslintrc.cjs`)

- Parser: `@typescript-eslint/parser`
- Extends: `eslint:recommended`, `plugin:@typescript-eslint/recommended`, `prettier` (via `eslint-config-prettier` to prevent style conflicts)
- `overrides` for environment globals:
  - `src/main/**` and `tests/unit/**` → `env: { node: true }`
  - `src/renderer/**` → `env: { browser: true }`
  - `tests/e2e/**` → `env: { node: true }`

**devDependencies:**

- `eslint`
- `@typescript-eslint/parser`
- `@typescript-eslint/eslint-plugin`
- `eslint-config-prettier`

---

## Prettier Configuration (`.prettierrc`)

```json
{
  "singleQuote": true,
  "printWidth": 100
}
```

`.prettierignore` excludes `dist/` and `node_modules/`.

---

## Unit Testing (Vitest)

**`vitest.config.ts`:**

- `environment: "node"`
- `include: ["tests/unit/**/*.test.ts"]`

**Scope:**

- Tests main process business logic extracted into pure functions
- Does not test Electron APIs directly (those are integration-level, covered by Playwright)
- Does not test renderer (no framework = no component logic to unit test yet)

**Initial smoke test (`tests/unit/main.test.ts`):**

- Imports a trivial utility function from `src/main/`
- Asserts expected output — verifies Vitest is wired up correctly

---

## E2E Testing (Playwright)

**`playwright.config.ts`:**

- Uses `@playwright/test` with Electron launch
- Launches built app from `dist/main/main.js`
- E2E tests always run against the production build — `pnpm build` must precede `pnpm test:e2e`
- `testDir: "tests/e2e"`

**Initial smoke test (`tests/e2e/app.test.ts`):**

- Launches the Electron app
- Asserts the main window opens
- Asserts a known element is visible in the renderer
- Closes the app cleanly

---

## Security Baseline

The preload script (`src/main/preload.ts`) is included from day one with:

- `contextIsolation: true` (Electron default since v12, kept explicit)
- `nodeIntegration: false`
- `contextBridge` wired up and ready for IPC when needed

This establishes correct security posture before any features are added.

---

## README Structure

The README documents the project in this order:

1. **Prerequisites** — Node version, pnpm install
2. **Install** — `pnpm install`
3. **Development** — `pnpm dev`
4. **Verification pipeline** — each command in order with a one-line description

---

## Key devDependencies

| Package                            | Purpose                                                 |
| ---------------------------------- | ------------------------------------------------------- |
| `electron`                         | Desktop runtime                                         |
| `typescript`                       | TypeScript compiler                                     |
| `vite`                             | Renderer bundler                                        |
| `vitest`                           | Unit test runner                                        |
| `@playwright/test`                 | E2E test runner                                         |
| `eslint`                           | Linter                                                  |
| `@typescript-eslint/parser`        | TS-aware ESLint parsing                                 |
| `@typescript-eslint/eslint-plugin` | TS ESLint rules                                         |
| `eslint-config-prettier`           | Disables ESLint style rules that conflict with Prettier |
| `prettier`                         | Code formatter                                          |
