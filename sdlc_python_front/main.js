const { app, BrowserWindow } = require('electron');
const path = require('path');

function createWindow() {
  const win = new BrowserWindow({
    width: 1366, // 14寸屏幕的常见宽度
    height: 768, // 14寸屏幕的常见高度
    webPreferences: {
      nodeIntegration: true
    }
  });

    win.loadURL('http://localhost:4000'); // 开发时加载本地开发服务器

    // win.loadFile(path.join(__dirname, 'dist/index.html')); // 生产时加载打包后的文件

}

app.whenReady().then(() => {
  createWindow();

  app.on('activate', function () {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

app.on('window-all-closed', function () {
  if (process.platform !== 'darwin') app.quit();
});