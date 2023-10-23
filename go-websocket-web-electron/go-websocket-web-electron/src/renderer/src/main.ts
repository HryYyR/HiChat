import { createApp } from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'
import axios from "axios"
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

// import router from './router/router'

import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css'
import ContextMenu from '@imengyu/vue3-context-menu'

export const fileurl = 'localhost:3006'
export const url = 'localhost:3004'

axios.defaults.baseURL = "http://localhost:3004"
axios.defaults.timeout = 4000;

// http request 请求拦截器
axios.interceptors.request.use(
    config => {
        // 这里判断localStorage里面是否存在token，如果有则在请求头里面设置
        if (localStorage.getItem("token")) {
            config.headers.Authorization = localStorage.getItem("token");
        }
        return config
    },
    err => {
        return err.response
    }
);

// // http response 响应拦截器
// axios.interceptors.response.use(
//     response => {
//         if (response.status != 200) {
//             return response.data.msg
//         }
//         return response
//     })


const app = createApp(App)
app.use(createPinia())
app.use(ElementPlus)
app.use(ContextMenu)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
// app.use(router)
app.mount('#app')
