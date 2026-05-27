import { contextBridge } from 'electron';

// Expose version info to the renderer via the context bridge.
// Add IPC channels here as the app grows.
contextBridge.exposeInMainWorld('versions', {
  node: () => process.versions.node,
  chrome: () => process.versions.chrome,
  electron: () => process.versions.electron,
});
