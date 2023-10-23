import axios from "axios"
import { ElMessage } from 'element-plus';

// const baseUrl = "http://hyyyh.top:3004"
const baseUrl = "http://localhost.top:3004"

const instance = axios.create({
    // baseURL 将自动加在 url`前面，除非 url 是一个绝对 URL。
    // 它可以通过设置一个 baseURL 便于为 axios 实例的方法传递相对 URL
    baseURL: baseUrl,
    // timeout设置一个请求超时时间，如果请求时间超过了timeout，请求将被中断，单位为毫秒（ms）
    timeout: 60000,
    // headers是被发送的自定义请求头，请求头内容需要根据后端要求去设置，这里我们使用本项目请求头。
    headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    }
})

// http request 请求拦截器
instance.interceptors.request.use(
    config => {
        // 这里判断localStorage里面是否存在token，如果有则在请求头里面设置
        if (localStorage.getItem("token")) {
            config.headers.Authorization = localStorage.getItem("token");
        }
        return config
    },
    err => {
        return err
    }
)


// http response 响应拦截器
instance.interceptors.response.use(
    response => {
        console.log(response);
        if (response.status !=200){
            ElMessage({
                type: "error",
                message:response.data.msg
            })
        }
            return response
    })
