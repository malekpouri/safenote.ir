import type { Handle } from '@sveltejs/kit';
import { LOCALE_COOKIE, resolveLocale } from '$lib/i18n/locale';

const securityHeaders: Record<string, string> = {
	'Referrer-Policy': 'no-referrer',
	'X-Content-Type-Options': 'nosniff',
	'X-Frame-Options': 'DENY',
	'Permissions-Policy': 'camera=(), microphone=(), geolocation=(), interest-cohort=()',
	'Cross-Origin-Opener-Policy': 'same-origin'
};

export const handle: Handle = async ({ event, resolve }) => {
	const locale = resolveLocale(event.cookies.get(LOCALE_COOKIE), event.request.headers.get('accept-language'));
	event.locals.locale = locale;

	const response = await resolve(event, {
		transformPageChunk: ({ html }) =>
			html.replace('%lang%', locale).replace('%dir%', locale === 'fa' ? 'rtl' : 'ltr')
	});

	for (const [k, v] of Object.entries(securityHeaders)) response.headers.set(k, v);
	const path = event.url.pathname;
	const isNotePage = path.startsWith('/n/') || /^\/[A-Za-z0-9]{6}$/.test(path);
	if (isNotePage || path.startsWith('/api/')) {
		response.headers.set('Cache-Control', 'no-store');
		response.headers.set('X-Robots-Tag', 'noindex, nofollow');
	}
	return response;
};
