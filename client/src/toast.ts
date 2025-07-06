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

type Input =
    | 'text'
    | 'email'
    | 'password'
    | 'number'
    | 'tel'
    | 'search'
    | 'range'
    | 'textarea'
    | 'select'
    | 'radio'
    | 'checkbox'
    | 'url'
    | 'date'
    | 'datetime-local'
    | 'time'
    | 'week'
    | 'month'

export async function customPrompt(title: string, inputType: Input, inputPlaceholder?: string): Promise<string | null> {
    const { value: result } = await Swal.fire({
        title: title,
        input: inputType,
        inputPlaceholder: inputPlaceholder,
        showCancelButton: true,
        confirmButtonText: 'Ok',
        cancelButtonText: 'Abbrechen',
        inputValidator: (value) => {
            if (!value) {
                return 'Bitte einen Wert eingeben';
            }
        }
    });

    return result;
}

export async function customAlert(title: string, text?: string): Promise<void> {
    await Swal.fire({
        title: title,
        text: text,
        icon: 'info',
        confirmButtonText: 'Ok',
        showCancelButton: false,
    });
}