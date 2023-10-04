 type Userdata = {
    CreatedTime: string
    Email: string
    ID: number
    LoginTime: string
    NikeName: string
    UserName: string
    GroupList: Array<GroupList>
}

 type GroupList = {
    GroupInfo: GroupInfo
    MessageList: Array<MessageListitem>
}

 type Group = Array<GroupList>

 type GroupinfoList = Array<GroupInfo>

 type GroupInfo = {
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

 type MessageListitem = {
    ID: number
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

export type {MessageListitem,Userdata,Group,GroupInfo,GroupList,GroupinfoList}