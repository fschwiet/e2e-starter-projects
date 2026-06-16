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
