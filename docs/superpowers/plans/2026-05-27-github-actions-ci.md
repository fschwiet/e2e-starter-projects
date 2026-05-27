# GitHub Actions CI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a GitHub Actions workflow that runs `pnpm test` on every push and pull request, exercising typecheck, unit tests, lint, build, Playwright E2E tests, and format check on an Ubuntu runner.

**Architecture:** A single workflow file at `.github/workflows/ci.yml` defines one `test` job on `ubuntu-latest`. The job checks out the repo, installs pnpm (version pinned via `packageManager` in `package.json`), installs Node 20 with pnpm-store caching, runs `pnpm install --frozen-lockfile`, and finally runs `xvfb-run pnpm test` so Electron can open a window under a virtual framebuffer. Concurrency is scoped per ref with `cancel-in-progress: true` to bound the queue.

**Tech Stack:** GitHub Actions, `actions/checkout@v4`, `pnpm/action-setup@v4`, `actions/setup-node@v4`, `xvfb-run`, pnpm 11.3.0, Node 20.

---

## File Map

| File                       | Purpose                                                              |
| -------------------------- | -------------------------------------------------------------------- |
| `.github/workflows/ci.yml` | GitHub Actions workflow that runs `pnpm test` on push and PR events. |

---

## Task 1: Create the CI Workflow

**Files:**

- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Create `.github/workflows/ci.yml`**

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
          node-version: 20
          cache: pnpm

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Run test suite
        run: xvfb-run pnpm test
```

Why each piece:

- `on.push.branches: ['**']` runs on every branch so feature branches get feedback before opening a PR.
- `on.pull_request` ensures PRs opened from forks run CI in the upstream repo (fork pushes do not trigger workflows there).
- `permissions: contents: read` follows least privilege — the job only needs to check out the repo.
- `concurrency` with `cancel-in-progress: true` cancels superseded runs on the same ref so rapid pushes don't queue.
- `pnpm/action-setup@v4` is placed **before** `actions/setup-node@v4` — `setup-node`'s `cache: pnpm` hook needs `pnpm` on `PATH` when it runs.
- `pnpm install --frozen-lockfile` fails CI if `pnpm-lock.yaml` is out of date rather than silently regenerating it.
- `xvfb-run` provides the virtual display Electron needs on a headless Linux runner; no Playwright browser install is required because the E2E suite uses `_electron.launch()` against the built binary.

- [ ] **Step 2: Validate YAML syntax locally**

Run:

```powershell
pnpm exec js-yaml .github/workflows/ci.yml > $null
```

Expected: command exits 0 with no output. If `js-yaml` is not available, fall back to:

```powershell
node -e "require('fs').readFileSync('.github/workflows/ci.yml','utf8'); console.log('ok')"
```

Expected: prints `ok`. (This only confirms the file exists and is readable; the push in Task 2 is the real validation.)

- [ ] **Step 3: Commit**

```powershell
git add .github/workflows/ci.yml
git commit -m "ci: add GitHub Actions workflow running pnpm test on push and PR"
```

---

## Task 2: Verify the Workflow Runs in GitHub Actions

**Files:** none (verification only)

- [ ] **Step 1: Push the branch**

```powershell
git push -u origin electron
```

Expected: push succeeds. If the branch already tracks a remote, `git push` is sufficient.

- [ ] **Step 2: Watch the workflow run**

Run:

```powershell
gh run watch
```

Expected: a run named `CI` appears for the latest commit on `electron`. The `test` job progresses through the steps in order: Checkout → Install pnpm → Setup Node.js → Install dependencies → Run test suite.

If `gh run watch` does not auto-select the latest run, list runs and watch the newest one explicitly:

```powershell
gh run list --workflow ci.yml --limit 1
gh run watch <run-id>
```

- [ ] **Step 3: Confirm the run succeeded**

Run:

```powershell
gh run list --workflow ci.yml --limit 1
```

Expected: the latest row shows `completed` and `success` for branch `electron`.

If the run failed, inspect logs with:

```powershell
gh run view --log-failed
```

Common failures and remedies:

- **`pnpm install --frozen-lockfile` fails** — `pnpm-lock.yaml` is stale. Run `pnpm install` locally, commit the lockfile, push again.
- **Electron fails to launch with `Missing X server or $DISPLAY`** — `xvfb-run` is not on the runner image. The `ubuntu-latest` image ships with it; if this still occurs, add `sudo apt-get update && sudo apt-get install -y xvfb` as a step before the test step.
- **`pnpm` not found during `setup-node` cache step** — confirm `pnpm/action-setup@v4` runs before `actions/setup-node@v4`.

- [ ] **Step 4: Confirm concurrency cancellation (optional spot check)**

Push a trivial follow-up commit (e.g., whitespace touch on the workflow file) and immediately push another. Expected: the first run shows status `cancelled` and only the latest run completes. This step is optional — skip if the first CI run already passed and the concurrency block is unchanged from the spec.

---

## Self-Review Notes

- Spec coverage: triggers (push/PR), `permissions: contents: read`, `concurrency` with cancel-in-progress, `ubuntu-latest`, `timeout-minutes: 20`, step ordering (pnpm before setup-node), `--frozen-lockfile`, `xvfb-run pnpm test`, Node 20, no Playwright browser install — all present in Task 1 Step 1.
- The spec's "Decisions and tradeoffs" section (fail-fast, no artifact upload, Linux-only, pnpm pinned via `packageManager`) requires no plan tasks — they are explanatory context, not work items.
- No placeholders, no TBDs, every command shown with expected output.
