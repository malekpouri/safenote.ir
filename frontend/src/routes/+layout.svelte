<script lang="ts">
	import '@fontsource-variable/inter';
	import 'vazirmatn/Vazirmatn-Variable-font-face.css';
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import Toast from '$lib/components/Toast.svelte';
	import { toasts, removeToast } from '$lib/stores/toast';
	import { locale, t } from '$lib/i18n';
	import type { LayoutData } from './$types';

	export let data: LayoutData;

	// The server resolved the locale (cookie / Accept-Language) and already set
	// <html lang dir>; sync the store before the first render so SSR and
	// hydration agree.
	locale.init(data.locale);
</script>

<svelte:head>
	<meta name="description" content={$t.app.description} />
	<meta property="og:locale" content={$locale === 'fa' ? 'fa_IR' : 'en_US'} />
</svelte:head>

<a
	href="#main"
	class="sr-only z-50 rounded-lg bg-slate-900 px-4 py-2 text-white focus:not-sr-only focus:fixed focus:start-4 focus:top-4"
>
	{$t.app.skip_to_content}
</a>

<div class="relative flex min-h-[100dvh] flex-col overflow-x-clip">
	<!-- A single soft glow keeps the page fresh without adding elements. -->
	<div
		aria-hidden="true"
		class="pointer-events-none absolute inset-x-0 top-0 -z-10 h-[28rem] bg-[radial-gradient(40rem_20rem_at_50%_-4rem,theme(colors.brand.100),transparent)] dark:bg-[radial-gradient(40rem_20rem_at_50%_-4rem,theme(colors.brand.950),transparent)]"
	></div>
	<Header />

	<main id="main" class="flex-1">
		<noscript>
			<p class="container-narrow muted pt-6">{$t.app.noscript}</p>
		</noscript>
		<slot />
	</main>

	<Footer />
</div>

<div
	class="pointer-events-none fixed inset-x-3 top-[calc(env(safe-area-inset-top)+0.75rem)] z-50 flex flex-col items-center gap-2 sm:bottom-6 sm:top-auto"
	aria-live="polite"
>
	{#each $toasts as toast (toast.id)}
		<Toast message={toast.message} type={toast.type} onRemove={() => removeToast(toast.id)} />
	{/each}
</div>
