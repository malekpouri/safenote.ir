<script lang="ts">
	import { locale, t } from '$lib/i18n';
	import Logo from './Logo.svelte';

	function toggleLanguage() {
		locale.set($locale === 'en' ? 'fa' : 'en');
	}
</script>

<header class="pt-[env(safe-area-inset-top)]">
	<div class="container-narrow flex h-16 items-center justify-between">
		<a href="/" class="-ms-1 inline-flex min-h-[44px] items-center gap-2 px-1" aria-label={$t.app.title}>
			<Logo class="h-7 w-7" />
			<span class="text-[17px] font-bold tracking-tight">
				{#if $locale === 'fa'}
					سیف<span class="text-brand-600 dark:text-brand-400">‌نوت</span>
				{:else}
					Safe<span class="text-brand-600 dark:text-brand-400">Note</span>
				{/if}
			</span>
		</a>
		<button
			type="button"
			on:click={toggleLanguage}
			class="chip -me-1"
			aria-label={$t.nav.switch_language_label}
			lang={$locale === 'en' ? 'fa' : 'en'}
		>
			{$t.nav.switch_language}
		</button>
	</div>
</header>
