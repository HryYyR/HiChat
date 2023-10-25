<template>
    <!-- 申请加入群聊对话框 -->
    <el-dialog @close="beforeCloseAddGroupEvent" v-model="props.addgroupdata.addGroupDialogVisible" title="添加群聊"
        width="40%">
        <div style="display: flex;">
            <el-input style="margin-right: 3px;" v-model="props.addgroupdata.addgroupinput" placeholder="支持模糊搜索"
                size="default" clearable @change=""></el-input>
            <el-button type="primary" size="default" @click="searchgroup">搜索</el-button>
        </div>
        <div>
            <div v-for="item in props.addgroupdata.addgroupsearchlist" :key="item.ID"
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
</template>

<script setup lang="ts">
import { fileurl } from '../../main'

const emit = defineEmits(['preapplyentergroup', 'searchgroup', 'beforeCloseAddGroupEvent'])
let props = defineProps({
    addgroupdata: {
        type: Object,
        required: true,
    },

})

const preapplyentergroup = (item) => {
    emit('preapplyentergroup',item)
}
const searchgroup = () => {
    emit('searchgroup')
}
const beforeCloseAddGroupEvent = () => {
    emit('beforeCloseAddGroupEvent')
}
</script>

<style scoped  lang="less">
@import url('./apply_join_group_dialog.vue.less');
</style>