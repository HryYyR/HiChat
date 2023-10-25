import axios from 'axios'
import { Md5 } from 'ts-md5';

// const staticurl = "http://hyyyh.top:3005"
// const fileurl = "http://hyyyh.top:3006"
// axios.defaults.baseURL = "http://hyyyh.top:3004"

const staticurl = "http://localhost:3005"
const fileurl = "http://localhost:3006"
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

    return axios({
        url: staticurl + "/register",
        method: "POST",
        data: registerinfo
    })
}

// 发送邮箱验证码
export function emailcodeapi(email: string) {
    const msg = {
        email: email
    }

    return axios({
        url: staticurl + "/emailcode",
        method: "POST",
        data: msg
    })
}

// 刷新获取用户的群聊列表
export function RefreshGroupListapi(id: number) {
    let userinfo = {
        ID: id
    }
    return axios({
        url: staticurl + "/user/getusergrouplist",
        method: "POST",
        data: userinfo
    })
}

// 刷新获取用户的好友列表
export function RefreshFriendListapi(id: number) {
    let userinfo = {
        ID: id
    }
    return axios({
        url: staticurl + "/user/getuserfriendlist",
        method: "POST",
        data: userinfo
    })
}

// 刷新获取用户的群聊通知列表
export function RefreshApplyJoinGroupListapi(id: number) {
    let userinfo = {
        ID: id
    }
    return axios({
        url: staticurl + "/user/getuserapplyjoingrouplist",
        method: "POST",
        data: userinfo
    })
}

// 刷新获取用户的好友申请列表
export function RefreshApplyAddFriendListapi(id: number) {
    let userinfo = {
        ID: id
    }
    return axios({
        url: staticurl + "/user/getuserapplyaddfriendlist",
        method: "POST",
        data: userinfo
    })
}



// 搜索群聊
export function searchGroupapi(text: string) {
    const msg = {
        searchstr: text
    }
    return axios.post("/user/searchGroup", msg)
}

// 处理加入群聊
export function joingroupapi(applyid: number, status: number) {
    let msg = {
        ApplyID: applyid,
        HandleStatus: status
    }
    return axios.post("/user/handlejoingroup", msg)
}

// 申请加入群聊
export function applyjoingroupapi(applydata: any) {
    return axios.post("/user/applyjoingroup", applydata)
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

// 修改用户信息
export function edituserdataapi(age: number, city: string) {
    let msg = {
        age: age,
        city: city
    }
    return axios({
        url: staticurl + "/user/edituserdata",
        method: "POST",
        data: msg
    })
}

// 获取用户信息
export function getuserdataapi(id: number) {
    let msg = {
        ID: id
    }
    return axios({
        url: staticurl + "/user/getuserdata",
        method: "POST",
        data: msg
    })
}

// 申请添加好友
export function applyadduserapi(PreApplyUserID:number, PreApplyUserName:string, ApplyUserID:number, ApplyUserName:string, ApplyMsg:string) {
    let msg = {
        ApplyUserID: ApplyUserID,  //申请人
        ApplyUserName: ApplyUserName,
        PreApplyUserID: PreApplyUserID, //被申请人
        PreApplyUserName: PreApplyUserName,
        ApplyMsg: ApplyMsg
    }
    return axios.post("/user/applyadduser", msg)
}

// 处理添加好友
export function adduserapi(applyid: number, status: number) {
    let msg = {
        ApplyID: applyid,
        HandleStatus: status
    }
    return axios.post("/user/handleadduser", msg)
}

// 上传资源
export function uploadresourceapi(file: FormData) {
    return axios({
        url: fileurl + "/uploadfile",
        method: "POST",
        data:file,
        headers:{
            "Content-Type":"multipart/form-data"
        }
    })
}