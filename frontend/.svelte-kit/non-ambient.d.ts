
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	export interface AppTypes {
		RouteId(): "/" | "/about" | "/api" | "/api/notes" | "/api/notes/[id]" | "/n" | "/n/[id]" | "/privacy";
		RouteParams(): {
			"/api/notes/[id]": { id: string };
			"/n/[id]": { id: string }
		};
		LayoutParams(): {
			"/": { id?: string };
			"/about": Record<string, never>;
			"/api": { id?: string };
			"/api/notes": { id?: string };
			"/api/notes/[id]": { id: string };
			"/n": { id?: string };
			"/n/[id]": { id: string };
			"/privacy": Record<string, never>
		};
		Pathname(): "/" | "/about" | "/about/" | "/api" | "/api/" | "/api/notes" | "/api/notes/" | `/api/notes/${string}` & {} | `/api/notes/${string}/` & {} | "/n" | "/n/" | `/n/${string}` & {} | `/n/${string}/` & {} | "/privacy" | "/privacy/";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): "/favicon.png" | "/robots.txt" | "/sitemap.xml" | string & {};
	}
}