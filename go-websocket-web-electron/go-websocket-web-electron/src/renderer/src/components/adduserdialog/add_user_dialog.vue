<template>
    <el-dialog :before-close="beforeCloseAddUserEvent" v-model="data.addUserDialogVisible" title="添加好友" width="40%">
        <div class="add_user_dialog_container">
            <div>
                <img :src="`http://${fileurl}/${targetuserdata.Avatar}`" alt="">
            </div>
            <div class="user_info">
                <p> {{ targetuserdata.UserName }}</p>
                <p>
                    <span v-if="targetuserdata.City">{{ targetuserdata.City }}</span>
                    <span v-if="targetuserdata.Age != 0">{{ targetuserdata.Age }}岁</span>
                </p>
            </div>
        </div>
        <el-input placeholder="申请消息" v-model="data.applycause" type="textarea" />
        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="applyadduser">
                    提交
                </el-button>
            </span>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue';
import { fileurl } from '../../main'
import { applyadduserapi } from '../../API/api'
import { ElMessage } from 'element-plus';
// import { Userdata } from '../../App.vue'

const emit = defineEmits(['changeHeaderDialog'])
const data = reactive({
    applycause: "",
    addUserDialogVisible: false
})

let props = defineProps({
    targetuserdata: {
        type: Object,
        required: true
    },
    userdata: {
        type: Object,
        required: true
    },
    addUserDialogVisible: {
        type: Boolean,
        required: true
    }
})
watch(props, (_, nv) => {
    // console.log(nv);
    data.addUserDialogVisible = nv.addUserDialogVisible
}, { deep: true })
const beforeCloseAddUserEvent = () => {
    data.applycause=""
    emit('changeHeaderDialog', props.targetuserdata)
}

const applyadduser = () => {
    let PreUserId = props.targetuserdata.ID
    let PreUserName = props.targetuserdata.UserName
    let BackUserId = props.userdata.ID
    let BackUserName = props.userdata.UserName
    console.log(PreUserId, PreUserName, BackUserId, BackUserName);
    applyadduserapi(PreUserId, PreUserName, BackUserId, BackUserName,data.applycause).then(res => {
        console.log(res.data);
        tip('success', res.data.msg)
    }).catch(() => {
        tip('error', "申请失败!")
    })
    emit('changeHeaderDialog', props.targetuserdata)
}// 提示
function tip(type: any, message: string) {
    ElMessage({
        "type": type,
        "message": message
    })
}
</script>

<style scoped lang="less">
@import url('./add_user_dialog.less');
</style>