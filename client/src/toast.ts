import 'simple-notify/dist/simple-notify.css';
import Swal from 'sweetalert2';
import { ClientResponseError } from 'pocketbase';

/**
 * Extracts a human-readable error message from various error types.
 * Handles PocketBase ClientResponseError, standard Error objects, and unknown types.
 */
export function extractErrorMessage(error: unknown): string {
    if (error instanceof ClientResponseError) {
        // PocketBase ClientResponseError - extract detailed information
        const parts: string[] = [];
        
        // Add the main message
        if (error.message) {
            parts.push(error.message);
        }
        
        // Add HTTP status info
        if (error.status) {
            parts.push(`(HTTP ${error.status})`);
        }
        
        // Extract field-specific errors from response.data
        if (error.response?.data) {
            const fieldErrors: string[] = [];
            for (const [field, value] of Object.entries(error.response.data)) {
                if (typeof value === 'object' && value !== null && 'message' in value) {
                    fieldErrors.push(`${field}: ${(value as { message: string }).message}`);
                } else if (typeof value === 'string') {
                    fieldErrors.push(`${field}: ${value}`);
                }
            }
            if (fieldErrors.length > 0) {
                parts.push('\n\nFelderfehler:\n' + fieldErrors.join('\n'));
            }
        }
        
        // Add response message if different from main message
        if (error.response?.message && error.response.message !== error.message) {
            parts.push(`\n\nServer: ${error.response.message}`);
        }
        
        return parts.join(' ') || 'Unbekannter PocketBase-Fehler';
    }
    
    if (error instanceof Error) {
        return error.message || 'Unbekannter Fehler';
    }
    
    if (typeof error === 'string') {
        return error;
    }
    
    // Try to stringify unknown objects
    try {
        return JSON.stringify(error, null, 2);
    } catch {
        return 'Unbekannter Fehler (nicht darstellbar)';
    }
}

/**
 * Shows an error toast with detailed error information extracted from the error object.
 * Use this for all PocketBase API error handling.
 */
export function showError(title: string, error: unknown): void {
    const message = extractErrorMessage(error);
    console.error(`${title}:`, error);
    
    // Use Swal.fire with html for multiline support
    Swal.fire({
        toast: true,
        icon: 'error',
        title: title,
        html: `<pre style="text-align: left; white-space: pre-wrap; word-wrap: break-word; font-size: 0.85em; max-height: 200px; overflow-y: auto;">${escapeHtml(message)}</pre>`,
        position: "top",
        showConfirmButton: true,
        confirmButtonText: 'OK',
        timer: 30000, // Longer timer for errors so users can read them
        timerProgressBar: true,
        width: 'auto',
        customClass: {
            popup: 'error-toast-popup'
        }
    });
}

/**
 * Shows an error alert (modal) with detailed error information.
 * Use this for critical errors that need user acknowledgment.
 */
export async function showErrorAlert(title: string, error: unknown): Promise<void> {
    const message = extractErrorMessage(error);
    console.error(`${title}:`, error);
    
    await Swal.fire({
        icon: 'error',
        title: title,
        html: `<pre style="text-align: left; white-space: pre-wrap; word-wrap: break-word; font-size: 0.9em; max-height: 400px; overflow-y: auto; background: #f5f5f5; padding: 1em; border-radius: 4px;">${escapeHtml(message)}</pre>`,
        confirmButtonText: 'OK',
        width: 600
    });
}

function escapeHtml(text: string): string {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

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