<template>
    <div class="message_time">
        {{ showttime }}
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
let props = defineProps({
    time: {
        type: String,
        required: true
    },
    pretime: {
        type: String,
        required: true
    }
})

const showttime = computed(() => {
    if (props.time == props.pretime) {
        return ""
    }
    let date = new Date(props.time)
    let predate = new Date(props.pretime)
    if (Math.abs(date.getMinutes() - predate.getMinutes())<10){
        return ""
    }

    return getRelativeTimeString(date)
})

function getRelativeTimeString(
    date: Date | number,
    lang = navigator.language
): string {
    const timeMs = typeof date === "number" ? date : date.getTime();
    const deltaSeconds = Math.round((timeMs - Date.now()) / 1000);
    const cutoffs = [60, 3600, 86400, 86400 * 7, 86400 * 30, 86400 * 365, Infinity];
    const units: Intl.RelativeTimeFormatUnit[] = ["second", "minute", "hour", "day", "week", "month", "year"];
    const unitIndex = cutoffs.findIndex(cutoff => cutoff > Math.abs(deltaSeconds));
    const divisor = unitIndex ? cutoffs[unitIndex - 1] : 1;
    const rtf = new Intl.RelativeTimeFormat(lang, { numeric: "auto" });
    return rtf.format(Math.floor(deltaSeconds / divisor), units[unitIndex]);
}



</script>

<style scoped lang="less">
@import url('./message_time.less');
</style>