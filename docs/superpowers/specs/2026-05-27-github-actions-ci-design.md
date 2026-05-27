---
name: github-actions-ci
description: GitHub Actions CI workflow to run pnpm test on every push and pull request
metadata:
  type: project
---

# GitHub Actions CI Design

## Goal

Run `pnpm test` automatically on every push and pull request, providing fast feedback on typecheck, unit tests, lint, build, Playwright E2E tests, and format check.

## Triggers

- `push` on all branches (`'**'`)
- `pull_request` — needed so PRs opened from forks run CI in the upstream repo (fork pushes do not trigger workflows there)

Duplicate runs on a PR branch (one from `push`, one from `pull_request`) are acceptable for a starter pack; the `concurrency` block below keeps the queue bounded.

## Workflow

**File:** `.github/workflows/ci.yml`

**Job:** `test` on `ubuntu-latest`, `timeout-minutes: 20`

**Workflow-level config:**

- `permissions: contents: read` — principle of least privilege; the job only needs to check out the repo
- `concurrency: { group: ci-${{ github.ref }}, cancel-in-progress: true }` — rapid pushes to the same branch cancel earlier runs

**Steps:**

1. `actions/checkout@v4` — check out the repo
2. `pnpm/action-setup@v4` — install pnpm (version comes from the `packageManager` field in `package.json`)
3. `actions/setup-node@v4` with Node.js 20 and `cache: pnpm` — installs Node and restores the pnpm store cache (ordering matters: pnpm must be on PATH before setup-node's cache hook runs)
4. `pnpm install --frozen-lockfile` — install dependencies; fail if `pnpm-lock.yaml` is out of date rather than silently regenerating it
5. `xvfb-run pnpm test` — run the full test suite inside a virtual display so Electron can open a window

## Why xvfb-run

The `pnpm test:e2e` step launches a real Electron window via Playwright. Linux runners have no display server by default. `xvfb-run` starts a temporary virtual framebuffer (Xvfb), sets `DISPLAY`, runs the command, then tears down — no extra setup or teardown steps required.

## What pnpm test covers

The `test` script in `package.json` runs in sequence:

1. `pnpm typecheck` — TypeScript type checking for main and renderer
2. `pnpm test:unit` — Vitest unit tests
3. `pnpm lint` — ESLint
4. `pnpm build` — TypeScript compile + Vite build (produces `dist/` for E2E)
5. `pnpm test:e2e` — Playwright E2E tests against the built Electron app
6. `pnpm format:check` — Prettier format check

No steps need to be split — the single `xvfb-run pnpm test` invocation covers all of these.

## Decisions and tradeoffs

**Fail-fast over full picture.** The `pnpm test` script chains steps with `&&`, so a typecheck failure hides later results (lint, format, etc.). Accepted because contributors are expected to re-run failures locally to get full output — CI just needs to flag that _something_ failed.

**pnpm version pinned via `packageManager`.** `package.json` declares `"packageManager": "pnpm@11.3.0"` so local and CI stay in sync. `pnpm/action-setup@v4` reads this field, so the workflow does not specify a version explicitly.

**No Playwright browser install step.** Tests use `_electron.launch()` against the built Electron binary, not chromium/firefox/webkit, so `npx playwright install` is not needed.

**Linux-only.** Cross-platform CI (macOS/Windows) is out of scope. Contributors on those platforms verify locally.

**No artifact upload on failure.** Playwright traces and screenshots stay on the runner. Adding artifact upload is a follow-up if E2E failures prove hard to debug from logs alone.

## Node.js Version

Node.js 20 (current LTS). The project does not specify an `engines` field, so LTS is a safe default.
