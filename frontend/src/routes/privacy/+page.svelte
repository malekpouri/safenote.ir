<script lang="ts">
	import { t } from '$lib/i18n';
	import { CONTACT_EMAIL, SITE_URL } from '$lib/site';
</script>

<svelte:head>
	<title>{$t.privacy.title} · {$t.app.title}</title>
	<link rel="canonical" href="{SITE_URL}/privacy" />
	<meta property="og:title" content="{$t.privacy.title} · {$t.app.title}" />
</svelte:head>

<article class="container-narrow prose-page pt-6 sm:pt-16">
	<h1 class="page-title">{$t.privacy.title}</h1>
	<p class="muted mt-2">{$t.privacy.updated}</p>
	<p class="mt-6">{$t.privacy.intro}</p>
	{#each $t.privacy.sections as section}
		<h2>{section.h}</h2>
		<p>{section.p}</p>
	{/each}
	<h2>{$t.privacy.contact_title}</h2>
	<p>
		{$t.privacy.contact_body}
		<a href="mailto:{CONTACT_EMAIL}" class="link" dir="ltr">{CONTACT_EMAIL}</a>
	</p>
</article>
