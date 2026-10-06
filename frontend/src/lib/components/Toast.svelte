<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import type { ToastType } from '$lib/stores/toast';

	export let message: string;
	export let type: ToastType = 'info';
	export let duration = 3000;
	export let onRemove: () => void;

	onMount(() => {
		const timer = setTimeout(onRemove, duration);
		return () => clearTimeout(timer);
	});
</script>

<button
	type="button"
	in:fade={{ duration: 150 }}
	out:fade={{ duration: 150 }}
	on:click={onRemove}
	class="pointer-events-auto max-w-sm rounded-lg px-4 py-2.5 text-sm shadow-lg {type === 'error'
		? 'bg-red-600 text-white'
		: 'bg-slate-900 text-white dark:bg-white dark:text-slate-900'}"
	role={type === 'error' ? 'alert' : 'status'}
>
	{message}
</button>
