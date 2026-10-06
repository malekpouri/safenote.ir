// Run with: npm test  (Node's built-in test runner with native TypeScript support)
import { test } from 'node:test';
import assert from 'node:assert/strict';
import { createCipheriv, createHash, randomBytes } from 'node:crypto';
import {
	NOTE_MAX_CHARS,
	V2_PREFIX,
	decryptNote,
	deriveKeys,
	encryptNote,
	generateLinkKey,
	generatePassword,
	generateSalt,
	isV2Key,
	legacyPasswordHash
} from '../src/lib/crypto.ts';

const SERVER_MAX_ENCRYPTED_LENGTH = 60000; // backend/internal/api/controllers.MaxEncryptedLength

test('link keys are 10 random base64url characters', () => {
	const keys = new Set(Array.from({ length: 200 }, generateLinkKey));
	assert.equal(keys.size, 200);
	for (const key of keys) {
		assert.match(key, /^[A-Za-z0-9_-]{10}$/);
		assert.ok(isV2Key(key));
	}
	assert.ok(!isV2Key('abc123')); // legacy v1 short key
	assert.match(generateSalt(), /^[A-Za-z0-9_-]{22}$/);
});

test('v2 round trip without password', async () => {
	const key = generateLinkKey();
	const { payload, accessToken, salt } = await encryptNote('hello سلام 👋', key);
	assert.ok(payload.startsWith(V2_PREFIX));
	assert.match(accessToken, /^[A-Za-z0-9_-]{43}$/);
	assert.equal(await decryptNote(payload, key, '', { salt }), 'hello سلام 👋');
});

test('v2 round trip with password; wrong password, key or salt fail', async () => {
	const key = generateLinkKey();
	const { payload, accessToken, salt } = await encryptNote('secret', key, 'pa55');
	assert.equal(await decryptNote(payload, key, 'pa55', { salt }), 'secret');
	await assert.rejects(decryptNote(payload, key, 'wrong', { salt }));
	await assert.rejects(decryptNote(payload, generateLinkKey(), 'pa55', { salt }));
	await assert.rejects(decryptNote(payload, key, 'pa55', { salt: generateSalt() }));

	// The access token depends on the link key, the password and the salt.
	assert.equal((await deriveKeys(key, 'pa55', salt)).accessToken, accessToken);
	assert.notEqual((await deriveKeys(key, 'wrong', salt)).accessToken, accessToken);
	assert.notEqual((await deriveKeys(key, '', salt)).accessToken, accessToken);
	assert.notEqual((await deriveKeys(key, 'pa55', generateSalt())).accessToken, accessToken);
});

test('pre-derived keys decrypt without re-running PBKDF2', async () => {
	const key = generateLinkKey();
	const { payload, salt } = await encryptNote('fast path', key);
	const keys = await deriveKeys(key, '', salt);
	assert.equal(await decryptNote(payload, key, '', { keys }), 'fast path');
});

test('legacy v1 payloads produced by the old Go/WASM code still decrypt', async () => {
	const shortKey = 'aB3xY9';
	const password = 'pw';
	const aesKey = createHash('sha256').update(shortKey + password).digest();
	const nonce = randomBytes(12);
	const cipher = createCipheriv('aes-256-gcm', aesKey, nonce);
	const ct = Buffer.concat([cipher.update('legacy note', 'utf8'), cipher.final(), cipher.getAuthTag()]);
	const payload = Buffer.concat([nonce, ct]).toString('base64');

	assert.equal(await decryptNote(payload, shortKey, password), 'legacy note');
	await assert.rejects(decryptNote(payload, shortKey, 'nope'));
});

test('legacy password hash matches the old hex SHA-256', async () => {
	assert.equal(await legacyPasswordHash('abc'), createHash('sha256').update('abc').digest('hex'));
});

test('a full-length note in any script fits the server limit', async () => {
	const key = generateLinkKey();
	for (const ch of ['a', 'س', '中']) {
		const { payload } = await encryptNote(ch.repeat(NOTE_MAX_CHARS), key);
		assert.ok(payload.length <= SERVER_MAX_ENCRYPTED_LENGTH, `${ch}: ${payload.length}`);
	}
});

test('generated passwords are random and of the requested length', () => {
	const a = generatePassword();
	assert.equal(a.length, 16);
	assert.notEqual(generatePassword(), a);
	assert.equal(generatePassword(24).length, 24);
});
