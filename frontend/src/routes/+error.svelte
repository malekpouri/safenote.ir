<script lang="ts">
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';

	$: notFound = $page.status === 404;
</script>

<svelte:head>
	<title>{notFound ? $t.error.not_found_title : $t.error.generic_title} · {$t.app.title}</title>
	<meta name="robots" content="noindex" />
</svelte:head>

<section class="container-narrow pt-6 sm:pt-16">
	<p class="muted tabular-nums">{$page.status}</p>
	<h1 class="page-title mt-1">{notFound ? $t.error.not_found_title : $t.error.generic_title}</h1>
	<p class="page-subtitle">{notFound ? $t.error.not_found_body : $t.error.generic_body}</p>
	<a href="/" class="btn-primary btn-lg mt-6 w-full sm:w-auto">{$t.error.home}</a>
</section>
