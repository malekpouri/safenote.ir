import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

const API_URL = env.INTERNAL_API_URL || 'http://backend:8080';

export async function GET({ params, fetch }) {
    try {
        const response = await fetch(`${API_URL}/api/notes/${params.id}`);

        if (!response.ok) {
            const errData = await response.json().catch(() => ({}));
            if (response.status === 404) {
                throw error(404, errData.error || 'Note not found');
            }
            throw error(response.status, errData.error || 'Backend request failed');
        }

        const data = await response.json();
        return new Response(JSON.stringify(data), {
            headers: { 'Content-Type': 'application/json' }
        });
    } catch (e: any) {
        console.error('Proxy error:', e);
        if (e.status) {
            throw e;
        }
        return new Response(JSON.stringify({ error: e.message || 'Internal Server Error' }), {
            status: 500,
            headers: { 'Content-Type': 'application/json' }
        });
    }
}

export async function DELETE({ params, request, fetch }) {
    try {
        const body = await request.json().catch(() => ({}));
        const response = await fetch(`${API_URL}/api/notes/${params.id}`, {
            method: 'DELETE',
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
        if (e.status) {
            throw e;
        }
        return new Response(JSON.stringify({ error: e.message || 'Internal Server Error' }), {
            status: 500,
            headers: { 'Content-Type': 'application/json' }
        });
    }
}
