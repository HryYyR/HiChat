<template>
    <el-dialog :before-close="closecreategroupdialog" v-model="props.creategroupdata.createGroupDialogVisible" title="创建群聊"
        width="40%">
        <el-steps :space="180" :active="props.creategroupdata.createstep" simple>
            <el-step title="群名称" />
            <el-step title="群头像" />
        </el-steps>
        <div class="creategroup_dialog">
            <el-input v-show="props.creategroupdata.createstep == 1" class="creategroup_groupname_input"
                v-model="props.creategroupdata.creategroupinput" placeholder="" size="default" clearable>
                <template #prepend>群聊名称:</template>
            </el-input>
            <el-upload v-show="props.creategroupdata.createstep == 2" class="avatar-uploader"
                :action="props.creategroupdata.headeruploadurl" :show-file-list="false"
                :on-success="uploadcreategroupheaderSuccess">
                <img v-if="props.creategroupdata.headerurl" :src="`http://${fileurl}/${props.creategroupdata.headerurl}`"
                    class="avatar" />
                <el-icon v-else class="avatar-uploader-icon">
                    <Plus />
                </el-icon>
            </el-upload>
        </div>
        <template #footer>
            <span class="dialog-footer">
                <el-button v-show="props.creategroupdata.createstep != 1" type="primary" @click="changestep(0)">
                    上一步
                </el-button>
                <el-button v-show="props.creategroupdata.createstep != 2" type="primary" @click="changestep(1)">
                    下一步
                </el-button>
                <el-button :disabled="props.creategroupdata.headerurl == '' || props.creategroupdata.creategroupinput == ''"
                    v-show="props.creategroupdata.createstep == 2" type="primary" @click="creategroup">
                    完成
                </el-button>
            </span>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import { fileurl } from '../../main'
const emit = defineEmits(['creategroup', 'uploadcreategroupheaderSuccess', 'changestep', 'closecreategroupdialog'])

const props = defineProps({
    creategroupdata: {
        type: Object,
        required: true
    }
})

const uploadcreategroupheaderSuccess = (response) => {
    emit('uploadcreategroupheaderSuccess', response)
}
const creategroup = () => {
    emit('creategroup')
}

const changestep = (i: number) => {
    emit('changestep', i)
}

const closecreategroupdialog = () => {
    emit('closecreategroupdialog')
}
closecreategroupdialog
</script>

<style scoped lang="less">
@import url('./create_group_dialog.less');
</style>