import axios from 'axios'
import { Md5 } from 'ts-md5';

const staticurl = "http://localhost:3005"
axios.defaults.baseURL = "http://localhost:3004"

// 登录
export function loginapi(username: string, password: string) {
    const md5: Md5 = new Md5()
    md5.appendAsciiStr(password)
    const encryptionpassword = md5.end()
    let info = {
        username: username,
        password: encryptionpassword
    }
    return axios({
        url: staticurl + "/login",
        method: "POST",
        data: info
    })
}

// 注册
export function registerapi(username: string, password: string, email: string, emailcode: string) {
    const md5: Md5 = new Md5()
    md5.appendAsciiStr(password)
    const encryptionpassword = md5.end()

    let registerinfo = {
        username: username,
        password: encryptionpassword,
        email: email,
        code: emailcode
    }

    return axios.post("/register", registerinfo)
}

// 刷新列表
export function RefreshGroupListapi(id: number) {
    let userinfo = {
        ID: id
    }
    return axios.post("/user/RefreshGroupList", userinfo)
}

// 搜索群聊
export function searchGroupapi(text: string) {
    const msg = {
        searchstr: text
    }
    return axios.post("/user/searchGroup", msg)
}

// 加入群聊
export function joingroupapi(GroupName: string) {
    let msg = {
        GroupName: GroupName,
    }
    return axios.post("/user/joingroup", msg)
}

// 创建群聊
export function creategroupapi(creategroupinput: string, headerurl: string) {
    let msg = {
        GroupName: creategroupinput,
        Avatar: headerurl
    }
    return axios.post("/user/creategroup", msg)
}

// 退出群聊
export function exitgroupapi(id: number) {
    let msg = {
        ID: id
    }
    return axios.post("/user/exitgroup", msg)
}

// 发送邮箱验证码
export function emailcodeapi(email: string) {
    const msg = {
        email: email
    }
    return axios.post("/emailcode", msg)
}