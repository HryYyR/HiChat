import { ElMessage } from "element-plus"
import {fileurl} from '../main'

export function nowtime() {
    let time = new Date()
    let year = time.getFullYear()
    let month: any = time.getMonth()
    let day: any = time.getDate()
    let hour: any = time.getHours()
    let minute: any = time.getMinutes()
    let second: any = time.getSeconds()
    if (month <= 9) {
        month = "0" + month
    }
    if (hour <= 9) {
        hour = "0" + hour
    }
    if (minute <= 9) {
        minute = "0" + minute
    }
    if (second <= 9) {
        second = "0" + second
    }
    let text = year + '-' + month + '-' + day + ' ' + hour + ':' + minute + ':' + second
    return text
}

export function tip(type: any, message: string) {
    ElMessage({
        "type": type,
        "message": message
    })
}

export function SendGroupResourceMsg(msg: string,MsgType :number,userdata:any,groupid:number):string{
    let  data = {
        UserID: userdata.ID,
        UserName: userdata.UserName,
        UserAvatar: userdata.Avatar == "" ? `http://${fileurl}/static/icon.png` : userdata.Avatar,
        GroupID: groupid,
        Msg: msg,
        MsgType: MsgType,
        IsReply: false,
        ReplyUserID: 0,
        Context: [],
        CreatedAt: new Date()
    }
    let strdata  = JSON.stringify(data)
    return  strdata
}