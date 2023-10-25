import { contextBridge, ipcRenderer } from 'electron'
import { electronAPI } from '@electron-toolkit/preload'

// Custom APIs for renderer
const api: api = {
  settitle: () => {
    ipcRenderer.send('settitle')
  },
  backtologin: () => {
    ipcRenderer.send('backtologin')
  },
  changWindowSize: () => {
    ipcRenderer.send('changWindowSize')
  },
  toMin: () => {
    ipcRenderer.send('toMin')
  },
  toMax: () => {
    ipcRenderer.send('toMax')
  },
  toClose: () => {
    ipcRenderer.send('toClose')
  },

}
type api = {
  settitle: Function
  backtologin: Function
  changWindowSize: Function
  toMin: Function
  toMax: Function
  toClose: Function
}

// Use `contextBridge` APIs to expose Electron APIs to
// renderer only if context isolation is enabled, otherwise
// just add to the DOM global.
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld('electron', electronAPI)
    contextBridge.exposeInMainWorld('api', api)
  } catch (error) {
    console.error(error)
  }
} else {
  // @ts-ignore (define in dts)
  window.electron = electronAPI
  // @ts-ignore (define in dts)
  window.api = api
}
