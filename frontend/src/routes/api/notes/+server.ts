import { proxy } from '$lib/server/backend';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = (event) => proxy(event, '/api/notes');
