import { test, expect, _electron as electron } from '@playwright/test';
import path from 'path';

test('main window opens and shows app title', async () => {
  const electronApp = await electron.launch({
    args: [path.join(__dirname, '../../dist/main/main.js')],
  });

  const window = await electronApp.firstWindow();

  await expect(window.locator('h1#app-title')).toBeVisible();
  await expect(window).toHaveTitle('new-application-name');

  await electronApp.close();
});
