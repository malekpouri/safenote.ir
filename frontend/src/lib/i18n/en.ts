export const en = {
	app: {
		title: 'SafeNote',
		tagline: 'Share encrypted notes that self-destruct.',
		description:
			'SafeNote is a free, open-source service for sending confidential notes. Notes are encrypted in your browser, can only be read with the link, and are destroyed after reading.',
		skip_to_content: 'Skip to content',
		noscript: 'SafeNote encrypts notes in your browser, so JavaScript is required.'
	},
	nav: {
		switch_language: 'فارسی',
		switch_language_label: 'تغییر زبان به فارسی'
	},
	footer: {
		about: 'About',
		privacy: 'Privacy',
		stats: 'Stats',
		source: 'GitHub',
		powered: 'Powered by'
	},
	home: {
		title_lead: 'Share secrets',
		title_accent: 'that vanish.',
		subtitle: 'Encrypted in your browser and destroyed after it’s read.',
		note_label: 'Secret note',
		placeholder: 'Write something private…',
		chars: '{n} / {max}',
		options: 'Options',
		views_label: 'Destroy after',
		views_option: '{n} view',
		views_option_plural: '{n} views',
		expires_label: 'Expires in',
		hours_1: '1 hour',
		hours_24: '24 hours',
		days_7: '7 days',
		days_30: '30 days',
		password_label: 'Password',
		password_optional: 'optional',
		password_placeholder: 'No password',
		password_set: 'password',
		password_show: 'Show password',
		password_hide: 'Hide password',
		password_generate: 'Generate a strong password',
		password_hint: 'Send the password to the recipient separately from the link.',
		submit: 'Create link',
		submitting: 'Encrypting…',
		footnote: 'End-to-end encrypted. We can’t read your note.',
		learn_more: 'How it works'
	},
	success: {
		title: 'Your link is ready',
		message: 'Copy it now. It contains the key and won’t be shown again.',
		link_label: 'Secure link',
		copy: 'Copy link',
		copied: 'Copied',
		share: 'Share',
		summary_views: 'Destroyed after {n} view',
		summary_views_plural: 'Destroyed after {n} views',
		summary_expires: 'expires {date}',
		summary_password: 'Password protected. Send the password separately.',
		new_note: 'New note',
		delete: 'Delete note',
		delete_confirm: 'Delete permanently?',
		delete_yes: 'Delete',
		cancel: 'Cancel'
	},
	how: {
		title: 'How it works',
		step1_title: 'Write & encrypt',
		step1_body: 'Your note is encrypted with AES-256-GCM in your browser using a random key that never leaves your device.',
		step2_title: 'Share the link',
		step2_body: 'The link holds the note ID and, after the “#”, the key. Browsers never send that part to the server.',
		step3_title: 'Read once, then gone',
		step3_body: 'The recipient’s browser decrypts the note, and the server deletes it permanently.'
	},
	faq: {
		title: 'Frequently asked questions',
		q1: 'Why not just send the secret over chat or email?',
		a1: 'Messages in chat apps and mailboxes stay around for years and are copied to backups and other devices. A SafeNote link leaves nothing readable behind once it has been opened.',
		q2: 'Can SafeNote read my notes?',
		a2: 'No. The decryption key is in the part of the link after “#”, which browsers never send to a server. We only store encrypted data that is useless without it.',
		q3: 'What happens if I lose the link?',
		a3: 'The note cannot be recovered by anyone, including us. Just create a new one.',
		q4: 'What if someone else opens the link first?',
		a4: 'The note will be destroyed and the real recipient will see that it no longer exists, so you will know it was intercepted.'
	},
	view: {
		page_title: 'Secret note',
		loading: 'Checking the note…',
		missing_key_title: 'This link is incomplete',
		missing_key_body: 'The decryption key after “#” is missing. Ask the sender for the full link.',
		not_found_title: 'Note not available',
		not_found_body: 'This note doesn’t exist, has already been read, or has expired.',
		ready_title: 'You received a secret note',
		ready_body: 'It may be destroyed as soon as you open it.',
		password_title: 'This note needs a password',
		password_label: 'Password',
		password_placeholder: 'Enter password',
		wrong_password: 'Wrong password. The note wasn’t opened, so you can try again.',
		wrong_password_legacy: 'Wrong password. Try again without leaving this page.',
		open: 'Open note',
		opening: 'Decrypting…',
		revealed_title: 'Secret note',
		burned: 'This note has been destroyed. Copy it if you need it.',
		views_left: 'Can be viewed {n} more time.',
		views_left_plural: 'Can be viewed {n} more times.',
		expires: 'Expires {date}',
		copy: 'Copy',
		new_note: 'Create a note',
		delete: 'Delete note',
		delete_confirm: 'Delete permanently?',
		delete_yes: 'Delete',
		cancel: 'Cancel',
		deleted_title: 'Note deleted',
		deleted_body: 'The note has been permanently removed from the server.'
	},
	toast: {
		copied: 'Copied to clipboard',
		copy_failed: 'Could not copy. Select the text and copy it manually.',
		deleted: 'Note deleted',
		error: 'Something went wrong. Please try again.',
		rate_limited: 'Too many requests. Please wait a minute and try again.',
		too_long: 'The note is too long.',
		decrypt_failed: 'Could not decrypt this note. The link may be damaged.'
	},
	about: {
		title: 'About SafeNote',
		lead: 'SafeNote lets you send passwords, keys and other sensitive information without leaving a readable copy in chat histories, mailboxes or backups.',
		contact_title: 'Contact',
		contact_body: 'Questions or feedback? Email',
		contact_or: 'or open an issue on'
	},
	privacy: {
		title: 'Privacy policy',
		updated: 'Last updated: October 2026',
		intro: 'Privacy is the reason SafeNote exists. This policy explains what happens to your data when you use the service.',
		sections: [
			{
				h: '1. Service description',
				p: 'SafeNote lets you create encrypted notes and share them as unique links. Sending the link is up to you. Depending on the channel you choose (email, SMS, messengers, social media), third parties could intercept the link and read the note before the intended recipient.'
			},
			{
				h: '2. Encryption and what we store',
				p: 'Notes are encrypted in your browser with AES-256-GCM. The key is generated on your device and placed in the part of the link after “#”, which browsers never send to servers. The key, together with your password if you set one, is hardened on your device with 600,000 rounds of PBKDF2; neither is ever transmitted. Our server stores only the encrypted note, its expiry time, its remaining view count, a random salt, and a one-way hash of an access token derived from the key. None of this can be used to decrypt the note.'
			},
			{
				h: '3. Deletion',
				p: 'A note is permanently deleted as soon as its last allowed view is used, when it expires (at most 30 days after creation), or when someone holding the link deletes it. Deleted notes cannot be recovered by anyone.'
			},
			{
				h: '4. IP addresses and logs',
				p: 'We do not log IP addresses or link them to notes. IP addresses are held briefly in memory only to rate-limit abuse and are discarded within minutes.'
			},
			{
				h: '5. Cookies and third parties',
				p: 'SafeNote sets a single cookie that remembers your language choice. We use no analytics, advertising, tracking pixels, or third-party fonts or scripts.'
			},
			{
				h: '6. Pseudonymous data',
				p: 'A note may contain personal data entered by its author. Because we never hold the key, we cannot access that data, and our database does not reveal who created a note.'
			},
			{
				h: '7. Disclaimer',
				p: 'SafeNote is not responsible for the content of notes. Do not use the service for unlawful purposes.'
			},
			{
				h: '8. Children',
				p: 'SafeNote is not intended for children under 16. Minors should obtain consent from a parent or guardian before using it.'
			},
			{
				h: '9. Changes to this policy',
				p: 'We may update this policy from time to time. Significant changes will be announced on the home page.'
			}
		],
		contact_title: '10. Contact',
		contact_body: 'Questions about this policy? Email us at'
	},
	admin: {
		title: 'Stats',
		subtitle: 'Live, anonymous counts.',
		active: 'Active notes',
		total: 'Notes created',
		refresh: 'Refresh',
		updated: 'Updated {time}',
		error: 'Stats are unavailable right now.'
	},
	error: {
		not_found_title: 'Page not found',
		not_found_body: 'The page you are looking for does not exist.',
		generic_title: 'Something went wrong',
		generic_body: 'An unexpected error occurred. Please try again.',
		home: 'Back to home'
	}
};

export type Dictionary = typeof en;
