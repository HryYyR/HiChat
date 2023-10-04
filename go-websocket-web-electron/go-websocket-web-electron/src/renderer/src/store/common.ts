// 定义关于counter的store
import { defineStore } from "pinia"

const useCounter = defineStore("counter", {
    state: () => ({
        token: "",
        userdata: <Userdata>{}
    })
})

// 将useCounter函数导出
export default useCounter

export type Userdata = {
    CreatedTime: string
    Email: string
    ID: number
    LoginTime: string
    NikeName: string
    UserName: string
    GroupList: Array<GroupList>
}

export type GroupList = {
    GroupInfo: GroupInfo
    MessageList: Array<MessageListitem>
}
export type GroupInfo = {
    Avatar: string
    CreatedAt: string
    CreaterID: number
    CreaterName: string
    DeletedAt: string
    Grade: number
    GroupName: string
    ID: number
    UUID: string
    UpdatedAt: string
}
export type MessageListitem = {
    ID:number
    Context: any
    CreatedAt: string
    GroupID: number
    IsReply: boolean
    Msg: string
    MsgType: number
    ReplyUserID: number
    UserID: number
    UserName: string
    UserUUID: string
}
