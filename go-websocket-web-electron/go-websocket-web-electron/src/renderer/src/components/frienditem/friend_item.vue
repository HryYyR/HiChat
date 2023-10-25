<template>
    <div class="list_item" @click="emit('setcurrentfriendlist',item)"
    >
        <div class="list_item_header">
            <img :src="`http://${fileurl}/${item.Avatar}`" alt="">
        </div>
        <div>
            <p>{{ props.item.NikeName }}</p>
            <!-- <p class="list_lastmsg" v-show="props.item.MessageList.length != 0">
                <span v-if="props.item.MessageList.at(-1)?.MsgType == 1">
                    {{ `${props.item.MessageList.at(-1)?.UserName}:${props.item.MessageList.at(-1)?.Msg}` }}
                </span>

                <span v-if="props.item.MessageList.at(-1)?.MsgType == 2">
                    {{ `${props.item.MessageList.at(-1)?.UserName}:[图片]` }}
                </span>

                <span v-if="props.item.MessageList.at(-1)?.MsgType == 3">
                    {{ `${props.item.MessageList.at(-1)?.UserName}:[语音]` }}
                </span>

            </p> -->
        </div>
        <div class="list_msginfo">
            <!-- <div>{{ props.item.MessageList.length != 0 ? props.item.MessageList.at(-1)?.CreatedAt.slice(11, 16) : '' }}
            </div>
            <div v-show="props.item.GroupInfo.UnreadMessage != 0">{{ props.item.MessageList.length != 0 ?
                props.item.GroupInfo.UnreadMessage : '' }}</div> -->
        </div>
    </div>
</template>

<script setup lang="ts">
import { Friend } from '../../models/models'
import { fileurl } from '../../main'
import { PropType } from 'vue'

const emit =defineEmits(['setcurrentfriendlist'])

const props = defineProps({
    item: {
        type: Object as PropType<Friend>,
        required: true
    }
})

</script>

<style scoped lang="less">
@import url('./friend_item.less');
@import url('../groupitem//groupitem.less');
</style>