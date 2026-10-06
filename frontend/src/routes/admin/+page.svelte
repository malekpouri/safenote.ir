<script lang="ts">
	import { onMount } from 'svelte';
	import { dateTime, fmt, num, t } from '$lib/i18n';

	let stats: { active_notes: number; total_notes_created: number } | null = null;
	let loading = true;
	let failed = false;
	let updatedAt: Date | null = null;

	async function fetchStats() {
		loading = true;
		try {
			const res = await fetch('/api/admin/stats');
			if (!res.ok) throw new Error();
			stats = await res.json();
			updatedAt = new Date();
			failed = false;
		} catch {
			failed = true;
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchStats();
		const timer = setInterval(() => {
			if (document.visibilityState === 'visible') fetchStats();
		}, 30_000);
		return () => clearInterval(timer);
	});

	$: items = [
		[$t.admin.active, stats?.active_notes],
		[$t.admin.total, stats?.total_notes_created]
	] as const;
</script>

<svelte:head>
	<title>{$t.admin.title} · {$t.app.title}</title>
	<meta name="robots" content="noindex" />
</svelte:head>

<section class="container-narrow pt-6 sm:pt-16">
	<h1 class="page-title">{$t.admin.title}</h1>
	<p class="page-subtitle">{$t.admin.subtitle}</p>

	{#if failed}
		<p class="mt-6 text-sm text-red-600 dark:text-red-400">{$t.admin.error}</p>
	{/if}

	<dl class="mt-8 grid grid-cols-2 gap-6">
		{#each items as [label, value]}
			<div>
				<dt class="muted">{label}</dt>
				<dd class="mt-1 text-3xl font-semibold tabular-nums tracking-tight">
					{value === undefined ? '—' : $num(value)}
				</dd>
			</div>
		{/each}
	</dl>

	<p class="muted mt-8 flex items-center gap-3">
		<button type="button" class="btn-text underline underline-offset-4" on:click={fetchStats} disabled={loading}>
			{$t.admin.refresh}
		</button>
		{#if updatedAt}<span>{fmt($t.admin.updated, { time: $dateTime(updatedAt) })}</span>{/if}
	</p>
</section>
