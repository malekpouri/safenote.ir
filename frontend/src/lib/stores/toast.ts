import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'info';

export interface ToastMessage {
    id: number;
    message: string;
    type: ToastType;
}

export const toasts = writable<ToastMessage[]>([]);

let toastId = 0;

export const addToast = (message: string, type: ToastType = 'info') => {
    const id = toastId++;
    toasts.update(all => [...all, { id, message, type }]);
};

export const removeToast = (id: number) => {
    toasts.update(all => all.filter(t => t.id !== id));
};
