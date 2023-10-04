<template>
    <div class="container" v-loading="data.loginloading" v-show="!data.islogin">
        <!-- 登录 -->
        <div class="view_container login_container" :style="{ marginTop: data.logindata.offset ? '-500px' : '0px' }">
            <div class="headericon">
                <el-icon size="30" color="white">
                    <Plus />
                </el-icon>
            </div>
            <el-input class="view_input" v-model="data.logindata.username" placeholder="账号" size="large"
                clearable></el-input>
            <el-input type="password" class="view_input" v-model="data.logindata.password" show-password placeholder="密码"
                size="large" clearable></el-input>
            <div class="option">
                <div class="option_item">
                    <span>记住密码: <el-switch size="small" v-model="data.logindata.rememberpassword" />
                    </span>
                </div>
                <div class="option_item">
                    <span>自动登录: <el-switch size="small" v-model="data.logindata.autologin" /></span>
                </div>
            </div>
            <div class="btn" @click="login">登录</div>
        </div>

        <!-- 注册 -->
        <div class="view_container register_container">
            <el-input class="view_input" v-model="data.registerdata.username" placeholder="用户名(由字母和数字组成,不能低于6位)"
                size="large" clearable></el-input>
            <el-input class="view_input" type="password" show-password v-model="data.registerdata.password"
                placeholder="密码(不能低于6位)" size="large" clearable></el-input>
            <el-input class="view_input" type="password" show-password v-model="data.registerdata.checkpassword"
                placeholder="确认密码" size="large" clearable></el-input>
            <el-input class="view_input" v-model="data.registerdata.email" placeholder="邮箱" size="large"
                clearable></el-input>
            <div class="view_input" style="display: flex;">
                <el-input placeholder="邮箱验证码" v-model="data.registerdata.emailcode" size="large"
                    style="margin-right: 10px;"></el-input><el-button type="success" size="large" @click="sendemailCode"
                    :disabled="data.registerdata.sendemailbtnvisible">发送验证码</el-button>


            </div>
            <div class="btn" @click="register">注册</div>
        </div>

        <!-- 公共组件 -->
        <div class="changeview" @click="toregister">{{ !data.logindata.offset ? "去注册" : "去登录" }}<el-icon
                v-if="!data.logindata.offset">
                <ArrowDownBold />
            </el-icon>
            <el-icon v-else>
                <ArrowUpBold />
            </el-icon>
        </div>

    </div>

    <!-- 主面板 -->
    <div v-show="data.islogin" class="index">
        <div class="left_list">
            <div class="userinfo">
                <div class="userinfo_item">
                    <div class="userinfo_header">
                        <img :src="data.userdata.Avatar.length==0?`http://${fileurl}/static/icon.png`:data.userdata.Avatar" >
                    </div>
                    <p class="userinfo_name">{{ data.userdata.UserName }}</p>
                    <div class="edit_selfinfo">个人信息</div>
                </div>
            </div>


            <!-- 群列表 -->
            <div class="group_list">
                <!-- 群工具 -->
                <div class="group_tools">
                    <input type="text" v-model="data.searchgroupinput" placeholder="搜索">


                    <el-dropdown trigger="click">
                        <el-icon class="open_addgroup_dialog_btn">
                            <Plus />
                        </el-icon>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item>
                                    <div @click="data.addgroupdata.addGroupDialogVisible = true">
                                        添加群聊
                                    </div>
                                </el-dropdown-item>
                                <el-dropdown-item>
                                    <div @click="data.creategroupdata.createGroupDialogVisible = true">
                                        创建群聊
                                    </div>
                                </el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>

                </div>
                <div class="group_item" v-for="(item) in data.grouplist" :key="item.GroupInfo.ID"
                    @click="setcurrentgrouplist(item)" @contextmenu.prevent.stop="openeditgroupmenu($event, item)"
                    :class="{ checkgroup: JSON.stringify(data.currentgroupdata) == '{}' ? false : item.GroupInfo.ID == data.currentgroupdata.GroupInfo.ID }">
                    <div>
                        <div class="group_item_header">
                            <img :src="item.GroupInfo.Avatar || `http://${fileurl}/static/default_group_avatar.png`">
                        </div>
                    </div>
                    <div>
                        <p>{{ item.GroupInfo.GroupName }}</p>
                        <p class="group_lastmsg">
                            {{ item.MessageList.length != 0 ?
                                item.MessageList.at(-1)?.MsgType == 1 ?
                                    `${item.MessageList.at(-1)?.UserName}:
                                                        ${item.MessageList.at(-1)?.Msg}` :
                                    item.MessageList.at(-1)?.Msg : ''

                            }}
                        </p>
                    </div>
                </div>
            </div>

            <div @click="outlogin" class="outlogin"><el-icon>
                    <ArrowLeftBold />
                </el-icon>退出登录</div>
        </div>

        <div class="right_list" v-if="JSON.stringify(data.currentgroupdata) != '{}'">
            <div class="rightlist_option">
                {{ data.currentgroupdata.GroupInfo.GroupName || "" }} ({{ data.currentgroupdata.GroupInfo.MemberCount }})
            </div>

            <!-- 消息列表 -->
            <div class="rightlist_container" ref="msglist">
                <div class="msg_item"
                    :style="{ justifyContent: item.MsgType == 1 ? item.UserID == data.userdata.ID ? 'flex-end' : 'flex-start' : 'center' }"
                    v-for="item in data.currentgroupdata.MessageList" :key="item.ID">

                    <!-- 左头像 -->
                    <div class="msg_header" v-if="item.UserID != data.userdata.ID && item.MsgType == 1">
                    </div>

                    <!-- 内容 -->
                    <div v-if="item.MsgType == 1" class="msg_text"
                        :style="{ alignItems: item.UserID == data.userdata.ID ? 'flex-end' : 'flex-start' }">
                        <p>{{ item.UserName }}</p>
                        <p class="msg_info" :class="item.UserID == data.userdata.ID ? 'selfinfo' : ''">{{ item.Msg }}</p>
                    </div>

                    <!-- 右头像 -->
                    <div class="msg_header" v-if="item.UserID == data.userdata.ID && item.MsgType == 1">
                    </div>

                    <!-- 用户退出消息 -->
                    <div class="msg_quit" v-if="item.MsgType == 201">
                        {{ item.Msg }}
                    </div>

                </div>

                <!-- 消息未读 -->
                <div class="message_unread" @click="scrolltonew(0, true)" v-show="data.messageunreaddata.unreadnumber != 0">
                    <span>
                        {{ data.messageunreaddata.unreadnumber }}条未读
                    </span>
                    <el-icon>
                        <ArrowDown />
                    </el-icon>
                </div>

            </div>
            <div class="rightlist_input">
                <div class="input_tool"></div>
                <textarea cols="30" rows="10" v-model="data.input"></textarea>
                <div @click="send" class="sendbtn" :style="{ color: data.input ? 'white' : 'rgba(255,255,255,0.4)' }">
                    发送
                </div>
            </div>
        </div>


        <!-- 添加群聊对话框 -->
        <el-dialog v-model="data.addgroupdata.addGroupDialogVisible" title="添加群聊" width="40%">
            <div style="display: flex;">
                <el-input style="margin-right: 3px;" v-model="data.addgroupdata.addgroupinput" placeholder="支持模糊搜索"
                    size="default" clearable @change=""></el-input>
                <el-button type="primary" size="default" @click="searchgroup">搜索</el-button>
            </div>
            <div>
                <div v-for="item in data.addgroupdata.addgroupsearchlist" :key="item.ID">
                    <span>{{ item.GroupName }}</span><el-button type="primary" size="default"
                        @click="applyentergroup(item)">加入</el-button>
                </div>
            </div>
        </el-dialog>

        <!-- 创建群聊对话框 -->
        <el-dialog v-model="data.creategroupdata.createGroupDialogVisible" title="创建群聊" width="40%">
            <el-steps :space="180" :active="data.creategroupdata.createstep" simple>
                <el-step title="群名称" />
                <el-step title="群头像" />
            </el-steps>
            <div class="creategroup_dialog">
                <el-input v-show="data.creategroupdata.createstep == 1" class="creategroup_groupname_input"
                    v-model="data.creategroupdata.creategroupinput" placeholder="" size="default" clearable>
                    <template #prepend>群聊名称:</template>
                </el-input>
                <el-upload v-show="data.creategroupdata.createstep == 2" class="avatar-uploader"
                    :action="data.creategroupdata.headeruploadurl" :show-file-list="false"
                    :on-success="uploadcreategroupheaderSuccess">
                    <img v-if="data.creategroupdata.headerurl" :src="data.creategroupdata.headerurl" class="avatar" />
                    <el-icon v-else class="avatar-uploader-icon">
                        <Plus />
                    </el-icon>
                </el-upload>
            </div>
            <template #footer>
                <span class="dialog-footer">
                    <el-button v-show="data.creategroupdata.createstep != 1" type="primary"
                        @click="() => data.creategroupdata.createstep--">
                        上一步
                    </el-button>
                    <el-button v-show="data.creategroupdata.createstep != 2" type="primary"
                        @click="() => data.creategroupdata.createstep++">
                        下一步
                    </el-button>
                    <el-button
                        :disabled="data.creategroupdata.headerurl == '' || data.creategroupdata.creategroupinput == ''"
                        v-show="data.creategroupdata.createstep == 2" type="primary" @click="creategroup">
                        完成
                    </el-button>
                </span>
            </template>
        </el-dialog>

        <!-- 确定退出(解散)群聊对话框 -->
        <el-dialog v-model="data.quitgroupdata.quitGroupDialogVisible" :title="data.quitgroupdata.title" width="30%">
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="data.quitgroupdata.quitGroupDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="quitgroup">
                        确定
                    </el-button>
                </span>
            </template>
        </el-dialog>

    </div>
</template>

<script setup lang="ts">
// const { ipcRenderer } = require('electron')
import { onMounted, reactive, ref, watch } from 'vue';
import useCounter from './store/common'
import { ElMessage, UploadProps } from 'element-plus';
// import { GroupList,Group,GroupInfo,GroupinfoList,Userdata } from './models/models';
import {
    loginapi,
    registerapi,
    RefreshGroupListapi,
    searchGroupapi,
    joingroupapi,
    creategroupapi,
    exitgroupapi,
    emailcodeapi
} from './API/api'
import ContextMenu from '@imengyu/vue3-context-menu'

const win: any = window
let Store = useCounter()
const url = 'localhost:3004'
// const staticurl = 'http://localhost:3005'
const fileurl = 'localhost:3006'
const msglist: any = ref(null)

onMounted(() => {
    win.api.settitle()
    initListener()

    data.logindata.username = localStorage.getItem("username") || ""
    data.logindata.password = localStorage.getItem("password") || ""
    if (localStorage.getItem("autologin") == "1" && data.logindata.username && data.logindata.password) {
        setTimeout(() => {
            login()
        }, 1000)
    }
})

const data = reactive({
    grouplist: <Group>[],//群信息
    currentgroupdata: <GroupList>{},
    userdata: <Userdata>{},  //用户信息
    logindata: { //登录信息
        username: "niko",
        password: "Hyyyh1527",
        rememberpassword: localStorage.getItem("rememberpassword") == "1" ? true : false,
        autologin: localStorage.getItem("autologin") == "1" ? true : false,
        offset: false
    },
    registerdata: { //注册信息
        username: "",
        password: "",
        checkpassword: "",
        email: "",
        emailcode: "",
        sendemailbtnvisible: false
    },
    ws: {
        wsconn: <any>null,  //ws连接
    },
    islogin: false, //是否登录
    messageunreaddata: {
        unreadnumber: 0
    },
    input: "hello!",  //聊天对话框
    searchgroupinput: "", //搜索群输入框
    loginloading: false, //是否加载中
    addgroupdata: {
        addgroupinput: "",   //添加群输入框
        addgroupsearchlist: <GroupinfoList>[], //添加群搜索列表
        addGroupDialogVisible: false,  //是否展示添加群对话框
    },
    creategroupdata: {
        headeruploadurl: `http://${fileurl}/uploadfile`,
        creategroupinput: "",   //添加群输入框
        createGroupDialogVisible: false,  //是否展示添加群对话框
        headerurl: "",
        createstep: 1
    },
    quitgroupdata: {
        quitGroupDialogVisible: false,
        title: "退出群聊",
        targetgroupdata: <GroupList>{

        }
    }

})

watch(data.logindata, (newValue, _) => {
    localStorage.setItem("rememberpassword", newValue.rememberpassword ? "1" : "0")
    localStorage.setItem("autologin", newValue.autologin ? "1" : "0")
})


// 发送消息
const send = () => {
    let message = {
        UserID: data.userdata.ID,
        UserName: data.userdata.UserName,
        GroupID: data.currentgroupdata.GroupInfo.ID,
        Msg: data.input,
        MsgType: 1,
        IsReply: false,
        ReplyUserID: 0,
        Context: [],
        CreatedAt: new Date()
    }
    data.ws.wsconn.send(JSON.stringify(message))
    data.input = ""
    scrolltonew(300, true)
    console.log(msglist);

}

const setcurrentgrouplist = (group: GroupList) => {
    const setcurrentlistener = () => {
        const { scrollHeight, scrollTop, offsetHeight } = msglist.value
        if (scrollHeight - scrollTop - 3 * 83.6 < offsetHeight
        ) {
            data.messageunreaddata.unreadnumber = 0
        }
    }

    data.currentgroupdata = group
    data.messageunreaddata.unreadnumber = 0 //清空未读
    scrolltonew()
    setTimeout(() => {
        msglist.value.addEventListener("scroll", setcurrentlistener)
    }, 0);
}

// 登录
const login = () => {
    data.loginloading = true
    const { username, password } = data.logindata
    loginapi(username, password).then(res => {
        console.log(res);
        // 数据处理
        Store.token = res.data.token
        localStorage.setItem("token", res.data.token)
        if (res.data.userdata.GroupList == null) {
            res.data.userdata.GroupList = []
        }
        data.userdata = res.data.userdata
        data.grouplist = res.data.userdata.GroupList.map((group: GroupList) => {
            if (group.MessageList == null) {
                group.MessageList = []
            }
            return group
        })

        // 数据保存
        localStorage.setItem("username", data.logindata.username)
        if (data.logindata.rememberpassword) {
            localStorage.setItem("password", data.logindata.password)
        } else {
            localStorage.removeItem("password")
        }

        // 连接ws
        data.ws.wsconn = new WebSocket(`ws://${url}/ws?token=${localStorage.getItem("token")}`),

            data.ws.wsconn.onopen = function () {
                // console.log(evt);
                console.log("connect success!");
            }

        data.ws.wsconn.onclose = function () {
            // console.log(evt);
            console.log("connect close!");
        }
        // 接收消息
        data.ws.wsconn.onmessage = function (evt: any) {
            var msgstr = evt.data.split('\n');
            let msg = JSON.parse(msgstr)
            console.log("收到消息:", msg);
            handleMsg(msg)
        }

        data.ws.wsconn.onerror = function (evt: any) {
            console.log(evt);
        }

        setTimeout(() => {
            data.loginloading = false
            win.api.changWindowSize()
            data.islogin = true
        }, 1000);

        // 设置显示

    }).catch((err) => {
        console.log(err);

        setTimeout(() => {
            tip("error", "账号或密码错误!")
            data.loginloading = false
        }, 1000);
        return
    })
}

// 处理消息
const handleMsg = (msg: any) => {

    const DefaultMsg = () => {
        data.grouplist.forEach((group) => {
            console.log(group.GroupInfo.ID, msg.GroupID);

            if (group.GroupInfo.ID == msg.GroupID) {
                if (group.MessageList == null) { group.MessageList = [] }
                group.MessageList.push(msg)
            }
        })

        if (msglist.value == null) return
        const { scrollHeight, scrollTop, offsetHeight } = msglist.value
        // console.log(scrollHeight, scrollTop, offsetHeight);
        if (msg.GroupID == data.currentgroupdata.GroupInfo.ID &&
            msg.UserID != data.userdata.ID
        ) {
            if (scrollTop + offsetHeight + (3 * 83.6) > scrollHeight) {
                scrolltonew(0, true)
            } else {
                data.messageunreaddata.unreadnumber += 1
            }
        }
    }

    const refreshGroupMsg = () => {
        console.log("收到刷新消息");
        refreshgrouplist()
        return
    }

    const QuitGroupMsg = async (msg) => {

        if (JSON.stringify(data.currentgroupdata) != "{}") {
            if (msg.UserID == data.currentgroupdata.GroupInfo.CreaterID) {
                console.log("清空当前列表");
                data.currentgroupdata = <GroupList>{}
            }
        }
        await refreshgrouplist()
        // 添加一个退出群聊提示消息
        data.grouplist.forEach(group => {
            console.log(group.GroupInfo.ID, msg.GroupID);

            if (group.GroupInfo.ID == msg.GroupID) {
                let quitmsg: MessageListitem = {
                    MsgType: 201,
                    Msg: `${msg.UserName}退出了群聊`,
                    CreatedAt: msg.CreatedAt,
                    ID: 0,
                    Context: null,
                    GroupID: msg.GroupID,
                    IsReply: false,
                    ReplyUserID: 0,
                    UserID: msg.UserID,
                    UserName: msg.UserName,
                    UserUUID: "",
                }
                group.MessageList.push(quitmsg)
            }
        })
    }

    const typelist = {
        1: DefaultMsg,
        200: refreshGroupMsg,
        201: QuitGroupMsg
    }
    const msgtype = msg.MsgType
    typelist[msgtype](msg)
}

// 退出登录
const outlogin = () => {
    // ipcRenderer.send('backtologin')
    win.api.backtologin()
    setTimeout(() => {
        data.islogin = false
    }, 50);
    data.ws.wsconn.close()
    data.currentgroupdata = <GroupList>{}
    data.input = ""
    data.searchgroupinput = ""
    data.addgroupdata.addgroupinput = ""
    data.addgroupdata.addgroupsearchlist = <GroupinfoList>[]
}

// 去注册页面
const toregister = () => {
    data.logindata.offset = !data.logindata.offset
}

// 注册
const register = () => {
    let { username, password, checkpassword, email, emailcode } = data.registerdata

    var usernamereg = /^[a-zA-Z0-9_]+$/;
    if (username.length < 6 ||
        password.length < 6 ||
        password != checkpassword ||
        email.length == 0 ||
        emailcode.length == 0 ||
        !usernamereg.test(username) ||
        /\s/.test(username) ||
        /\s/.test(password)
    ) {
        tip("error", "信息有误,请检查后重试!")
        return
    }

    registerapi(username, password, email, emailcode).then(res => {
        console.log(res);
        if (res.status != 200) {
            tip("error", "注册失败,请稍后再试!")
            return
        }
        tip("success", "注册成功!")
        data.logindata.offset = !data.logindata.offset
        data.logindata.username = username
        data.logindata.password = password
        data.registerdata = {
            username: "",
            password: "",
            checkpassword: "",
            email: "",
            emailcode: "",
            sendemailbtnvisible: false
        }
    }).catch(err => {
        tip("error", err.response.data.msg)
    })
}

// 刷新列表
const refreshgrouplist = async () => {
    let datares = await RefreshGroupListapi(data.userdata.ID)
    if (datares.status != 200) {
        return alert("获取群列表失败!")
    }
    console.log("刷新群列表", datares.data);
    if (datares.data.usergrouplist != null) {
        data.grouplist = datares.data.usergrouplist.map((group: any) => {
            if (group.MessageList == null) {
                group.MessageList = []
            }
            return group
        })
    } else {
        data.grouplist = []
    }

    // 解决当有人退出然后重进后,发送消息丢失响应式的bug
    if (JSON.stringify(data.currentgroupdata) != "{}") {
        const currentid = data.currentgroupdata.GroupInfo.ID
        data.grouplist.forEach(group => {
            if (group.GroupInfo.ID == currentid) {
                data.currentgroupdata = group
            }
        })
    }


}

// 提示
function tip(type: any, message: string) {
    ElMessage({
        "type": type,
        "message": message
    })
}

// 搜索群聊
const searchgroup = () => {
    searchGroupapi(data.addgroupdata.addgroupinput).then((res) => {
        if (res.status != 200) {
            tip("Error", res.data.msg)
            return
        }
        data.addgroupdata.addgroupsearchlist = res.data.grouplist == null ? [] : res.data.grouplist

    }).catch(_ => {
        tip("Error", "发起请求失败！")
    })
}

// 加入群聊
const applyentergroup = async (group: GroupInfo) => {
    let res = await joingroupapi(group.GroupName)
    console.log(res);

    if (res.status != 200) {
        tip("Error", res.data.msg)
        return
    }
    tip("success", "加入成功！")
    data.addgroupdata.addGroupDialogVisible = false
    data.addgroupdata.addgroupinput = ""
    setTimeout(() => {
        refreshgrouplist()
    }, 500);
}

// 创建群聊
const creategroup = async () => {
    const { creategroupinput, headerurl } = data.creategroupdata
    let res = await creategroupapi(creategroupinput, headerurl)
    if (res.status != 200) {
        tip('error', res.data.msg)
        return
    }
    tip('success', res.data.msg)
    refreshgrouplist()
    data.creategroupdata = {
        headeruploadurl: `http://${fileurl}/uploadfile`,
        creategroupinput: "",
        createGroupDialogVisible: false,
        headerurl: "",
        createstep: 1
    }
}

// 退出群聊
const quitgroup = async () => {
    const GroupInfo = data.quitgroupdata.targetgroupdata.GroupInfo
    let res = await exitgroupapi(GroupInfo.ID)
    data.quitgroupdata.quitGroupDialogVisible = false
    if (res.status != 200) {
        tip("error", res.data.msg)
        return
    }
    data.currentgroupdata = <GroupList>{}
    tip("success", GroupInfo.CreaterID == data.userdata.ID ? "解散成功!" : "退出成功!")
}

// 打开右键菜单
const openeditgroupmenu = (e: any, item: GroupList) => {
    if (e.type == "contextmenu") {
        ContextMenu.showContextMenu({
            x: e.clientX,
            y: e.clientY,
            items: [
                {
                    label: item.GroupInfo.CreaterID == data.userdata.ID ? "解散群聊" : "退出群聊",
                    onClick: () => {
                        data.quitgroupdata.title = item.GroupInfo.CreaterID == data.userdata.ID ? "解散群聊" : "退出群聊"
                        data.quitgroupdata.targetgroupdata = item
                        data.quitgroupdata.quitGroupDialogVisible = true
                    }
                }
            ]
        });
    }

}

// 初始化键盘监听
const initListener = () => {
    window.addEventListener('keydown', (event) => {
        if (event.key == "Enter") {
            data.input += "\n"
            return
        }
    })
}

// 发送邮箱验证码
const sendemailCode = () => {
    const email = data.registerdata.email

    var reg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/;
    if (!reg.test(email)) {
        tip("error", "邮箱格式不正确!")
        return
    }

    data.registerdata.sendemailbtnvisible = true
    setTimeout(() => {
        data.registerdata.sendemailbtnvisible = false
    }, 60000);

    emailcodeapi(email).then(res => {
        console.log(res);
        tip("success", res.data.msg)

    }).catch(err => {
        tip("error", err.response.data.msg)
        console.log(err);

    })
}

//  上传群头像
const uploadcreategroupheaderSuccess: UploadProps['onSuccess'] = (response) => {
    console.log(response);

    if (response.code == 1) {
        tip("error", response.msg)
        return
    }
    tip("success", response.msg)
    data.creategroupdata.headerurl = `http://${fileurl}/${response.fileurl}`
}

// 滚动到最新
const scrolltonew = (delay: number = 0, smooth: boolean = false) => {
    data.messageunreaddata.unreadnumber = 0
    setTimeout(() => {
        msglist.value.scrollTo({ top: 100000, behavior: smooth ? "smooth" : "instant" })
    }, delay);
}

export type Userdata = {
    ID: number
    NikeName: string
    UserName: string
    Email: string
    CreatedTime: string
    LoginTime: string
    Avatar:string
    GroupList: Array<GroupList>
}

export type GroupList = {
    GroupInfo: GroupInfo
    MessageList: Array<MessageListitem>
}

export type Group = Array<GroupList>

export type GroupinfoList = Array<GroupInfo>

export type GroupInfo = {
    Avatar: string
    CreatedAt: string
    CreaterID: number
    CreaterName: string
    DeletedAt: string
    Grade: number
    MemberCount: number
    GroupName: string
    ID: number
    UUID: string
    UpdatedAt: string
}
export type MessageListitem = {
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


</script>
<style  lang="less">
@import url('./index.less');
</style>