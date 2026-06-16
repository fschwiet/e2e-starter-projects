import { describe, it, expect } from 'vitest';
import { helloWorld } from '../../src/commands/helloWorld';

describe('helloWorld', () => {
  it('returns the greeting string', () => {
    expect(helloWorld()).toBe('Hello, world!');
  });
});
