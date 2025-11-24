import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

const API_URL = env.INTERNAL_API_URL || 'http://backend:8080';

export async function POST({ request, fetch }) {
    const body = await request.json();

    try {
        const response = await fetch(`${API_URL}/api/notes`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(body)
        });

        if (!response.ok) {
            const errData = await response.json().catch(() => ({}));
            throw error(response.status, errData.error || 'Backend request failed');
        }

        const data = await response.json();
        return new Response(JSON.stringify(data), {
            headers: { 'Content-Type': 'application/json' }
        });
    } catch (e: any) {
        console.error('Proxy error:', e);
        return new Response(JSON.stringify({ error: e.message || 'Internal Server Error' }), {
            status: 500,
            headers: { 'Content-Type': 'application/json' }
        });
    }
}
