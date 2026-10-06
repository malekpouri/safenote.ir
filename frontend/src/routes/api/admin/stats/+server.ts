import { proxy } from '$lib/server/backend';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = (event) => proxy(event, '/api/admin/stats');
