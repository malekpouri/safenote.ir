export type Locale = 'en' | 'fa';

export const LOCALES: Locale[] = ['en', 'fa'];
export const LOCALE_COOKIE = 'locale';

export function isLocale(value: unknown): value is Locale {
	return value === 'en' || value === 'fa';
}

/** Cookie first, then the browser's Accept-Language, defaulting to English. */
export function resolveLocale(cookie: string | undefined | null, acceptLanguage: string | null): Locale {
	if (isLocale(cookie)) return cookie;
	const preferred = (acceptLanguage ?? '')
		.split(',')
		.map((part) => part.split(';')[0].trim().toLowerCase())
		.find((tag) => tag.startsWith('fa') || tag.startsWith('en'));
	return preferred?.startsWith('fa') ? 'fa' : 'en';
}

export const dirOf = (locale: Locale) => (locale === 'fa' ? 'rtl' : 'ltr');
