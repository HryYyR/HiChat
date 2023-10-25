<template>
    <div class="header"></div>
    <div class="tabbar">
        <el-icon @click="toMin">
            <Minus />
        </el-icon>
        <el-icon @click="toMax">
            <CopyDocument />
        </el-icon>
        <el-icon @click="toClose">
            <Close />
        </el-icon>
    </div>
</template>

<script setup lang="ts">

const win: any = window
console.log(win);

const toMin = () => {
    win.api.toMin()
}
const toClose = () => {
    win.api.toClose()
}
const toMax = () => {
    win.api.toMax()
}
</script>

<style scoped lang="less" >
.header{
    position: absolute;
    -webkit-app-region: drag;
    width: 55vw;
    height: 45px;
    right: 110px;
}
.tabbar {
    width: 110px;
    position: absolute;
    right: 0;
    display: flex;
    z-index: 999;
    flex-direction: row;
    justify-content: flex-end;
    height: 45px;

    i {
        padding: 6px 10px;
    }

    i:hover {
        background-color: rgb(228, 228, 228);
    }
}</style>