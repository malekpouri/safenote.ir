import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'info';

export interface ToastMessage {
	id: number;
	message: string;
	type: ToastType;
}

export const toasts = writable<ToastMessage[]>([]);

let toastId = 0;
const MAX_VISIBLE = 3;

export const addToast = (message: string, type: ToastType = 'info') => {
	const id = toastId++;
	toasts.update((all) => [...all, { id, message, type }].slice(-MAX_VISIBLE));
};

export const removeToast = (id: number) => {
	toasts.update((all) => all.filter((t) => t.id !== id));
};
