import Notify from 'simple-notify'
import 'simple-notify/dist/simple-notify.css'

export function toast(status: "error" | "success" | "info" | "warning", message: string, text?: string) {
    return new Notify({
        status: status,
        title: message,
        text: text,
        autotimeout: 5000
    })
}