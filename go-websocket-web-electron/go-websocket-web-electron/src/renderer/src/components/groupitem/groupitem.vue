<template>
    <div class="list_item" @click="setcurrentgrouplist(props.item)"
        @contextmenu.prevent.stop="openeditgroupmenu($event, props.item)"
        :class="{ checkgroup: JSON.stringify(currentgroupdata) == '{}' ? false : props.item.GroupInfo.ID == currentgroupdata.GroupInfo.ID }">
        <div>
            <div class="list_item_header">
                <img
                    :src="`http://${fileurl}/${props.item.GroupInfo.Avatar}` || `http://${fileurl}/static/default_group_avatar.png`">
            </div>
        </div>
        <div>
            <p>{{ props.item.GroupInfo.GroupName }}</p>
            <p class="list_lastmsg" v-show="props.item.MessageList.length != 0">
                <span v-if="props.item.MessageList.at(-1)?.MsgType == 1">
                    {{ `${props.item.MessageList.at(-1)?.UserName}:${props.item.MessageList.at(-1)?.Msg}` }}
                </span>

                <span v-if="props.item.MessageList.at(-1)?.MsgType == 2">
                    {{ `${props.item.MessageList.at(-1)?.UserName}:[图片]` }}
                </span>

                <span v-if="props.item.MessageList.at(-1)?.MsgType == 3">
                    {{ `${props.item.MessageList.at(-1)?.UserName}:[语音]` }}
                </span>


                <!-- {{ props.item.MessageList.length != 0 ?
                    props.item.MessageList.at(-1)?.MsgType == 1 ?
                        `${props.item.MessageList.at(-1)?.UserName}: ${props.item.MessageList.at(-1)?.Msg}` :
                        props.item.MessageList.at(-1)?.MsgType ==
                            2 ? `${props.item.MessageList.at(-1)?.UserName}:[文件]` : props.item.MessageList.at(-1)?.Msg :
                    ''
                }} -->
            </p>
        </div>
        <div class="list_msginfo">
            <div>{{ props.item.MessageList.length != 0 ? props.item.MessageList.at(-1)?.CreatedAt.slice(11, 16) : '' }}
            </div>
            <div v-show="props.item.GroupInfo.UnreadMessage != 0">{{ props.item.MessageList.length != 0 ?
                props.item.GroupInfo.UnreadMessage : '' }}</div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { fileurl } from '../../main'


const emit = defineEmits(['setcurrentgrouplist', 'openeditgroupmenu'])

let props = defineProps({
    item: {
        type: Object,
        required: true
    },
    currentgroupdata: {
        type: Object,
        required: true

    }
})


const setcurrentgrouplist = (item) => {
    emit('setcurrentgrouplist', item)
}

const openeditgroupmenu = (e, item) => {
    emit('openeditgroupmenu', e, item)
}
</script>

<style scoped lang="less">
@import url('./groupitem.less');
</style>