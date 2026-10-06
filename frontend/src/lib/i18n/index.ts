import { derived, writable } from 'svelte/store';
import { en } from './en';
import { fa } from './fa';
import { LOCALE_COOKIE, dirOf, type Locale } from './locale';

export type { Locale } from './locale';

const translations = { en, fa };

function createLocaleStore() {
	const { subscribe, set } = writable<Locale>('en');

	function apply(locale: Locale, persist: boolean) {
		set(locale);
		if (typeof document === 'undefined') return;
		document.documentElement.lang = locale;
		document.documentElement.dir = dirOf(locale);
		if (persist) {
			document.cookie = `${LOCALE_COOKIE}=${locale}; path=/; max-age=31536000; samesite=lax`;
		}
	}

	return {
		subscribe,
		/** Sync with the server-resolved locale (no cookie write). */
		init: (locale: Locale) => apply(locale, false),
		/** User-initiated change; remembered via cookie so SSR matches next time. */
		set: (locale: Locale) => apply(locale, true)
	};
}

export const locale = createLocaleStore();

export const t = derived(locale, ($locale) => translations[$locale]);

export const dir = derived(locale, dirOf);

/** Replaces {name} placeholders. */
export function fmt(template: string, vars: Record<string, string | number>): string {
	return template.replace(/\{(\w+)\}/g, (_, k) => String(vars[k] ?? `{${k}}`));
}

/** Locale-aware number formatting (Persian digits for fa). */
export const num = derived(locale, ($locale) => {
	const f = new Intl.NumberFormat($locale === 'fa' ? 'fa-IR' : 'en-US');
	return (n: number) => f.format(n);
});

export const dateTime = derived(locale, ($locale) => {
	const f = new Intl.DateTimeFormat($locale === 'fa' ? 'fa-IR' : 'en-US', {
		dateStyle: 'medium',
		timeStyle: 'short'
	});
	return (d: Date | string) => f.format(new Date(d));
});
