# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

SafeNote (https://www.safenote.ir) shares encrypted notes that delete themselves. Encryption happens in the browser, so the server only ever stores ciphertext. The UI is bilingual: English, and Persian (fa) with an RTL layout. Planned and finished work is tracked in `ROADMAP.md`.

## Commands

Full stack (served on host port **8081**; the backend is not exposed to the host):
```bash
docker compose up --build
```

Backend (Go + Fiber + mattn/go-sqlite3, needs CGO and a C compiler):
```bash
cd backend
DB_PATH=./dev.db go run ./cmd/server      # env: DB_PATH, PORT, RATE_LIMIT_CREATE, RATE_LIMIT_READ
go test ./...                             # all tests live in internal/api/router_test.go
go test ./internal/api -run TestConcurrentOpenOnlyOnce -race
```

Frontend (SvelteKit 2 + Svelte 4 + Tailwind 3, adapter-node):
```bash
cd frontend
INTERNAL_API_URL=http://localhost:8080 npm run dev
npm run check     # svelte-check (types in .svelte and .ts)
npm test          # node --test 'tests/*.test.ts' (Node's native TS support, no test framework)
npm run build
```

The build prints warnings that SvelteKit imports `untrack`, `fork` and `settled` from svelte. This is a known, harmless Kit 2.49 / Svelte 4 interaction.

## Architecture

**Request path:** browser → SvelteKit server routes (`frontend/src/routes/api/**/+server.ts`) → Go backend.
- Every API route is a one-line call to `proxy()` in `src/lib/server/backend.ts`. It passes the backend's status and JSON through unchanged and forwards the client IP as `X-Real-IP`, which the backend rate limiter trusts.
- A new backend endpoint needs a matching proxy route.

**Encryption (`frontend/src/lib/crypto.ts`, the core of the design):**
- `linkKey` is 10 random base64url chars (60 bits) in the URL fragment. Links are `/<id>#<linkKey>` with a 6-char ID (`src/params/noteid.ts`). The older `/n/<id>` route still works. Both render `src/lib/components/NoteViewer.svelte`. Links are kept short on purpose (SMS).
- `master = PBKDF2-SHA256(linkKey ‖ password, salt = "safenote/v2" ‖ salt, 600k)`. The 16-byte salt is random per note; the server stores it in `kdf_salt` and returns it from `GET /api/notes/:id`.
- The stretching is what makes the short key safe: it makes offline brute force by someone holding the database impractical. Don't lower the iterations or the key length without redoing that math.
- HKDF over `master` derives two values:
  - the AES-256-GCM `encKey`;
  - an `accessToken`, which the server stores only as SHA-256.
- Ciphertext is stored as `"v2:" + base64url(iv ‖ ct ‖ tag)`.
- The access token gates `POST /api/notes/:id/open` and `DELETE`. A wrong password therefore returns 401 **before** a view is consumed.
- Changing any derivation constant breaks every existing link.

**Legacy v1 notes:**
- Notes created before the redesign have no `v2:` prefix and an empty `access_token_hash`.
- Their key is `SHA-256(shortKey + password)`.
- They open without a token, and are deleted using the old `password_hash`.
- `meta.legacy` from `GET /api/notes/:id` tells the viewer which path to take.
- Roadmap item 1.10 removes this path once old notes have expired.

**Backend (`backend/internal/`):**
- `api/router.go` builds the Fiber app: limiter, routes, and `/health`.
- `api/controllers` validates input (views 1–100, expiration 5 min–30 days, `MaxEncryptedLength`).
- `repositories`: SQLite in WAL mode with `SetMaxOpenConns(1)`. The single connection is what makes `Consume` (read → authorize → decrement → delete on the last view) atomic. Keep it.
- Schema changes happen in `migrate()` (`CREATE IF NOT EXISTS` plus `ALTER TABLE` when a column is missing). There is no migration tool.
- Times are stored in UTC. Liveness checks happen in Go (`isLive`), not in SQL.
- `scheduler` purges expired notes at startup and every 10 minutes.

**Frontend:**
- **Locale:** `hooks.server.ts` resolves it (cookie `locale`, then `Accept-Language`) and rewrites `%lang%`/`%dir%` in `app.html`. `+layout.svelte` calls `locale.init(data.locale)`.
- **Strings:** every UI string goes in **both** `src/lib/i18n/en.ts` and `fa.ts`. `fa` is typed as `Dictionary`, so a missing key fails `npm run check`. Use `fmt()` for `{placeholders}` and the `$num` / `$dateTime` stores for locale-aware formatting.
- **CSP:** configured in `svelte.config.js` (`kit.csp`). No external origins are allowed: fonts are self-hosted via `@fontsource-variable/inter` and `vazirmatn`. Adding any CDN or external asset requires a CSP change.
- **Design:** deliberately minimal. Every page is one narrow column (`container-narrow`). Personality comes from a few deliberate touches: the logo, brand-colored buttons, one soft background glow, a gradient accent in the title, and option pills (`chip`). Avoid adding sections, grids or extra navigation. Shared classes (`page-title`, `input`, `btn-primary`, `btn-text`, `card`, `chip`, `text-gradient`, `link`, `muted`, `prose-page`…) live in `src/app.css` under `@layer components`. Keep form fields at 16px on phones (iOS zoom) and tap targets at 44px or more.
- **Dark mode:** Tailwind `darkMode: 'media'`. Every color class needs a `dark:` counterpart.
- **Shared values:** site constants (URL, GitHub, contact, view/expiry options) are in `src/lib/site.ts`.

## Repo notes

- `database/` is the production bind mount. `*.db` files are git-ignored and must never be committed.
- `frontend/node_modules` and `.svelte-kit` are ignored. The `.dockerignore` files keep them out of images.
