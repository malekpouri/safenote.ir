<script lang="ts">
	import { t } from '$lib/i18n';
	import { GITHUB_URL, POWERED_BY } from '$lib/site';

	const year = new Date().getFullYear();

	$: links = [
		{ href: '/about', label: $t.footer.about },
		{ href: '/privacy', label: $t.footer.privacy },
		{ href: '/admin', label: $t.footer.stats },
		{ href: GITHUB_URL, label: $t.footer.source, external: true }
	];
</script>

<footer class="container-narrow pb-safe mt-16">
	<div class="flex flex-col items-center gap-2 border-t border-slate-100 pt-6 text-sm text-slate-400 sm:flex-row sm:justify-between dark:border-slate-900">
		<nav class="flex flex-wrap justify-center gap-x-5" aria-label="Footer">
			{#each links as link}
				<a
					href={link.href}
					class="inline-flex min-h-[44px] items-center transition-colors hover:text-slate-900 sm:min-h-0 dark:hover:text-white"
					target={link.external ? '_blank' : undefined}
					rel={link.external ? 'noopener noreferrer' : undefined}
				>
					{link.label}
				</a>
			{/each}
		</nav>
		<p>
			&copy; {year} {$t.app.title} · {$t.footer.powered}
			<a href={POWERED_BY.url} target="_blank" rel="noopener noreferrer" class="transition-colors hover:text-slate-900 dark:hover:text-white">{POWERED_BY.name}</a>
		</p>
	</div>
</footer>
