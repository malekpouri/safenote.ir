import { writable, derived } from 'svelte/store';
import { en } from './en';
import { fa } from './fa';

export type Locale = 'en' | 'fa';

const translations = {
    en,
    fa
};

function createI18nStore() {
    const { subscribe, set, update } = writable<Locale>('en');

    return {
        subscribe,
        set: (locale: Locale) => {
            set(locale);
            if (typeof document !== 'undefined') {
                document.documentElement.lang = locale;
                document.documentElement.dir = locale === 'fa' ? 'rtl' : 'ltr';
            }
        },
        toggle: () => update(l => {
            const next = l === 'en' ? 'fa' : 'en';
            if (typeof document !== 'undefined') {
                document.documentElement.lang = next;
                document.documentElement.dir = next === 'fa' ? 'rtl' : 'ltr';
            }
            return next;
        })
    };
}

export const locale = createI18nStore();

export const t = derived(locale, ($locale) => {
    return translations[$locale];
});

export const dir = derived(locale, ($locale) => {
    return $locale === 'fa' ? 'rtl' : 'ltr';
});
