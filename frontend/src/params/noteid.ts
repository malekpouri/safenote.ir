import type { ParamMatcher } from '@sveltejs/kit';

/** Short note links: /<6 base62 chars>. Static routes (/about, …) still win. */
export const match: ParamMatcher = (param) => /^[A-Za-z0-9]{6}$/.test(param);
