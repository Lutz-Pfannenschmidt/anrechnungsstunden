import 'simple-notify/dist/simple-notify.css';
import Swal from 'sweetalert2';

export function toast(status: "success" | "error" | "warning" | "info" | "question", title: string, text?: string) {
    Swal.fire({
        toast: true,
        icon: status,
        title: title,
        text: text,
        position: "top",
        showConfirmButton: false,
        timer: 5000,
        timerProgressBar: true,
        didOpen: (toast) => {
            toast.onmouseenter = Swal.stopTimer;
            toast.onmouseleave = Swal.resumeTimer;
        }
    });
}