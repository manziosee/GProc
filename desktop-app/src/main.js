const { app, BrowserWindow, Menu, ipcMain } = require('electron');
const path = require('path');
const axios = require('axios');

class GProcDesktop {
  constructor() {
    this.mainWindow = null;
    this.gprocEndpoint = 'http://localhost:8080';
  }

  createWindow() {
    this.mainWindow = new BrowserWindow({
      width: 1200,
      height: 800,
      webPreferences: {
        nodeIntegration: true,
        contextIsolation: false
      },
      title: 'GProc Desktop - Local Process Manager'
    });

    this.mainWindow.loadFile(path.join(__dirname, 'renderer/index.html'));
    this.setupIPC();
  }

  setupIPC() {
    ipcMain.handle('get-processes', async () => {
      try {
        const response = await axios.get(`${this.gprocEndpoint}/api/v1/processes`);
        return response.data;
      } catch (error) {
        return { error: error.message };
      }
    });

    ipcMain.handle('start-process', async (event, name, command, options) => {
      try {
        const response = await axios.post(`${this.gprocEndpoint}/api/v1/processes`, {
          name, command, ...options
        });
        return response.data;
      } catch (error) {
        return { error: error.message };
      }
    });
  }
}

const gprocApp = new GProcDesktop();

app.whenReady().then(() => {
  gprocApp.createWindow();
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});