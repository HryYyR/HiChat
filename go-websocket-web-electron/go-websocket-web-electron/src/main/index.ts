import { app, shell, BrowserWindow, ipcMain } from 'electron'
import { join } from 'path'
import { electronApp, optimizer, is } from '@electron-toolkit/utils'
import icon from '../../resources/icon.png?asset'



function createWindow(): void {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    width: 400,
    height: 550,
    show: false,
    // frame: false, //取消window自带的关闭最小化等
    resizable: false, //禁止改变主窗口尺寸
    autoHideMenuBar: true,
    title: "Hichat",
    ...(process.platform === 'linux' ? { icon } : {}),
    webPreferences: {
      contextIsolation: true,  //容器隔离
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false,
      nodeIntegration: true,
      devTools: true,
    }
  })

  function initMainWindow(mainWindow: BrowserWindow) {
    mainWindow.setResizable(true);
    mainWindow.setMaximizable(true);
    mainWindow.setTitle('Hichat');

    mainWindow.setMinimumSize(800, 500);
    mainWindow.setSize(1050, 700)

    mainWindow.center()
  }

  function initLoginWindow(mainWindow: BrowserWindow) {
    mainWindow.setResizable(false);
    mainWindow.setMaximizable(false);
    mainWindow.setTitle('登录');

    mainWindow.setMinimumSize(400, 550);
    mainWindow.setSize(400, 550)

    mainWindow.center()
  };

  function delayShowWindow(initFn, delay: number) {
    mainWindow.setOpacity(0);
    initFn(mainWindow);
    // 在最小化之后修改size会无效，所以要在最小化之前修改大小
    mainWindow.minimize();
    setTimeout(() => {
      mainWindow.setOpacity(1);
      mainWindow.show();
      mainWindow.focus();
    }, delay);
  }

  ipcMain.on('changWindowSize', (_, delay = 500) => {
    if (delay) {
      delayShowWindow(initMainWindow, delay);
    } else {
      initMainWindow(mainWindow);
    }

  })

  ipcMain.on('backtologin', (_, delay = 500) => {
    if (delay) {
      delayShowWindow(initLoginWindow, delay);
    } else {
      initLoginWindow(mainWindow);
    }
  })

  ipcMain.on('settitle', () => {
    mainWindow.setTitle("Hichat")
  })



  mainWindow.on('ready-to-show', () => {
    mainWindow.show()
  })




  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  // HMR for renderer base on electron-vite cli.
  // Load the remote URL for development or the local html file for production.
  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  // Set app user model id for windows
  electronApp.setAppUserModelId('com.electron')

  // Default open or close DevTools by F12 in development
  // and ignore CommandOrControl + R in production.
  // see https://github.com/alex8088/electron-toolkit/tree/master/packages/utils
  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  createWindow()

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// In this file you can include the rest of your app"s specific main process
// code. You can also put them in separate files and require them here.
