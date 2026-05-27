# Electron Build Pipeline Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Scaffold a minimal Electron app with TypeScript, ESLint, Prettier, Vitest unit tests, and Playwright E2E tests, with a documented local verification pipeline.

**Architecture:** Two independent TypeScript compilation targets — `tsc` for the main/preload process (CommonJS, Node), Vite for the renderer (ESM, browser). Testable main-process logic lives in a pure utility module (`windowConfig.ts`) so Vitest can run it in Node without the Electron runtime. Playwright launches the production-built app for E2E tests.

**Tech Stack:** Electron, TypeScript, Vite, Vitest, Playwright, ESLint (`@typescript-eslint`), Prettier, pnpm

---

## File Map

| File | Purpose |
|------|---------|
| `package.json` | Scripts, devDependencies, Electron `main` entry |
| `.gitignore` | Ignore `dist/` and `node_modules/` |
| `tsconfig.main.json` | tsc config for main process (CommonJS, `dist/main`) |
| `tsconfig.renderer.json` | TypeScript config consumed by Vite (ESNext/bundler) |
| `vite.config.ts` | Vite config — root `src/renderer`, output `dist/renderer` |
| `vitest.config.ts` | Vitest config — Node environment, `tests/unit/**` |
| `playwright.config.ts` | Playwright config — `tests/e2e/**`, no browsers |
| `.eslintrc.cjs` | ESLint with `@typescript-eslint`, env overrides per directory |
| `.prettierrc` | Prettier formatting options |
| `.prettierignore` | Exclude `dist/` and `node_modules/` from formatting |
| `src/main/windowConfig.ts` | Pure utility: returns `BrowserWindowConstructorOptions` — unit-testable |
| `src/main/main.ts` | Electron entry: creates window, loads renderer URL |
| `src/main/preload.ts` | Preload: exposes `versions` via `contextBridge` |
| `src/renderer/index.html` | HTML entry point with `#app-title` heading |
| `src/renderer/index.ts` | Renderer entry point (minimal) |
| `src/renderer/style.css` | Base styles |
| `tests/unit/main.test.ts` | Vitest smoke test for `getWindowOptions` |
| `tests/e2e/app.test.ts` | Playwright smoke test: window opens, heading visible |
| `README.md` | Prerequisites, install, dev workflow, verification pipeline |

---

## Task 1: Repository Foundation

**Files:**
- Create: `package.json`
- Create: `.gitignore`

- [ ] **Step 1: Create `.gitignore`**

```
node_modules/
dist/
```

- [ ] **Step 2: Create `package.json`**

```json
{
  "name": "jsonloserbaby",
  "version": "0.0.1",
  "description": "Minimal Electron app",
  "main": "dist/main/main.js",
  "scripts": {
    "dev:renderer": "vite",
    "dev:main": "tsc -p tsconfig.main.json",
    "build": "tsc -p tsconfig.main.json && vite build",
    "typecheck": "tsc -p tsconfig.main.json --noEmit && tsc -p tsconfig.renderer.json --noEmit",
    "test:unit": "vitest run",
    "test:e2e": "playwright test",
    "test": "pnpm typecheck && pnpm test:unit && pnpm lint && pnpm build && pnpm test:e2e && pnpm format:check",
    "lint": "eslint src tests --ext .ts",
    "format": "prettier --write .",
    "format:check": "prettier --check ."
  },
  "devDependencies": {}
}
```

- [ ] **Step 3: Install all dependencies**

Run:
```
pnpm add -D electron typescript vite vitest @playwright/test @types/node eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin eslint-config-prettier prettier
```

Expected: pnpm installs all packages, `node_modules/` is created, `pnpm-lock.yaml` is created.

- [ ] **Step 4: Commit**

```
git add package.json .gitignore pnpm-lock.yaml
git commit -m "chore: add package.json and install dependencies"
```

---

## Task 2: TypeScript Configuration

**Files:**
- Create: `tsconfig.main.json`
- Create: `tsconfig.renderer.json`

- [ ] **Step 1: Create `tsconfig.main.json`**

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "CommonJS",
    "moduleResolution": "node",
    "rootDir": "src/main",
    "outDir": "dist/main",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "sourceMap": true
  },
  "include": ["src/main/**/*"]
}
```

- [ ] **Step 2: Create `tsconfig.renderer.json`**

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"]
  },
  "include": ["src/renderer/**/*"]
}
```

- [ ] **Step 3: Verify configs parse correctly**

Run:
```
pnpm tsc -p tsconfig.main.json --noEmit --listFiles
```

Expected: tsc reports no source files found but exits without error (no syntax errors in config).

Run:
```
pnpm tsc -p tsconfig.renderer.json --noEmit --listFiles
```

Expected: same — exits cleanly.

- [ ] **Step 4: Commit**

```
git add tsconfig.main.json tsconfig.renderer.json
git commit -m "chore: add TypeScript configurations for main and renderer"
```

---

## Task 3: Vitest + Unit Smoke Test (TDD)

**Files:**
- Create: `vitest.config.ts`
- Create: `tests/unit/main.test.ts`
- Create: `src/main/windowConfig.ts`

- [ ] **Step 1: Create `vitest.config.ts`**

```typescript
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'node',
    include: ['tests/unit/**/*.test.ts'],
  },
});
```

- [ ] **Step 2: Write the failing unit test**

Create `tests/unit/main.test.ts`:

```typescript
import { describe, it, expect } from 'vitest';
import { getWindowOptions } from '../../src/main/windowConfig';

describe('getWindowOptions', () => {
  it('returns 800x600 dimensions', () => {
    const opts = getWindowOptions('/some/preload.js');
    expect(opts.width).toBe(800);
    expect(opts.height).toBe(600);
  });

  it('enables contextIsolation and disables nodeIntegration', () => {
    const opts = getWindowOptions('/some/preload.js');
    expect(opts.webPreferences?.contextIsolation).toBe(true);
    expect(opts.webPreferences?.nodeIntegration).toBe(false);
  });

  it('sets the preload path', () => {
    const opts = getWindowOptions('/some/preload.js');
    expect(opts.webPreferences?.preload).toBe('/some/preload.js');
  });
});
```

- [ ] **Step 3: Run the test — verify it fails**

Run:
```
pnpm test:unit
```

Expected: FAIL — `Error: Cannot find module '../../src/main/windowConfig'`

- [ ] **Step 4: Create `src/main/windowConfig.ts`**

```typescript
import type { BrowserWindowConstructorOptions } from 'electron';

export function getWindowOptions(preloadPath: string): BrowserWindowConstructorOptions {
  return {
    width: 800,
    height: 600,
    webPreferences: {
      preload: preloadPath,
      contextIsolation: true,
      nodeIntegration: false,
    },
  };
}
```

- [ ] **Step 5: Run the test — verify it passes**

Run:
```
pnpm test:unit
```

Expected: PASS — all 3 tests pass.

- [ ] **Step 6: Commit**

```
git add vitest.config.ts tests/unit/main.test.ts src/main/windowConfig.ts
git commit -m "test: add Vitest config and unit smoke test for getWindowOptions"
```

---

## Task 4: Main Process Source

**Files:**
- Create: `src/main/main.ts`
- Create: `src/main/preload.ts`

- [ ] **Step 1: Create `src/main/preload.ts`**

```typescript
import { contextBridge } from 'electron';

// Expose version info to the renderer via the context bridge.
// Add IPC channels here as the app grows.
contextBridge.exposeInMainWorld('versions', {
  node: () => process.versions.node,
  chrome: () => process.versions.chrome,
  electron: () => process.versions.electron,
});
```

- [ ] **Step 2: Create `src/main/main.ts`**

```typescript
import { app, BrowserWindow } from 'electron';
import path from 'path';
import { getWindowOptions } from './windowConfig';

function createWindow(): void {
  const preloadPath = path.join(__dirname, 'preload.js');
  const win = new BrowserWindow(getWindowOptions(preloadPath));

  const devServerUrl = process.env['VITE_DEV_SERVER_URL'];
  if (devServerUrl) {
    win.loadURL(devServerUrl);
  } else {
    win.loadFile(path.join(__dirname, '../renderer/index.html'));
  }
}

app.whenReady().then(() => {
  createWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});
```

- [ ] **Step 3: Type-check main process**

Run:
```
pnpm tsc -p tsconfig.main.json --noEmit
```

Expected: exits with no errors.

- [ ] **Step 4: Commit**

```
git add src/main/main.ts src/main/preload.ts
git commit -m "feat: add main process entry point and preload script"
```

---

## Task 5: Renderer Source + Vite Config

**Files:**
- Create: `vite.config.ts`
- Create: `src/renderer/index.html`
- Create: `src/renderer/index.ts`
- Create: `src/renderer/style.css`

- [ ] **Step 1: Create `vite.config.ts`**

```typescript
import { defineConfig } from 'vite';
import { resolve } from 'path';

export default defineConfig({
  root: resolve(__dirname, 'src/renderer'),
  build: {
    outDir: resolve(__dirname, 'dist/renderer'),
    emptyOutDir: true,
  },
});
```

- [ ] **Step 2: Create `src/renderer/index.html`**

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>jsonloserbaby</title>
    <link rel="stylesheet" href="./style.css" />
  </head>
  <body>
    <h1 id="app-title">jsonloserbaby</h1>
    <script type="module" src="./index.ts"></script>
  </body>
</html>
```

- [ ] **Step 3: Create `src/renderer/index.ts`**

```typescript
// Renderer entry point.
// window.versions (node, chrome, electron) is available via the preload script.
console.log('Renderer loaded');
```

- [ ] **Step 4: Create `src/renderer/style.css`**

```css
body {
  font-family: system-ui, sans-serif;
  margin: 2rem;
}
```

- [ ] **Step 5: Run the full build**

Run:
```
pnpm build
```

Expected:
- `dist/main/main.js`, `dist/main/preload.js`, `dist/main/windowConfig.js` are created
- `dist/renderer/index.html` and bundled assets are created
- No errors

- [ ] **Step 6: Commit**

```
git add vite.config.ts src/renderer/
git commit -m "feat: add renderer source and Vite config"
```

---

## Task 6: Playwright E2E Test

**Files:**
- Create: `playwright.config.ts`
- Create: `tests/e2e/app.test.ts`

- [ ] **Step 1: Create `playwright.config.ts`**

```typescript
import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: 'tests/e2e',
  timeout: 30000,
});
```

- [ ] **Step 2: Create `tests/e2e/app.test.ts`**

```typescript
import { test, expect, _electron as electron } from '@playwright/test';
import path from 'path';

test('main window opens and shows app title', async () => {
  const electronApp = await electron.launch({
    args: [path.join(__dirname, '../../dist/main/main.js')],
  });

  const window = await electronApp.firstWindow();

  await expect(window.locator('h1#app-title')).toBeVisible();
  await expect(window).toHaveTitle('jsonloserbaby');

  await electronApp.close();
});
```

- [ ] **Step 3: Run E2E tests**

Run:
```
pnpm test:e2e
```

Expected: PASS — 1 test passes. The Electron window opens, the heading is found, the title matches.

If the test fails with "dist/main/main.js not found", run `pnpm build` first (Task 5 Step 5).

- [ ] **Step 4: Commit**

```
git add playwright.config.ts tests/e2e/app.test.ts
git commit -m "test: add Playwright E2E config and app smoke test"
```

---

## Task 7: ESLint Configuration

**Files:**
- Create: `.eslintrc.cjs`

- [ ] **Step 1: Create `.eslintrc.cjs`**

```javascript
module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  plugins: ['@typescript-eslint'],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'prettier',
  ],
  overrides: [
    {
      files: ['src/main/**/*', 'tests/unit/**/*', 'tests/e2e/**/*'],
      env: { node: true },
    },
    {
      files: ['src/renderer/**/*'],
      env: { browser: true },
    },
    {
      files: ['*.cjs'],
      env: { node: true },
      rules: {
        '@typescript-eslint/no-require-imports': 'off',
      },
    },
  ],
};
```

- [ ] **Step 2: Run lint**

Run:
```
pnpm lint
```

Expected: no errors. If any errors are reported, fix them before proceeding. Common issues:
- `@typescript-eslint/no-unused-vars` — remove or prefix with `_` any unused variables
- `@typescript-eslint/no-explicit-any` — replace `any` with a proper type

- [ ] **Step 3: Commit**

```
git add .eslintrc.cjs
git commit -m "chore: add ESLint configuration"
```

---

## Task 8: Prettier Configuration

**Files:**
- Create: `.prettierrc`
- Create: `.prettierignore`

- [ ] **Step 1: Create `.prettierrc`**

```json
{
  "singleQuote": true,
  "printWidth": 100
}
```

- [ ] **Step 2: Create `.prettierignore`**

```
dist/
node_modules/
pnpm-lock.yaml
```

- [ ] **Step 3: Format all files**

Run:
```
pnpm format
```

Expected: Prettier rewrites any files that don't match the config. Review the diff (`git diff`) and ensure the changes look correct (quote style, line length).

- [ ] **Step 4: Verify format check passes**

Run:
```
pnpm format:check
```

Expected: exits with code 0, no files flagged.

- [ ] **Step 5: Commit**

```
git add .prettierrc .prettierignore
git add -u
git commit -m "chore: add Prettier configuration and format all files"
```

(`git add -u` stages any files that Prettier reformatted in Step 3.)

---

## Task 9: README

**Files:**
- Create: `README.md`

- [ ] **Step 1: Create `README.md`** with the following content:

````markdown
# jsonloserbaby

Minimal Electron app with TypeScript, ESLint, Prettier, Vitest unit tests, and Playwright E2E tests.

## Prerequisites

- [Node.js](https://nodejs.org/) v20+
- [pnpm](https://pnpm.io/) — install with `npm install -g pnpm`

## Install

```
pnpm install
```

## Development

Run these three steps in separate terminals (or sequentially):

1. Compile the main process:
   ```
   pnpm dev:main
   ```
2. Start the Vite dev server for the renderer:
   ```
   pnpm dev:renderer
   ```
3. Launch Electron:
   ```
   npx electron .
   ```

## Verification Pipeline

Run these commands in order to verify a change is correct:

| Step | Command             | What it checks                                    |
| ---- | ------------------- | ------------------------------------------------- |
| 1    | `pnpm typecheck`    | TypeScript types (main + renderer, no emit)       |
| 2    | `pnpm test:unit`    | Unit tests (Vitest)                               |
| 3    | `pnpm lint`         | ESLint rules                                      |
| 4    | `pnpm build`        | Production build (tsc + Vite)                     |
| 5    | `pnpm test:e2e`     | End-to-end tests against built app (Playwright)   |
| 6    | `pnpm format:check` | Prettier formatting                               |

Run the full pipeline in one command:

```
pnpm test
```
````

- [ ] **Step 2: Commit**

```
git add README.md
git commit -m "docs: add README with verification pipeline"
```

---

## Task 10: Full Pipeline Verification

- [ ] **Step 1: Run the complete verification pipeline**

Run:
```
pnpm test
```

Expected output (all steps pass):
```
> pnpm typecheck   ✓ (no errors)
> pnpm test:unit   ✓ 3 tests passed
> pnpm lint        ✓ (no errors)
> pnpm build       ✓ (dist/ populated)
> pnpm test:e2e    ✓ 1 test passed
> pnpm format:check ✓ (no files changed)
```

- [ ] **Step 2: If any step fails, fix it**

- **typecheck fails:** fix TypeScript errors in the reported file
- **test:unit fails:** check that `src/main/windowConfig.ts` exports `getWindowOptions` correctly
- **lint fails:** address each ESLint error; re-run `pnpm lint` to confirm clean
- **build fails:** check `tsc` and `vite build` output separately (`pnpm tsc -p tsconfig.main.json` then `pnpm vite build`)
- **test:e2e fails:** ensure `dist/main/main.js` exists; check that `index.html` contains `<h1 id="app-title">` and `<title>jsonloserbaby</title>`
- **format:check fails:** run `pnpm format`, then `git add -u`, then `pnpm format:check` again

- [ ] **Step 3: Commit the clean state**

```
git add -A
git commit -m "chore: verify full pipeline passes"
```
