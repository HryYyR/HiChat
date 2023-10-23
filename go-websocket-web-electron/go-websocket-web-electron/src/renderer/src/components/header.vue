<template>
    <div class="header">123
        <span class="float-right btn-opt text-white">
            <i class="el-icon-minus mr-10" @click="toMin()"></i>
            <i class="el-icon-close" @click="toClose()"></i>
        </span>
    </div>
</template>

<script setup lang="ts">
import { ipcRenderer } from 'electron';

const toMin = () => {
    ipcRenderer.send('min-app')
}
const toClose = () => {
    ipcRenderer.send('close-app')
}
</script>

<style scoped lang="less" >
.header {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center; 
    -webkit-app-region: drag; //设置可拖动
}
</style>