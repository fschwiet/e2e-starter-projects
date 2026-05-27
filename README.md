# Overview

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

1. Compile the src/main process then run src/renderer with vite:
   ```
   pnpm dev
   ```
2. Launch Electron:
   ```
   npx electron .
   ```

## Verification Pipeline

Run these commands in order to verify a change is correct:

| Step | Command             | What it checks                                  |
| ---- | ------------------- | ----------------------------------------------- |
| 1    | `pnpm typecheck`    | TypeScript types (main + renderer, no emit)     |
| 2    | `pnpm test:unit`    | Unit tests (Vitest)                             |
| 3    | `pnpm lint`         | ESLint rules                                    |
| 4    | `pnpm build`        | Production build (tsc + Vite)                   |
| 5    | `pnpm test:e2e`     | End-to-end tests against built app (Playwright) |
| 6    | `pnpm format:check` | Prettier formatting                             |

Run the full pipeline in one command:

```
pnpm test
```
