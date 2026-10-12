import { describe, it, expect } from 'vitest';

// Probability (0..1) that the "sometimes passes" test passes.
const PASS_PROBABILITY = 0.5;

describe('nondeterministic', () => {
  it('always passes', () => {
    expect(true).toBe(true);
  });

  it('always fails', () => {
    expect(true).toBe(false);
  });

  it('passes sometimes', () => {
    expect(Math.random()).toBeLessThan(PASS_PROBABILITY);
  });
});
