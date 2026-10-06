<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { fade } from 'svelte/transition';
	import { decryptNote, deriveKeys, isV2Key, legacyPasswordHash, type NoteKeys } from '$lib/crypto';
	import { addToast } from '$lib/stores/toast';
	import { dateTime, fmt, num, t } from '$lib/i18n';
	import { apiErrorMessage } from '$lib/api';
	import Icon from '$lib/components/Icon.svelte';
	import Spinner from '$lib/components/Spinner.svelte';

	type State = 'loading' | 'missing-key' | 'not-found' | 'ready' | 'revealed' | 'deleted';

	interface Meta {
		is_password_protected: boolean;
		legacy: boolean;
		expires_at: string;
		salt: string;
	}

	export let id: string;
	let state: State = 'loading';
	let key = '';
	let meta: Meta | null = null;
	let password = '';
	let showPassword = false;
	let passwordError = '';
	let opening = false;
	let content = '';
	let viewsLeft = 0;
	let expiresAt = '';
	let keys: NoteKeys | null = null;
	let legacyPayload: string | null = null;
	let confirmingDelete = false;
	let deleting = false;
	let passwordInput: HTMLInputElement;

	onMount(async () => {
		key = decodeURIComponent(window.location.hash.slice(1));
		if (!key) {
			state = 'missing-key';
			return;
		}
		try {
			const res = await fetch(`/api/notes/${encodeURIComponent(id)}`);
			if (res.status === 404) {
				state = 'not-found';
				return;
			}
			if (!res.ok) {
				const data = await res.json().catch(() => ({}));
				throw new Error(apiErrorMessage(res.status, $t, data.error));
			}
			meta = await res.json();
			if (meta && !meta.legacy && !isV2Key(key)) {
				state = 'missing-key';
				return;
			}
			state = 'ready';
			if (meta?.is_password_protected) {
				await tick();
				passwordInput?.focus();
			}
		} catch (e) {
			state = 'not-found';
			addToast(e instanceof Error ? e.message : $t.toast.error, 'error');
		}
	});

	async function getKeys(): Promise<NoteKeys> {
		keys = await deriveKeys(key, password, meta?.salt ?? '');
		return keys;
	}

	async function reveal() {
		if (!meta || opening) return;
		if (meta.is_password_protected && !password) {
			passwordInput?.focus();
			return;
		}
		opening = true;
		passwordError = '';
		try {
			let payload: string;
			if (meta.legacy) {
				// Legacy notes cannot verify the password server-side, so keep the
				// payload around and let the user retry locally.
				if (legacyPayload === null) {
					const data = await open({});
					if (!data) return;
					legacyPayload = data.encrypted_data;
					viewsLeft = data.views_remaining;
					expiresAt = data.expires_at;
				}
				payload = legacyPayload!;
			} else {
				const k = await getKeys();
				const data = await open({ access_token: k.accessToken });
				if (!data) return;
				payload = data.encrypted_data;
				viewsLeft = data.views_remaining;
				expiresAt = data.expires_at;
			}

			try {
				content = await decryptNote(payload, key, password, { keys: keys ?? undefined });
			} catch {
				if (meta.legacy && meta.is_password_protected) {
					passwordError = $t.view.wrong_password_legacy;
				} else {
					addToast($t.toast.decrypt_failed, 'error');
				}
				return;
			}
			legacyPayload = null;
			state = 'revealed';
		} catch (e) {
			addToast(e instanceof Error ? e.message : $t.toast.error, 'error');
		} finally {
			opening = false;
		}
	}

	/** Returns the opened note, or null when the UI has already handled a failure. */
	async function open(body: Record<string, string>) {
		const res = await fetch(`/api/notes/${encodeURIComponent(id)}/open`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body)
		});
		const data = await res.json().catch(() => ({}));
		if (res.status === 401) {
			if (meta?.is_password_protected) {
				passwordError = $t.view.wrong_password;
				await tick();
				passwordInput?.select();
			} else {
				addToast($t.toast.decrypt_failed, 'error');
			}
			return null;
		}
		if (res.status === 404) {
			state = 'not-found';
			return null;
		}
		if (!res.ok) throw new Error(apiErrorMessage(res.status, $t, data.error));
		return data as { encrypted_data: string; views_remaining: number; expires_at: string };
	}

	async function deleteNote() {
		if (!meta) return;
		deleting = true;
		passwordError = '';
		try {
			const body: Record<string, string> = {};
			if (meta.legacy) {
				if (meta.is_password_protected) body.password_hash = await legacyPasswordHash(password);
			} else {
				body.access_token = (keys ?? (await getKeys())).accessToken;
			}
			const res = await fetch(`/api/notes/${encodeURIComponent(id)}`, {
				method: 'DELETE',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(body)
			});
			if (res.status === 401) {
				passwordError = $t.view.wrong_password;
				keys = null;
				return;
			}
			if (!res.ok && res.status !== 404) {
				const data = await res.json().catch(() => ({}));
				throw new Error(apiErrorMessage(res.status, $t, data.error));
			}
			content = '';
			state = 'deleted';
			addToast($t.toast.deleted, 'success');
		} catch (e) {
			addToast(e instanceof Error ? e.message : $t.toast.error, 'error');
		} finally {
			deleting = false;
			confirmingDelete = false;
		}
	}

	async function copyContent() {
		try {
			await navigator.clipboard.writeText(content);
			addToast($t.toast.copied, 'success');
		} catch {
			addToast($t.toast.copy_failed, 'error');
		}
	}

	$: canDelete =
		(state === 'ready' && (!meta?.is_password_protected || !!password)) || (state === 'revealed' && viewsLeft > 0);
</script>

<svelte:head>
	<title>{$t.view.page_title} · {$t.app.title}</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<section class="container-narrow pt-6 sm:pt-14">
	{#if state === 'loading'}
		<p class="muted flex items-center justify-center gap-2 py-16" in:fade>
			<Spinner class="h-4 w-4 text-brand-500" />
			{$t.view.loading}
		</p>
	{:else if state === 'missing-key' || state === 'not-found' || state === 'deleted'}
		<div class="text-center" in:fade={{ duration: 200 }}>
			<span
				class="mx-auto flex h-14 w-14 items-center justify-center rounded-full {state === 'deleted'
					? 'animate-pop bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-400'
					: 'bg-slate-100 text-slate-400 dark:bg-slate-800'}"
			>
				<Icon name={state === 'deleted' ? 'check' : state === 'missing-key' ? 'key' : 'fire'} class="h-7 w-7" />
			</span>
			<h1 class="page-title mt-5">
				{state === 'missing-key' ? $t.view.missing_key_title : state === 'deleted' ? $t.view.deleted_title : $t.view.not_found_title}
			</h1>
			<p class="page-subtitle mx-auto max-w-sm">
				{state === 'missing-key' ? $t.view.missing_key_body : state === 'deleted' ? $t.view.deleted_body : $t.view.not_found_body}
			</p>
			<a href="/" class="btn-primary btn-lg mt-7 w-full sm:w-auto sm:px-8">
				<Icon name="plus" class="h-4 w-4" />
				{$t.view.new_note}
			</a>
		</div>
	{:else if state === 'ready' && meta}
		<form on:submit|preventDefault={reveal} in:fade={{ duration: 200 }}>
			<div class="text-center">
				<span class="mx-auto flex h-14 w-14 animate-pop items-center justify-center rounded-full bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400">
					<Icon name={meta.is_password_protected ? 'lock' : 'mail'} class="h-7 w-7" />
				</span>
				<h1 class="page-title mt-5">{meta.is_password_protected ? $t.view.password_title : $t.view.ready_title}</h1>
				<p class="page-subtitle flex items-center justify-center gap-1.5">
					<Icon name="fire" class="h-4 w-4 shrink-0 text-orange-500" />
					{$t.view.ready_body}
				</p>
			</div>

			{#if meta.is_password_protected}
				<div class="mt-7">
					<label for="password" class="sr-only">{$t.view.password_label}</label>
					<div class="relative">
						<input
							bind:this={passwordInput}
							id="password"
							type={showPassword ? 'text' : 'password'}
							value={password}
							on:input={(e) => {
								password = e.currentTarget.value;
								passwordError = '';
								keys = null;
							}}
							class="input pe-12"
							class:!border-red-400={passwordError}
							placeholder={$t.view.password_placeholder}
							autocomplete="off"
							autocapitalize="off"
							autocorrect="off"
							spellcheck="false"
							enterkeyhint="go"
							aria-invalid={!!passwordError}
							aria-describedby={passwordError ? 'password-error' : undefined}
						/>
						<button
							type="button"
							class="icon-btn absolute end-0.5 top-1/2 -translate-y-1/2"
							on:click={() => (showPassword = !showPassword)}
							aria-label={showPassword ? $t.home.password_hide : $t.home.password_show}
						>
							<Icon name={showPassword ? 'eye-off' : 'eye'} class="h-4 w-4" />
						</button>
					</div>
					{#if passwordError}
						<p id="password-error" class="mt-2 text-sm text-red-600 dark:text-red-400" in:fade>{passwordError}</p>
					{/if}
				</div>
			{/if}

			<button
				type="submit"
				class="btn-primary btn-lg mt-6 w-full"
				disabled={opening || (meta.is_password_protected && !password)}
			>
				{#if opening}
					<Spinner class="h-4 w-4" />
					{$t.view.opening}
				{:else}
					<Icon name="eye" class="h-4 w-4" />
					{$t.view.open}
				{/if}
			</button>
		</form>
	{:else if state === 'revealed'}
		<div in:fade={{ duration: 200 }}>
			<h1 class="page-title">{$t.view.revealed_title}</h1>
			<p class="page-subtitle flex items-start gap-1.5">
				{#if viewsLeft === 0}
					<Icon name="fire" class="mt-1 h-4 w-4 shrink-0 text-orange-500" />
					{$t.view.burned}
				{:else}
					<Icon name="clock" class="mt-1 h-4 w-4 shrink-0 text-brand-500" />
					<span>
						{fmt(viewsLeft === 1 ? $t.view.views_left : $t.view.views_left_plural, { n: $num(viewsLeft) })}
						{#if expiresAt}{fmt($t.view.expires, { date: $dateTime(expiresAt) })}.{/if}
					</span>
				{/if}
			</p>

			<div class="card mt-6 p-2">
				<div
					class="max-h-[55dvh] overflow-auto overscroll-contain whitespace-pre-wrap break-words px-3 py-2 text-[15px] leading-7 sm:max-h-[28rem]"
					dir="auto"
					tabindex="0"
					role="textbox"
					aria-readonly="true"
					aria-label={$t.view.revealed_title}
				>{content}</div>
				<button type="button" class="btn-primary btn-lg mt-2 w-full" on:click={copyContent}>
					<Icon name="copy" class="h-4 w-4" />
					{$t.view.copy}
				</button>
			</div>
		</div>
	{/if}

	{#if state === 'revealed' || canDelete}
		<div class="mt-8 flex flex-wrap items-center justify-between gap-x-4">
			{#if state === 'revealed'}
				<a href="/" class="btn-text">
					<Icon name="plus" class="h-4 w-4" />
					{$t.view.new_note}
				</a>
			{:else}
				<span></span>
			{/if}
			{#if canDelete}
				{#if confirmingDelete}
					<span class="flex items-center gap-4" in:fade={{ duration: 100 }}>
						<span class="muted">{$t.view.delete_confirm}</span>
						<button type="button" class="btn-text-danger font-semibold" on:click={deleteNote} disabled={deleting}>
							{#if deleting}<Spinner />{/if}
							{$t.view.delete_yes}
						</button>
						<button type="button" class="btn-text" on:click={() => (confirmingDelete = false)}>{$t.view.cancel}</button>
					</span>
				{:else}
					<button type="button" class="btn-text-danger" on:click={() => (confirmingDelete = true)}>
						{$t.view.delete}
					</button>
				{/if}
			{/if}
		</div>
	{/if}
</section>
