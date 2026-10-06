import defaultTheme from 'tailwindcss/defaultTheme';

/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'media',
	theme: {
		extend: {
			fontFamily: {
				sans: ['"Inter Variable"', 'Vazirmatn', ...defaultTheme.fontFamily.sans],
				fa: ['Vazirmatn', ...defaultTheme.fontFamily.sans]
			},
			colors: {
				brand: {
					50: '#eef2ff',
					100: '#e0e7ff',
					200: '#c7d2fe',
					300: '#a5b4fc',
					400: '#818cf8',
					500: '#6366f1',
					600: '#4f46e5',
					700: '#4338ca',
					800: '#3730a3',
					900: '#312e81',
					950: '#1e1b4b'
				}
			},
			boxShadow: {
				card: '0 1px 2px rgb(15 23 42 / 0.04), 0 8px 24px -6px rgb(15 23 42 / 0.10)'
			},
			keyframes: {
				pop: {
					'0%': { opacity: '0', transform: 'scale(0.6)' },
					'60%': { opacity: '1', transform: 'scale(1.08)' },
					'100%': { transform: 'scale(1)' }
				},
				'fade-up': {
					'0%': { opacity: '0', transform: 'translateY(6px)' },
					'100%': { opacity: '1', transform: 'translateY(0)' }
				}
			},
			animation: {
				'fade-up': 'fade-up 0.35s ease-out both',
				pop: 'pop 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) both'
			}
		}
	},
	plugins: []
};
