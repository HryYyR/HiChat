<template>
    <div class="userinfo">

        <div class="userinfo_item">
            <div class="userinfo_header">
                <img :src="`http://${fileurl}/${props.userdata.Avatar}`" @click="userdetaildialoghandleClose">
            </div>
            <p class="userinfo_name">{{ props.userdata.UserName }}</p>
        </div>

        <el-dialog :modal="false" v-model="data.UserDetailDialogVisible" title="个人信息" width="50%"
            :before-close="userdetaildialoghandleClose">

            <p><el-input v-model="props.userdata.NikeName" size="default" disabled>
                    <template #prepend>
                        昵称:
                    </template>
                </el-input>
            </p>

            <p><el-input v-model="props.userdata.Email" size="default" disabled>
                    <template #prepend>
                        邮箱:
                    </template>
                </el-input>
            </p>

            <p><el-input v-model="data.city" size="default" @change="">
                    <template #prepend>
                        城市:
                    </template>
                </el-input>
            </p>
            <p><el-input v-model="data.age" size="default" @change="">
                    <template #prepend>
                        年龄:
                    </template>

                </el-input>
            </p>
            <p><el-input v-model="props.userdata.CreatedTime" size="default" disabled>
                    <template #prepend>
                        注册时间:
                    </template>
                </el-input>
            </p>

            <p><el-input v-model="props.userdata.LoginTime" size="default" disabled>
                    <template #prepend>
                        登陆时间:
                    </template>

                </el-input>
            </p>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="edituserdata">
                        保存
                    </el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import {  fileurl } from '../../main'
import { edituserdataapi } from '../../API/api'
import { ElMessage } from 'element-plus';
let emit = defineEmits(['userdetaildialoghandleClose','edituserdata'])
let props = defineProps({
    userdata: {
        type: Object,
        required: true
    },

    UserDetailDialogVisible: {
        type: Boolean,
        default: false
    }
})

const data = reactive({
    UserDetailDialogVisible: false,
    age: 0,
    city: ""
})
const userdetaildialoghandleClose = () => {
    emit('userdetaildialoghandleClose')
    data.UserDetailDialogVisible = props.UserDetailDialogVisible
    data.city = props.userdata.City
    data.age = props.userdata.Age
}

const edituserdata = () => {
    console.log(data.city, data.age);
    edituserdataapi(data.age, data.city).then(res => {
        tip('success', res.data.msg)
        emit('edituserdata',data.age,data.city)
    }).catch((err) => {
        tip('error', err.response.data.msg)
    })
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
@import url('./userinfo.less');
</style>
