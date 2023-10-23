<template>
    <div class="userinfo">
        <div class="userinfo_item">
            <div class="userinfo_header">
                <img
                    :src="`http://${fileurl}/${props.userheader}`">
            </div>
            <p class="userinfo_name">{{ props.username }}</p>
            <div class="edit_selfinfo">个人信息</div>
        </div>
    </div>
</template>

<script setup lang="ts">
// import {fileurl} from '../../main'
const fileurl= "localhost:3006"
let props = defineProps({
    username: {
        type: String,
        default:""
    },
    userheader:{
        type: String,
        required: true,
        default:`http://${fileurl}/static/icon.png`
    }
})
</script>

<style scoped lang="less">
@import url('./userinfo.less');

</style>
