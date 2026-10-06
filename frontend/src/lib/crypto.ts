/**
 * Client-side encryption for SafeNote, built on the WebCrypto API.
 *
 * v2 scheme:
 *   linkKey    = 10 random base64url chars (60 bits) in the URL fragment (never sent to the server)
 *   salt       = 16 random bytes per note, stored by the server (not secret)
 *   master     = PBKDF2-SHA256(linkKey || password, salt = "safenote/v2" || salt, 600k)
 *   encKey     = HKDF-SHA256(master, info = "safenote/v2/enc")    -> AES-256-GCM
 *   accessTok  = HKDF-SHA256(master, info = "safenote/v2/access") -> sent to the server, which stores SHA-256(token)
 *   payload    = "v2:" || base64url(iv[12] || ciphertext || tag)
 *
 * The link key is short so links fit in an SMS. Stretching it with PBKDF2
 * makes every guess cost 600k hash rounds, and the per-note salt rules out
 * precomputation, so even someone holding the database cannot brute-force a
 * note. Everyone else cannot get the ciphertext at all: opening a note needs
 * the access token, which also derives from the key.
 *
 * v1 (legacy) payloads have no prefix:
 *   key = SHA-256(shortKey + password), payload = base64(nonce[12] || ciphertext || tag)
 */

export const V2_PREFIX = 'v2:';
export const PBKDF2_ITERATIONS = 600_000;
export const NOTE_MAX_CHARS = 10_000;

const LINK_KEY_CHARS = 10;
const SALT_BYTES = 16;
const IV_BYTES = 12;
const BASE64URL = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_';
const enc = new TextEncoder();
const dec = new TextDecoder();
const subtle = () => globalThis.crypto.subtle;
type Bytes = Uint8Array<ArrayBuffer>;

export function toBase64Url(bytes: Uint8Array): string {
	let bin = '';
	for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i]);
	return btoa(bin).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

export function fromBase64Url(s: string): Bytes {
	const b64 = s.replace(/-/g, '+').replace(/_/g, '/') + '='.repeat((4 - (s.length % 4)) % 4);
	const bin = atob(b64);
	const out = new Uint8Array(bin.length);
	for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
	return out;
}

function concat(...parts: Uint8Array[]): Bytes {
	const out = new Uint8Array(parts.reduce((n, p) => n + p.length, 0));
	let offset = 0;
	for (const p of parts) {
		out.set(p, offset);
		offset += p.length;
	}
	return out;
}

/** A fresh link key: 10 uniformly random base64url characters (60 bits). */
export function generateLinkKey(): string {
	// 64 symbols, so the low 6 bits of each random byte map uniformly.
	return Array.from(globalThis.crypto.getRandomValues(new Uint8Array(LINK_KEY_CHARS)), (b) => BASE64URL[b & 63]).join('');
}

/** A fresh per-note KDF salt (base64url, 22 characters). */
export function generateSalt(): string {
	return toBase64Url(globalThis.crypto.getRandomValues(new Uint8Array(SALT_BYTES)));
}

/** v2 link keys are 10 base64url chars; anything else is a legacy v1 key. */
export function isV2Key(key: string): boolean {
	return /^[A-Za-z0-9_-]{10}$/.test(key);
}

export interface NoteKeys {
	encKey: CryptoKey;
	accessToken: string;
}

export async function deriveKeys(linkKey: string, password: string, salt: string): Promise<NoteKeys> {
	if (!isV2Key(linkKey)) throw new Error('Invalid link key');
	const saltBytes = fromBase64Url(salt);
	if (saltBytes.length !== SALT_BYTES) throw new Error('Invalid salt');

	// The link key has a fixed length, so linkKey || password is unambiguous.
	const secret = await subtle().importKey('raw', enc.encode(linkKey + password), 'PBKDF2', false, ['deriveBits']);
	const master = await subtle().deriveBits(
		{ name: 'PBKDF2', hash: 'SHA-256', salt: concat(enc.encode('safenote/v2'), saltBytes), iterations: PBKDF2_ITERATIONS },
		secret,
		256
	);

	const ikm = await subtle().importKey('raw', master, 'HKDF', false, ['deriveBits', 'deriveKey']);
	const hkdf = (info: string) => ({ name: 'HKDF', hash: 'SHA-256', salt: new Uint8Array(0), info: enc.encode(info) });

	const encKey = await subtle().deriveKey(hkdf('safenote/v2/enc'), ikm, { name: 'AES-GCM', length: 256 }, false, [
		'encrypt',
		'decrypt'
	]);
	const tokenBits = new Uint8Array(await subtle().deriveBits(hkdf('safenote/v2/access'), ikm, 256));
	return { encKey, accessToken: toBase64Url(tokenBits) };
}

export interface EncryptedNote {
	payload: string;
	accessToken: string;
	salt: string;
}

export async function encryptNote(text: string, linkKey: string, password = ''): Promise<EncryptedNote> {
	const salt = generateSalt();
	const { encKey, accessToken } = await deriveKeys(linkKey, password, salt);
	const iv = globalThis.crypto.getRandomValues(new Uint8Array(IV_BYTES));
	const ct = new Uint8Array(await subtle().encrypt({ name: 'AES-GCM', iv }, encKey, enc.encode(text)));
	return { payload: V2_PREFIX + toBase64Url(concat(iv, ct)), accessToken, salt };
}

async function decryptV2(payload: string, encKey: CryptoKey): Promise<string> {
	const data = fromBase64Url(payload.slice(V2_PREFIX.length));
	if (data.length < IV_BYTES + 16) throw new Error('Ciphertext too short');
	const pt = await subtle().decrypt({ name: 'AES-GCM', iv: data.slice(0, IV_BYTES) }, encKey, data.slice(IV_BYTES));
	return dec.decode(pt);
}

async function decryptV1(payload: string, shortKey: string, password: string): Promise<string> {
	const digest = await subtle().digest('SHA-256', enc.encode(shortKey + password));
	const key = await subtle().importKey('raw', digest, 'AES-GCM', false, ['decrypt']);
	const bin = atob(payload);
	const data = new Uint8Array(bin.length);
	for (let i = 0; i < bin.length; i++) data[i] = bin.charCodeAt(i);
	if (data.length < IV_BYTES + 16) throw new Error('Ciphertext too short');
	const pt = await subtle().decrypt({ name: 'AES-GCM', iv: data.slice(0, IV_BYTES) }, key, data.slice(IV_BYTES));
	return dec.decode(pt);
}

/**
 * Decrypts either format. v2 needs either pre-derived keys (avoids running
 * PBKDF2 twice) or the note's salt. Throws on a wrong key or password.
 */
export async function decryptNote(
	payload: string,
	linkKey: string,
	password = '',
	v2: { keys?: NoteKeys; salt?: string } = {}
): Promise<string> {
	if (payload.startsWith(V2_PREFIX)) {
		if (!v2.keys && !v2.salt) throw new Error('Missing salt');
		const { encKey } = v2.keys ?? (await deriveKeys(linkKey, password, v2.salt!));
		return decryptV2(payload, encKey);
	}
	return decryptV1(payload, linkKey, password);
}

/** Hex SHA-256, used only to authorize deleting legacy v1 notes. */
export async function legacyPasswordHash(password: string): Promise<string> {
	const digest = new Uint8Array(await subtle().digest('SHA-256', enc.encode(password)));
	return Array.from(digest, (b) => b.toString(16).padStart(2, '0')).join('');
}

const PASSWORD_ALPHABET = 'abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#$%&*?';

/** Uniformly random password without look-alike characters. */
export function generatePassword(length = 16): string {
	const n = PASSWORD_ALPHABET.length;
	const limit = 256 - (256 % n);
	let out = '';
	while (out.length < length) {
		for (const b of globalThis.crypto.getRandomValues(new Uint8Array(length * 2))) {
			if (b < limit) out += PASSWORD_ALPHABET[b % n];
			if (out.length === length) break;
		}
	}
	return out;
}
