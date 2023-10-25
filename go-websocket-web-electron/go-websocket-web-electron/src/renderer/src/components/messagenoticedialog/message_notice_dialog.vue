<template>
    <!-- 消息通知对话框 -->
    <el-dialog style="background-color: rgb(229,229,229);" v-model="props.applymsgdata.applyMsgDialogVisible" title="消息通知"
        width="60%">

        <el-menu class="menu_group" default-active="0" :unique-opened="true" background-color="#e5e5e5" mode="horizontal" :ellipsis="false"
            @select="handleSelect">
            <el-menu-item class="menuitem" index="0">群聊通知 <span v-show="filterapplyjoingrouplist != 0" >{{filterapplyjoingrouplist}}</span></el-menu-item>
            <el-menu-item class="menuitem" index="1">用户通知 <span  v-show="filterapplyadduserlist != 0">{{filterapplyadduserlist}}</span></el-menu-item>
        </el-menu>


        <div v-show="data.selectmenu == 0" class="apply_msg_list_dialog">
            <div v-for="item in props.userdata.ApplyList" :key="item.ID">
                <p>
                <p style="display: flex;">
                    <span> {{ item.ApplyUserName }} </span>
                <p style="margin: 0 5px;">申请加入群聊</p> <span> {{ JSON.stringify(props.grouplist) != '[]'? props.grouplist.filter(i => i.GroupInfo.ID ==
                    item.GroupID)[0].GroupInfo.GroupName:'群聊不存在'
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


        <div v-show="data.selectmenu == 1" class="apply_msg_list_dialog">
            <div v-for="item in props.userdata.ApplyUserList" :key="item.ID">
                <p>
                <p style="display: flex;">
                    <span> {{ item.ApplyUserID == userdata.ID ? "你" : item.ApplyUserName }} </span>
                <p style="margin: 0 5px;">申请添加好友</p> <span v-show="item.PreApplyUserID != userdata.ID"> {{
                    item.PreApplyUserName }}</span>
                <p style="font-size: 0.8rem;line-height: 1.2rem;margin-left: 10px;color: rgb(168, 168, 168);">{{
                    item.CreatedAt.slice(11, 19) }}</p>
                </p>
                <p v-show="item.ApplyUserID != userdata.ID">
                    留言:{{ item.ApplyMsg }}
                </p>
                </p>
                <div v-show="item.HandleStatus == 0 && item.ApplyUserID != userdata.ID">
                    <el-button type="primary" size="default" @click="handleapplyaddusermsg(item, 1)">同意</el-button>
                    <el-button type="danger" size="default" @click="handleapplyaddusermsg(item, -1)">拒绝</el-button>
                </div>
                <div v-show="item.HandleStatus == 1">已同意</div>
                <div v-show="item.HandleStatus == -1">已拒绝</div>
                <div v-show="item.HandleStatus == 0 && item.ApplyUserID == userdata.ID">等待验证</div>
            </div>
        </div>


    </el-dialog>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import { GroupList } from '../../models/models'

const emit = defineEmits(['handleapplymsg','handleapplyaddusermsg'])
let props = defineProps({
    userdata: {
        type: Object,
        required: true,
    },
    grouplist: {
        type: Array<GroupList>,
        required: true,
    },
    applymsgdata: {
        type: Object,
        required: true,
    },
    filterapplyjoingrouplist: {
        type: Number,
        required: true,
    },
    filterapplyadduserlist: {
        type: Number,
        required: true,
    },
})

const data = reactive({
    selectmenu: 0
})

const handleapplymsg = (item, num:number) => {
    emit('handleapplymsg', item, num)
}

const handleapplyaddusermsg = (item, num:number) => {
    emit('handleapplyaddusermsg',item,num)
}

const handleSelect = (e) => {
    data.selectmenu = e
}                     
</script>

<style scoped lang="less">
@import url('./message_notice_dialog.less');
</style>