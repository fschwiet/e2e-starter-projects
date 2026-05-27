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
