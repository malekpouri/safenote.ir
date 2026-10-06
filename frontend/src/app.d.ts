// See https://kit.svelte.dev/docs/types#app
import type { Locale } from '$lib/i18n';

declare global {
	namespace App {
		interface Locals {
			locale: Locale;
		}
		interface PageData {
			locale: Locale;
		}
	}
}

export {};
