<template>
    <div class="msg_item"
        :style="{ justifyContent: (item.MsgType == 1 || item.MsgType == 2|| item.MsgType==3) ? item.UserID == props.userdata.ID ? 'flex-end' : 'flex-start' : 'center' }">
        <MessageTimeVue :time="item.CreatedAt" :pretime="preitem.CreatedAt" />
        <!-- 左头像 -->
        <div class="msg_header" @contextmenu.prevent.stop="openHeaderHandleMenu($event, item)"
            v-if="item.UserID != props.userdata.ID && (item.MsgType == 1 || item.MsgType == 2|| item.MsgType==3)">
            <img :src="`http://${fileurl}/${item.UserAvatar}`" alt="">
        </div>

        <!-- 内容 -->
        <pre style="text-wrap: wrap;" v-if="item.MsgType == 1 || item.MsgType == 2 || item.MsgType==3 " class="msg_text"
            :style="{ alignItems: item.UserID == props.userdata.ID ? 'flex-end' : 'flex-start' }">
        <p>{{ item.UserName }}</p>

        <el-image
        class="msg_info" 
        :src="`http://${fileurl}/${item.Msg}`"
        :zoom-rate="1.1"
        :max-scale="1"
        :min-scale="0.2"
        v-if="item.MsgType == 2"
        :class="item.UserID == props.userdata.ID ? 'selfinfo' : ''" 
        :preview-src-list="[`http://${fileurl}/${item.Msg}`]"
        fit="cover"
        hide-on-click-modal
        />

        <audio v-if="item.MsgType == 3" :src="`http://${fileurl}/${item.Msg}`" controls ></audio>

        <p 
        v-if="item.MsgType == 1"
        @contextmenu.prevent.stop="openMsgHandleMenu($event, item)"  
        class="msg_info" 
        :class="item.UserID == props.userdata.ID ? 'selfinfo' : ''" 
        v-text="item.Msg"></p>
        </pre>



        <!-- 右头像 -->
        <div class="msg_header" v-if="item.UserID == props.userdata.ID && (item.MsgType == 1 || item.MsgType == 2 || item.MsgType == 3)"
            @contextmenu.prevent.stop="openHeaderHandleMenu($event, item)">
            <img :src="`http://${fileurl}/${item.UserAvatar}`" alt="">
        </div>

        <!-- 用户退出消息 -->
        <div class="msg_quit" v-if="item.MsgType == 201 || item.MsgType == 202">
            {{ item.Msg }}
        </div>

    </div>
</template>

<script setup lang="ts">
import ContextMenu from '@imengyu/vue3-context-menu'
import { getuserdataapi } from '../../API/api';
import { fileurl } from '../../main'
import MessageTimeVue from './messagetime/message_time.vue'
import { ElMessage } from 'element-plus';
const emit = defineEmits(['openMsgHandleMenu', 'changeHeaderDialog'])
let props = defineProps({
    item: {
        type: Object,
        required: true,
    },
    preitem: {
        type: Object,
        required: true,

    },
    currentgroupdata: {
        type: Object,
        required: true,
    },
    userdata: {
        type: Object,
        required: true,
    }
})

const openMsgHandleMenu = (e, item) => {
    // emit('openMsgHandleMenu', e, item)
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
const openHeaderHandleMenu = (e, item) => {

    if (e.type == "contextmenu") {
        if (item.UserID == props.userdata.ID) return
        ContextMenu.showContextMenu({
            x: e.clientX,
            y: e.clientY,
            items: [
                {
                    label: "添加好友",
                    onClick: () => {
                        console.log("添加好友");
                        getuserdataapi(item.UserID).then(res => {
                            console.log(res.data.data);

                            emit('changeHeaderDialog', res.data.data)
                        }).catch(err => {
                            tip('error', err.response.Msg)
                        })
                    },
                    hidden:props.userdata.FriendList ?props.userdata.FriendList.filter(i=>i.Id == item.UserID).length !=0:false
                    
                },
                {
                    label: "查看资料",
                    onClick: () => {
                        console.log(item);
                    }
                },
            ]
        });
    }

}

// 提示
function tip(type: any, message: string) {
    ElMessage({
        "type": type,
        "message": message
    })
}

</script>

<style scoped lang="less">
@import url('./message_item.less');
</style> 