import type { Dictionary } from '$lib/i18n/en';

/** Maps an API failure to a localized, user-facing message. */
export function apiErrorMessage(status: number, t: Dictionary, fallback?: string): string {
	if (status === 429) return t.toast.rate_limited;
	if (status === 413) return t.toast.too_long;
	if (status === 404) return t.view.not_found_body;
	if (status >= 500) return t.toast.error;
	return fallback || t.toast.error;
}
