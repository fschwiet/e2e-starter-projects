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
