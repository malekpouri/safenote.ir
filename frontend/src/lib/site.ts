export const SITE_URL = 'https://www.safenote.ir';
export const SITE_HOST = 'www.safenote.ir';

/**
 * Origin used in shared links. The bare domain is 4 characters shorter (handy
 * for SMS). It 301-redirects to www, and browsers keep the #key across the
 * redirect.
 */
export const SHARE_ORIGIN = 'https://safenote.ir';

/** SafeNote's own domain gets the short origin; self-hosted copies keep theirs. */
export function shareOrigin(location: Location): string {
	return /(^|\.)safenote\.ir$/.test(location.hostname) ? SHARE_ORIGIN : location.origin;
}

export const GITHUB_URL = 'https://github.com/malekpouri/safenote.ir';
export const CONTACT_EMAIL = 'admin@utux.ir';
export const POWERED_BY = { name: 'Utux', url: 'https://utux.ir' };

export const VIEW_OPTIONS = [1, 2, 3, 5, 10] as const;
/** Minutes. */
export const EXPIRATION_OPTIONS = [60, 1440, 10080, 43200] as const;
export const DEFAULT_EXPIRATION = 43200;
