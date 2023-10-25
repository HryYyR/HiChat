<template>
    <!-- 登录 -->
    <div class="view_container login_container" :style="{ marginTop: props.logindata.offset ? '-500px' : '0px' }">
        <div class="headericon">
            <el-icon size="30" color="white">
                <Plus />
            </el-icon>
        </div>
        <el-input class="view_input" v-model="props.logindata.username" placeholder="账号" size="large"
            clearable></el-input>
        <el-input type="password" class="view_input" v-model="props.logindata.password" show-password placeholder="密码"
            size="large" clearable></el-input>
        <div class="option">
            <div class="option_item">
                <span>记住密码: <el-switch size="small" v-model="props.logindata.rememberpassword" />
                </span>
            </div>
            <div class="option_item">
                <span>自动登录: <el-switch size="small" v-model="props.logindata.autologin" /></span>
            </div>
        </div>
        <div class="btn" @click="login">登录</div>
    </div>
</template>

<script setup lang="ts">

const emit = defineEmits(['login'])
let props = defineProps({
    logindata: {
        type: Object,
        required: true,
    },
})

const login = ()=>{
    emit('login')
}


</script>

<style scoped lang="less">
@import url('./login.less');
</style>