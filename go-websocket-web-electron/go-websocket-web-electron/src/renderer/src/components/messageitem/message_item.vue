<template>
    <div class="msg_item"
        :style="{ justifyContent: item.MsgType == 1 ? item.UserID == props.userdata.ID ? 'flex-end' : 'flex-start' : 'center' }">

        <!-- 左头像 -->
        <div class="msg_header" v-if="item.UserID != props.userdata.ID && item.MsgType == 1">
            <img :src="`http://${fileurl}/${item.UserAvatar}`" alt="">
        </div>

        <!-- 内容 -->
        <pre style="text-wrap: wrap;" v-if="item.MsgType == 1" class="msg_text"
            :style="{ alignItems: item.UserID == props.userdata.ID ? 'flex-end' : 'flex-start' }">
        <p>{{ item.UserName }}</p>
        <p 
        @contextmenu.prevent.stop="openMsgHandleMenu($event, item)"  
        class="msg_info" 
        :class="item.UserID == props.userdata.ID ? 'selfinfo' : ''" 
        v-text="item.Msg"></p>
    </pre>


        <!-- 右头像 -->
        <div class="msg_header" v-if="item.UserID == props.userdata.ID && item.MsgType == 1">
            <img :src="`http://${fileurl}/${item.UserAvatar}`" alt="">
        </div>

        <!-- 用户退出消息 -->
        <div class="msg_quit" v-if="item.MsgType == 201 || item.MsgType == 202">
            {{ item.Msg }}
        </div>

    </div>
</template>

<script setup lang="ts">
import { fileurl } from '../../main'
const emit = defineEmits(['openMsgHandleMenu'])
let props = defineProps({
    item: {
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
    emit('openMsgHandleMenu', e, item)
}

</script>

<style scoped lang="less">
@import url('./message_item.less');
</style> 