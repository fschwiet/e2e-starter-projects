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
