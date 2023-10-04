

export function nowtime() {
    let time = new Date()
    let year = time.getFullYear()
    let month: any = time.getMonth()
    let day: any = time.getDate()
    let hour: any = time.getHours()
    let minute: any = time.getMinutes()
    let second: any = time.getSeconds()
    if (month <= 9) {
        month = "0" + month
    }
    if (hour <= 9) {
        hour = "0" + hour
    }
    if (minute <= 9) {
        minute = "0" + minute
    }
    if (second <= 9) {
        second = "0" + second
    }
    let text = year + '-' + month + '-' + day + ' ' + hour + ':' + minute + ':' + second
    return text
}