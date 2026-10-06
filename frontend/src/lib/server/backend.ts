import { env } from '$env/dynamic/private';
import type { RequestEvent } from '@sveltejs/kit';

const API_URL = env.INTERNAL_API_URL || 'http://backend:8080';

/**
 * Best-effort client IP for the backend's rate limiter. Behind a reverse proxy
 * the last X-Forwarded-For entry is the one the proxy appended.
 */
function clientIp(event: RequestEvent): string {
	const xff = event.request.headers.get('x-forwarded-for');
	if (xff) {
		const last = xff.split(',').map((s) => s.trim()).filter(Boolean).pop();
		if (last) return last;
	}
	const real = event.request.headers.get('x-real-ip');
	if (real) return real;
	try {
		return event.getClientAddress();
	} catch {
		return '';
	}
}

/**
 * Forwards a request to the Go backend and passes its status and JSON body
 * through unchanged.
 */
export async function proxy(event: RequestEvent, path: string, method = event.request.method): Promise<Response> {
	const headers: Record<string, string> = { 'Content-Type': 'application/json' };
	const ip = clientIp(event);
	if (ip) headers['X-Real-IP'] = ip;

	const hasBody = method !== 'GET' && method !== 'HEAD';
	try {
		const res = await event.fetch(`${API_URL}${path}`, {
			method,
			headers,
			body: hasBody ? await event.request.text() : undefined
		});
		return new Response(await res.text(), {
			status: res.status,
			headers: { 'Content-Type': 'application/json', 'Cache-Control': 'no-store' }
		});
	} catch (e) {
		console.error('Backend unreachable:', e);
		return new Response(JSON.stringify({ error: 'Service temporarily unavailable' }), {
			status: 502,
			headers: { 'Content-Type': 'application/json', 'Cache-Control': 'no-store' }
		});
	}
}
