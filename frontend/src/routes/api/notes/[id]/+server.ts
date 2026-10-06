import { proxy } from '$lib/server/backend';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = (event) => proxy(event, `/api/notes/${encodeURIComponent(event.params.id)}`);
export const DELETE: RequestHandler = (event) => proxy(event, `/api/notes/${encodeURIComponent(event.params.id)}`);
