<script lang="ts">
	import { num, t } from '$lib/i18n';
	import { CONTACT_EMAIL, GITHUB_URL, SITE_URL } from '$lib/site';

	$: steps = [
		[$t.how.step1_title, $t.how.step1_body],
		[$t.how.step2_title, $t.how.step2_body],
		[$t.how.step3_title, $t.how.step3_body]
	];
	$: faqs = [
		[$t.faq.q1, $t.faq.a1],
		[$t.faq.q2, $t.faq.a2],
		[$t.faq.q3, $t.faq.a3],
		[$t.faq.q4, $t.faq.a4]
	];
</script>

<svelte:head>
	<title>{$t.about.title} · {$t.app.title}</title>
	<link rel="canonical" href="{SITE_URL}/about" />
	<meta property="og:title" content="{$t.about.title} · {$t.app.title}" />
</svelte:head>

<article class="container-narrow prose-page pt-6 sm:pt-16">
	<h1 class="page-title">{$t.about.title}</h1>
	<p class="mt-3">{$t.about.lead}</p>

	<h2 id="how">{$t.how.title}</h2>
	<ol class="space-y-3">
		{#each steps as [title, body], i}
			<li class="flex gap-3">
				<span class="w-4 shrink-0 tabular-nums text-slate-400">{$num(i + 1)}.</span>
				<span><strong class="font-medium text-slate-900 dark:text-white">{title}.</strong> {body}</span>
			</li>
		{/each}
	</ol>

	<h2 id="faq">{$t.faq.title}</h2>
	<dl class="space-y-5">
		{#each faqs as [q, a]}
			<div>
				<dt class="text-[15px] font-medium">{q}</dt>
				<dd class="mt-1">{a}</dd>
			</div>
		{/each}
	</dl>

	<h2>{$t.about.contact_title}</h2>
	<p>
		{$t.about.contact_body}
		<a href="mailto:{CONTACT_EMAIL}" class="link" dir="ltr">{CONTACT_EMAIL}</a>
		{$t.about.contact_or}
		<a href={GITHUB_URL} target="_blank" rel="noopener noreferrer" class="link">GitHub</a>.
	</p>
</article>
