<script lang="ts">
	import { tick } from 'svelte';
	import { fade, slide } from 'svelte/transition';
	import { NOTE_MAX_CHARS, encryptNote, generateLinkKey, generatePassword } from '$lib/crypto';
	import { addToast } from '$lib/stores/toast';
	import { dateTime, fmt, num, t } from '$lib/i18n';
	import { DEFAULT_EXPIRATION, EXPIRATION_OPTIONS, SITE_URL, VIEW_OPTIONS, shareOrigin } from '$lib/site';
	import { apiErrorMessage } from '$lib/api';
	import Icon from '$lib/components/Icon.svelte';
	import Spinner from '$lib/components/Spinner.svelte';

	let note = '';
	let views: number = VIEW_OPTIONS[0];
	let expiration: number = DEFAULT_EXPIRATION;
	let password = '';
	let showPassword = false;
	let showOptions = false;
	let loading = false;

	let created: {
		id: string;
		link: string;
		accessToken: string;
		views: number;
		expiresAt: string;
		hasPassword: boolean;
	} | null = null;
	let copied = false;
	let confirmingDelete = false;
	let deleting = false;
	let linkInput: HTMLTextAreaElement;
	const canShare = typeof navigator !== 'undefined' && 'share' in navigator;
	const isDesktop = () => typeof window !== 'undefined' && window.matchMedia('(min-width: 640px)').matches;

	$: expirationLabels = {
		60: $t.home.hours_1,
		1440: $t.home.hours_24,
		10080: $t.home.days_7,
		43200: $t.home.days_30
	} as Record<number, string>;
	$: viewsLabel = (n: number) => fmt(n === 1 ? $t.home.views_option : $t.home.views_option_plural, { n: $num(n) });
	$: nearLimit = note.length > NOTE_MAX_CHARS - 1000;

	async function createNote() {
		if (!note.trim() || loading) return;
		loading = true;
		try {
			const linkKey = generateLinkKey();
			const { payload, accessToken, salt } = await encryptNote(note, linkKey, password);
			const res = await fetch('/api/notes', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					encrypted_data: payload,
					access_token: accessToken,
					salt,
					is_password_protected: !!password,
					views_remaining: views,
					expiration
				})
			});
			const data = await res.json().catch(() => ({}));
			if (!res.ok) throw new Error(apiErrorMessage(res.status, $t, data.error));

			created = {
				id: data.id,
				link: `${shareOrigin(window.location)}/${data.id}#${linkKey}`,
				accessToken,
				views,
				expiresAt: data.expires_at,
				hasPassword: !!password
			};
			note = '';
			await tick();
			window.scrollTo({ top: 0 });
			// On phones, programmatic selection pops up selection handles; only do it on desktop.
			if (isDesktop()) linkInput?.select();
		} catch (e) {
			addToast(e instanceof Error ? e.message : $t.toast.error, 'error');
		} finally {
			loading = false;
		}
	}

	function onKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
			e.preventDefault();
			createNote();
		}
	}

	async function copyLink() {
		if (!created) return;
		try {
			await navigator.clipboard.writeText(created.link);
			copied = true;
			setTimeout(() => (copied = false), 2000);
		} catch {
			linkInput?.select();
			addToast($t.toast.copy_failed, 'error');
		}
	}

	async function shareLink() {
		if (!created) return;
		try {
			await navigator.share({ title: $t.app.title, url: created.link });
		} catch {
			/* user cancelled */
		}
	}

	async function deleteNote() {
		if (!created) return;
		deleting = true;
		try {
			const res = await fetch(`/api/notes/${created.id}`, {
				method: 'DELETE',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ access_token: created.accessToken })
			});
			if (!res.ok && res.status !== 404) {
				const data = await res.json().catch(() => ({}));
				throw new Error(apiErrorMessage(res.status, $t, data.error));
			}
			addToast($t.toast.deleted, 'success');
			reset();
		} catch (e) {
			addToast(e instanceof Error ? e.message : $t.toast.error, 'error');
		} finally {
			deleting = false;
			confirmingDelete = false;
		}
	}

	function reset() {
		created = null;
		password = '';
		showPassword = false;
		showOptions = false;
		confirmingDelete = false;
	}

	/** Pills open the options panel and focus the matching field. */
	async function openOption(field: 'views' | 'expiration' | 'password') {
		showOptions = true;
		await tick();
		document.getElementById(field)?.focus();
	}

	function useGeneratedPassword() {
		password = generatePassword();
		showPassword = true;
	}
</script>

<svelte:head>
	<title>{$t.app.title} · {$t.app.tagline}</title>
	<link rel="canonical" href="{SITE_URL}/" />
	<meta property="og:title" content="{$t.app.title} · {$t.app.tagline}" />
	<meta property="og:description" content={$t.app.description} />
	<meta name="twitter:title" content="{$t.app.title} · {$t.app.tagline}" />
	<meta name="twitter:description" content={$t.app.description} />
</svelte:head>

<section class="container-narrow pt-6 sm:pt-14">
	{#if created}
		<div in:fade={{ duration: 200 }} class="text-center">
			<span class="mx-auto flex h-14 w-14 animate-pop items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-400">
				<Icon name="check" class="h-7 w-7" strokeWidth={2.5} />
			</span>
			<h1 class="page-title mt-5">{$t.success.title}</h1>
			<p class="page-subtitle">{$t.success.message}</p>
		</div>

		<div class="card mt-7 p-2" in:fade={{ duration: 200, delay: 100 }}>
			<label for="link" class="sr-only">{$t.success.link_label}</label>
			<!-- A textarea so the whole link, including the key after #, wraps and stays visible on phones. -->
			<textarea
				bind:this={linkInput}
				id="link"
				readonly
				rows={isDesktop() ? 1 : 2}
				value={created.link}
				dir="ltr"
				class="block w-full resize-none break-all border-0 bg-transparent px-3 py-2 font-mono text-base leading-7 text-slate-800 focus:outline-none focus:ring-0 sm:text-sm dark:text-slate-100"
				on:click={(e) => e.currentTarget.select()}
			></textarea>
			<div class="flex gap-2">
				<button type="button" class="btn-primary btn-lg flex-1" on:click={copyLink}>
					<Icon name={copied ? 'check' : 'copy'} class="h-4 w-4" />
					{copied ? $t.success.copied : $t.success.copy}
				</button>
				{#if canShare}
					<button type="button" class="btn-secondary btn-lg" on:click={shareLink} aria-label={$t.success.share}>
						<Icon name="share" class="h-4 w-4" />
					</button>
				{/if}
			</div>
		</div>

		<ul class="mt-4 flex flex-wrap justify-center gap-2 text-[13px] text-slate-500 dark:text-slate-400">
			<li class="inline-flex items-center gap-1.5 rounded-full bg-slate-100 px-3 py-1.5 dark:bg-slate-800/70">
				<Icon name="fire" class="h-3.5 w-3.5 text-orange-500" />
				{fmt(created.views === 1 ? $t.success.summary_views : $t.success.summary_views_plural, { n: $num(created.views) })}
			</li>
			<li class="inline-flex items-center gap-1.5 rounded-full bg-slate-100 px-3 py-1.5 dark:bg-slate-800/70">
				<Icon name="clock" class="h-3.5 w-3.5 text-brand-500" />
				{fmt($t.success.summary_expires, { date: $dateTime(created.expiresAt) })}
			</li>
		</ul>
		{#if created.hasPassword}
			<p class="muted mt-3 flex items-center justify-center gap-1.5 text-center">
				<Icon name="key" class="h-4 w-4 shrink-0 text-amber-500" />
				{$t.success.summary_password}
			</p>
		{/if}

		<div class="mt-8 flex flex-wrap items-center justify-between gap-x-4">
			<button type="button" class="btn-text" on:click={reset}>
				<Icon name="plus" class="h-4 w-4" />
				{$t.success.new_note}
			</button>
			{#if confirmingDelete}
				<span class="flex items-center gap-4" in:fade={{ duration: 100 }}>
					<span class="muted">{$t.success.delete_confirm}</span>
					<button type="button" class="btn-text-danger font-semibold" on:click={deleteNote} disabled={deleting}>
						{#if deleting}<Spinner />{/if}
						{$t.success.delete_yes}
					</button>
					<button type="button" class="btn-text" on:click={() => (confirmingDelete = false)}>{$t.success.cancel}</button>
				</span>
			{:else}
				<button type="button" class="btn-text-danger" on:click={() => (confirmingDelete = true)}>
					{$t.success.delete}
				</button>
			{/if}
		</div>
	{:else}
		<div in:fade={{ duration: 200 }}>
			<h1 class="page-title">
				{$t.home.title_lead}
				<span class="text-gradient whitespace-nowrap">{$t.home.title_accent}</span>
			</h1>
			<p class="page-subtitle">{$t.home.subtitle}</p>

			<form on:submit|preventDefault={createNote} class="mt-6">
				<div class="card transition focus-within:border-brand-300 focus-within:ring-4 focus-within:ring-brand-500/10 dark:focus-within:border-brand-600">
					<label for="note" class="sr-only">{$t.home.note_label}</label>
					<textarea
						id="note"
						bind:value={note}
						on:keydown={onKeydown}
						rows="6"
						maxlength={NOTE_MAX_CHARS}
						class="block min-h-[10rem] w-full resize-y max-[359px]:h-[7.5rem] max-[359px]:min-h-[7.5rem] rounded-t-2xl border-0 bg-transparent px-4 pb-2 pt-4 text-base leading-7 placeholder:text-slate-400 focus:outline-none focus:ring-0 sm:text-[15px] dark:placeholder:text-slate-600"
						placeholder={$t.home.placeholder}
						autocomplete="off"
						spellcheck="false"
						required
					></textarea>

					<div class="flex flex-wrap items-center gap-1.5 px-3 pb-3 sm:gap-2">
						<button type="button" class="chip" class:chip-on={views > 1} on:click={() => openOption('views')}>
							<Icon name="fire" class="h-3.5 w-3.5 text-orange-500" />
							{viewsLabel(views)}
						</button>
						<button type="button" class="chip" on:click={() => openOption('expiration')}>
							<Icon name="clock" class="h-3.5 w-3.5 text-brand-500" />
							{expirationLabels[expiration]}
						</button>
						<button
							type="button"
							class="chip"
							class:chip-on={!!password}
							on:click={() => openOption('password')}
							aria-label={$t.home.password_label}
						>
							<Icon name={password ? 'lock' : 'key'} class="h-3.5 w-3.5 text-amber-500" />
							<!-- Icon-only on very narrow phones so the three pills stay on one row. -->
							<span class="max-[359px]:hidden">{$t.home.password_label}</span>
						</button>
						{#if nearLimit}
							<span class="ms-auto pe-1 text-xs tabular-nums text-slate-400" aria-live="polite">
								{fmt($t.home.chars, { n: $num(note.length), max: $num(NOTE_MAX_CHARS) })}
							</span>
						{/if}
					</div>

					{#if showOptions}
						<div
							id="note-options"
							class="space-y-4 rounded-b-2xl border-t border-slate-100 bg-slate-50/70 px-4 pb-4 pt-4 dark:border-slate-800 dark:bg-slate-950/40"
							transition:slide={{ duration: 160 }}
						>
							<div class="grid grid-cols-2 gap-3">
								<div>
									<label for="views" class="label">{$t.home.views_label}</label>
									<select id="views" bind:value={views} class="input">
										{#each VIEW_OPTIONS as v}
											<option value={v}>{viewsLabel(v)}</option>
										{/each}
									</select>
								</div>
								<div>
									<label for="expiration" class="label">{$t.home.expires_label}</label>
									<select id="expiration" bind:value={expiration} class="input">
										{#each EXPIRATION_OPTIONS as minutes}
											<option value={minutes}>{expirationLabels[minutes]}</option>
										{/each}
									</select>
								</div>
							</div>
							<div>
								<label for="password" class="label">
									{$t.home.password_label}
									<span class="text-slate-400 dark:text-slate-500">({$t.home.password_optional})</span>
								</label>
								<div class="relative">
									<input
										id="password"
										type={showPassword ? 'text' : 'password'}
										value={password}
										on:input={(e) => (password = e.currentTarget.value)}
										class="input pe-[5.5rem]"
										placeholder={$t.home.password_placeholder}
										autocomplete="new-password"
										autocapitalize="off"
										autocorrect="off"
										spellcheck="false"
										enterkeyhint="done"
									/>
									<div class="absolute inset-y-0 end-0.5 flex items-center">
										<button
											type="button"
											class="icon-btn hover:text-brand-600"
											on:click={useGeneratedPassword}
											aria-label={$t.home.password_generate}
											title={$t.home.password_generate}
										>
											<Icon name="sparkles" class="h-4 w-4" />
										</button>
										<button
											type="button"
											class="icon-btn"
											on:click={() => (showPassword = !showPassword)}
											aria-label={showPassword ? $t.home.password_hide : $t.home.password_show}
											title={showPassword ? $t.home.password_hide : $t.home.password_show}
										>
											<Icon name={showPassword ? 'eye-off' : 'eye'} class="h-4 w-4" />
										</button>
									</div>
								</div>
								{#if password}
									<p class="mt-1.5 text-xs text-slate-500 dark:text-slate-400">{$t.home.password_hint}</p>
								{/if}
							</div>
						</div>
					{/if}
				</div>

				<button type="submit" class="btn-primary btn-lg mt-4 w-full" disabled={loading || !note.trim()}>
					{#if loading}
						<Spinner class="h-4 w-4" />
						{$t.home.submitting}
					{:else}
						<Icon name="lock" class="h-4 w-4" />
						{$t.home.submit}
					{/if}
				</button>
			</form>

			<p class="muted mt-6 flex flex-wrap items-center justify-center gap-x-1.5 text-center">
				<Icon name="shield" class="h-4 w-4 text-emerald-500" />
				{$t.home.footnote}
				<a href="/about" class="link whitespace-nowrap">{$t.home.learn_more}</a>
			</p>
		</div>
	{/if}
</section>
