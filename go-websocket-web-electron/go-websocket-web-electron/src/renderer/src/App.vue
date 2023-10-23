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
                    :disabled="data.registerdata.sendemailbtnvisible">{{ data.registerdata.sendcodebtn }}</el-button>


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
            <UserInfoVue :username="data.userdata.UserName" :userheader="data.userdata.Avatar" />

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

            <!-- 申请消息列表 -->
            <div class="apply_msg_list" v-if="data.islogin" @click="data.applymsgdata.applyMsgDialogVisible = true">
                <div class="apply_msg_list_left">
                    <p>消息通知</p><span v-show="data.userdata.ApplyList.filter(i => i.HandleStatus == 0).length != 0">{{
                        data.userdata.ApplyList.filter(i => i.HandleStatus == 0).length
                    }}</span>
                </div>
                <p class="apply_msg_list_right"><el-icon>
                        <ArrowRightBold />
                    </el-icon></p>
            </div>
            <!-- 群列表 -->
            <div class="group_list">
                <GroupItemVue v-for="(item) in data.grouplist" :key="item.GroupInfo.ID" :item="item"
                    :currentgroupdata="data.currentgroupdata" @setcurrentgrouplist="setcurrentgrouplist"
                    @openeditgroupmenu="openeditgroupmenu" />
            </div>

            <div @click="outlogin" class="outlogin">
                <el-icon>
                    <ArrowLeftBold />
                </el-icon>退出登录
            </div>
        </div>

        <div class="right_list" v-if="JSON.stringify(data.currentgroupdata) != '{}'">
            <div class="rightlist_option">
                {{ data.currentgroupdata.GroupInfo.GroupName || "" }} ({{ data.currentgroupdata.GroupInfo.MemberCount }})
                {{ data.currentgroupdata.MessageList.length }}
            </div>

            <!-- 消息列表 -->
            <div class="rightlist_container" ref="msglist">
                <MessageItemVue
                v-for="item in data.currentgroupdata.MessageList" :key="item.ID"
                :item="item"
                :userdata="data.userdata"
                :currentgroupdata="data.currentgroupdata"
                @openMsgHandleMenu="openMsgHandleMenu"
                />

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
                <div class="input_tool">

                    <el-icon class="tool_image" @click="sendimage">
                        <Picture />
                    </el-icon>

                </div>
                <textarea cols="30" rows="10" v-model="data.input"></textarea>
                <div @click="send" class="sendbtn" :style="{ color: data.input ? 'white' : 'rgba(255,255,255,0.4)' }">
                    发送
                </div>
            </div>
        </div>


        <!-- 申请加入群聊对话框 -->
        <el-dialog :before-close="beforeCloseAddGroupEvent" v-model="data.addgroupdata.addGroupDialogVisible" title="添加群聊"
            width="40%">
            <div style="display: flex;">
                <el-input style="margin-right: 3px;" v-model="data.addgroupdata.addgroupinput" placeholder="支持模糊搜索"
                    size="default" clearable @change=""></el-input>
                <el-button type="primary" size="default" @click="searchgroup">搜索</el-button>
            </div>
            <div>
                <div v-for="item in data.addgroupdata.addgroupsearchlist" :key="item.ID"
                    class="apply_join_geoup_dialog_grouplist_item">
                    <img :src="`http://${fileurl}/${item.Avatar}`" alt="">
                    <p>{{ item.GroupName }}</p>
                    <p class="group_number"><el-icon>
                            <User />
                        </el-icon>{{ item.MemberCount }}</p>
                    <el-button type="primary" size="default" @click="preapplyentergroup(item)">申请</el-button>

                </div>

            </div>
        </el-dialog>

        <!-- 填写申请加入群聊理由对话框 -->
        <el-dialog :before-close="beforeCloseAddGroupEvent" v-model="data.addgroupdata.preaddGroupDialogVisible"
            title="申请理由" width="40%">
            <el-input placeholder="申请理由" v-model="data.applyjoingroupdata.Msg" type="textarea" />
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="applyentergroup">
                        提交
                    </el-button>
                </span>
            </template>
        </el-dialog>

        <!-- 创建群聊对话框 -->
        <CreateGroupDialog :creategroupdata="data.creategroupdata" @changestep="changestep" @creategroup="creategroup"
            @uploadcreategroupheaderSuccess="uploadcreategroupheaderSuccess" />

        <!-- 确定退出(解散)群聊对话框 -->
        <el-dialog v-model="data.quitgroupdata.quitGroupDialogVisible" :title="data.quitgroupdata.title" width="30%">
            <p>退出后不会通知群聊中其他成员，且不会再接收此群消息</p>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="data.quitgroupdata.quitGroupDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="quitgroup">
                        确定
                    </el-button>
                </span>
            </template>
        </el-dialog>

        <!-- 消息通知对话框 -->
        <el-dialog style="background-color: rgb(229,229,229);" v-model="data.applymsgdata.applyMsgDialogVisible"
            title="消息通知" width="60%">

            <div class="apply_msg_list_dialog">
                <div v-for="item in data.userdata.ApplyList" :key="item.ID">
                    <p>
                    <p style="display: flex;">
                        <span> {{ item.ApplyUserName }} </span>
                    <p style="margin: 0 5px;">申请加入群聊</p> <span> {{ data.grouplist.filter(i => i.GroupInfo.ID ==
                        item.GroupID)[0].GroupInfo.GroupName
                    }}</span>
                    <p style="font-size: 0.8rem;line-height: 1.2rem;margin-left: 10px;color: rgb(168, 168, 168);">{{
                        item.CreatedAt.slice(11, 19) }}</p>
                    </p>
                    <p>
                        留言:{{ item.ApplyMsg }}
                    </p>
                    </p>
                    <div v-show="item.HandleStatus == 0">
                        <el-button type="primary" size="default" @click="handleapplymsg(item, 1)">同意</el-button>
                        <el-button type="danger" size="default" @click="handleapplymsg(item, -1)">拒绝</el-button>
                    </div>
                    <div v-show="item.HandleStatus == 1">已同意</div>
                    <div v-show="item.HandleStatus == -1">已拒绝</div>
                </div>
            </div>
        </el-dialog>

    </div>
</template>

<script setup lang="ts">
// const { ipcRenderer } = require('electron')
import {url,fileurl} from './main'
import { onMounted, reactive, ref, watch } from 'vue';
import useCounter from './store/common'
import { ElMessage, UploadProps } from 'element-plus';
import UserInfoVue from './components/userinfo/userinfo.vue'
import GroupItemVue from './components/groupitem/groupitem.vue'
import MessageItemVue from './components/messageitem/message_item.vue'
import CreateGroupDialog from './components/creategroupdialog/create_group_dialog.vue'
// import { GroupList,Group,GroupInfo,GroupinfoList,Userdata } from './models/models';
import {
    loginapi,
    registerapi,
    RefreshGroupListapi,
    searchGroupapi,
    joingroupapi,
    creategroupapi,
    exitgroupapi,
    emailcodeapi,
    applyjoingroupapi
} from './API/api'
import ContextMenu from '@imengyu/vue3-context-menu'

const win: any = window
let Store = useCounter()
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
        sendcodebtn: "发送验证码",
        sendemailbtnvisible: false
    },
    ws: {
        wsconn: <any>null,  //ws连接
    },
    islogin: false, //是否登录
    messageunreaddata: {
        unreadnumber: 0
    },
    input: "hello!",  //聊天输入框
    searchgroupinput: "", //搜索群输入框
    loginloading: false, //是否加载中
    addgroupdata: {
        addgroupinput: "",   //添加群输入框
        addgroupsearchlist: <GroupinfoList>[], //添加群搜索列表
        addGroupDialogVisible: false,  //是否展示添加群对话框
        preaddGroupDialogVisible: false,  //是否展示添加群理由对话框
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
    },
    applyjoingroupdata: {
        GroupName: "",
        GroupID: -1,
        Msg: "",
        ApplyWay: 1
    },
    applymsgdata: {
        applyMsgDialogVisible: false
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
        UserAvatar: data.userdata.Avatar == "" ? `http://${fileurl}/static/icon.png` : data.userdata.Avatar,
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
    scrolltonew(200, true)
    console.log(msglist);

}

// 设置选中,清除未读消息,监听滚动
const setcurrentgrouplist = (group: GroupList) => {
    const setcurrentlistener = () => {
        const { scrollHeight, scrollTop, offsetHeight } = msglist.value
        if (scrollHeight - scrollTop - 3 * 83.6 < offsetHeight
        ) {
            data.messageunreaddata.unreadnumber = 0
        }
    }


    data.currentgroupdata = group
    if (data.currentgroupdata.GroupInfo.UnreadMessage != 0) clearcurrentmsg()
    data.messageunreaddata.unreadnumber = 0 //清空未读
    scrolltonew()
    setTimeout(() => {
        msglist.value.addEventListener("scroll", setcurrentlistener)
    }, 0);

    group.GroupInfo.UnreadMessage = 0

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
            // console.log(group.GroupInfo.ID, msg.GroupID);
            if (group.GroupInfo.ID == msg.GroupID) {
                if (group.MessageList == null) { group.MessageList = [] }
                group.MessageList.push(msg)
            }
        })
    }

    const refreshGroupMsg = async () => {
        console.log("收到刷新消息");
        await refreshgrouplist()
        return
    }

    const QuitGroupMsg = async () => {
        if (JSON.stringify(data.currentgroupdata) != "{}") {
            if (msg.UserID == data.currentgroupdata.GroupInfo.CreaterID) {
                console.log("清空当前列表");
                data.currentgroupdata = <GroupList>{}
            }
        }
        await refreshgrouplist()
    }

    const JoginGroupMsg = async () => {
        await refreshgrouplist()
    }

    const typelist = {
        1: DefaultMsg,
        200: refreshGroupMsg,
        201: QuitGroupMsg,
        202: JoginGroupMsg
    }
    const msgtype = msg.MsgType
    typelist[msgtype](msg)

    if (msg.UserID != data.userdata.ID) {
        if (JSON.stringify(data.currentgroupdata) != "{}") {
            if (data.currentgroupdata.GroupInfo.ID != msg.GroupID) {
                data.grouplist.forEach(group => {
                    if (group.GroupInfo.ID == msg.GroupID) {
                        group.GroupInfo.UnreadMessage++
                    }
                })
            } else {
                clearcurrentmsg()
            }
        } else {
            data.grouplist.forEach(group => {
                if (group.GroupInfo.ID == msg.GroupID) {
                    group.GroupInfo.UnreadMessage++
                }
            })
        }

    }

    if (msglist.value == null || JSON.stringify(data.currentgroupdata) == '{}') return
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
            sendemailbtnvisible: false,
            sendcodebtn: "发送验证码"
        }
    }).catch(err => {
        tip("error", err.response.data.msg)
    })
}

const sendimage = () => {

}

// 刷新列表
const refreshgrouplist = async () => {
    let datares = await RefreshGroupListapi(data.userdata.ID)
    if (datares.status != 200) {
        return alert("获取群列表失败!")
    }
    console.log("刷新群列表", datares.data);

    data.userdata.ApplyList = datares.data.applylist == null ? [] : datares.data.applylist
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

const changestep = (i: number) => {
    i == 0 ? (data.creategroupdata.createstep--) : (data.creategroupdata.createstep++)
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
        console.log(res.data.grouplist);

        data.addgroupdata.addgroupsearchlist = res.data.grouplist == null ? [] : res.data.grouplist

    }).catch(_ => {
        tip("Error", "发起请求失败！")
    })
}

// 关闭添加群聊对话框之前
const beforeCloseAddGroupEvent = (done: any) => {
    data.addgroupdata.addgroupsearchlist = <GroupinfoList>{}
    data.addgroupdata.addgroupinput = ""
    done()
}

// 填写添加群聊理由前(绑定当前选择数据)
const preapplyentergroup = (group: GroupInfo) => {
    data.applyjoingroupdata.GroupID = group.ID
    data.applyjoingroupdata.GroupName = group.GroupName
    data.addgroupdata.preaddGroupDialogVisible = true
}

// 申请加入群聊
const applyentergroup = async () => {
    console.log(data.applyjoingroupdata);

    applyjoingroupapi(data.applyjoingroupdata).then(res => {
        console.log(res);
        tip("success", res.data.msg)
    }).catch(err => {
        tip("error", err.response.data.msg)
    })
    data.addgroupdata.addgroupinput = ""
    data.addgroupdata.preaddGroupDialogVisible = false
    data.addgroupdata.addGroupDialogVisible = false
    data.addgroupdata.addgroupsearchlist = <GroupinfoList>{}

}

// 创建群聊
const creategroup = async () => {
    const { creategroupinput, headerurl } = data.creategroupdata
    creategroupapi(creategroupinput, headerurl).then(res => {
        tip('success', res.data.msg)
        refreshgrouplist()
        data.creategroupdata = {
            headeruploadurl: `http://${fileurl}/uploadfile`,
            creategroupinput: "",
            createGroupDialogVisible: false,
            headerurl: "",
            createstep: 1
        }
    }).catch(err => {
        tip('error', err.response.data.msg)
    })
    // if (res.status != 200) {
    //     tip('error', res.data.msg)
    //     return
    // }
    // tip('success', res.data.msg)
    // refreshgrouplist()
    // data.creategroupdata = {
    //     headeruploadurl: `http://${fileurl}/uploadfile`,
    //     creategroupinput: "",
    //     createGroupDialogVisible: false,
    //     headerurl: "",
    //     createstep: 1
    // }
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

const handleapplymsg = (apply: ApplyItem, status: number) => {

    joingroupapi(apply.ID, status).then(res => {
        console.log(res.data);
        tip("success", res.data.msg)
        apply.HandleStatus = status

    }).catch(error => {
        console.log(error);
        tip("error", error)

    })
}

// 打开群聊右键菜单
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

// 打开消息右键菜单
const openMsgHandleMenu = (e: any, item: MessageListitem) => {
    if (e.type == "contextmenu") {
        ContextMenu.showContextMenu({
            x: e.clientX,
            y: e.clientY,
            items: [
                {
                    label: "复制",
                    onClick: () => {
                        let text = window.getSelection()?.toString() || ""
                        if (text.length == 0) {
                            navigator.clipboard.writeText(item.Msg)
                        } else {
                            navigator.clipboard.writeText(text)

                        }
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
    let i = 59
    data.registerdata.sendcodebtn = `60s`
    let downtime = setInterval(() => {
        data.registerdata.sendcodebtn = `${i}s`
        i--
    }, 1000)
    setTimeout(() => {
        clearInterval(downtime)
        data.registerdata.sendcodebtn = "发送验证码"
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
    data.creategroupdata.headerurl = response.fileurl
}

// 滚动到最新
const scrolltonew = (delay: number = 0, smooth: boolean = false) => {
    data.messageunreaddata.unreadnumber = 0
    setTimeout(() => {
        msglist.value.scrollTo({ top: 100000, behavior: smooth ? "smooth" : "instant" })
    }, delay);
}

const clearcurrentmsg = () => {
    let message = {
        UserID: data.userdata.ID,
        UserName: data.userdata.UserName,
        GroupID: data.currentgroupdata.GroupInfo.ID,
        MsgType: 401,
        CreatedAt: new Date()
    }
    data.ws.wsconn.send(JSON.stringify(message))
}

export type Userdata = {
    ID: number
    NikeName: string
    UserName: string
    Email: string
    CreatedTime: string
    LoginTime: string
    Avatar: string
    GroupList: Array<GroupList>
    ApplyList: Array<ApplyItem>
}
type ApplyItem = {
    ID: number
    GroupID: number
    ApplyMsg: string
    ApplyUserID: number
    ApplyUserName: string
    ApplyWay: number
    CreatedAt: string
    DeletedAt: string
    UpdatedAt: string
    HandleStatus: number

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
    UnreadMessage: number
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
    UserAvatar: string
}


</script>
<style  lang="less">
@import url('./index.less');
</style>