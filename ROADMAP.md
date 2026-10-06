# SafeNote Roadmap

This roadmap tracks the fixes from the October 2026 code review and the UI overhaul. Items are grouped by phase and ordered by priority within each phase.

Legend: ✅ done · 🔜 planned / follow-up

---

## Phase 1: Encryption redesign (critical)

The original scheme used a 6-character link key (about 36 bits) and derived the AES key as `SHA-256(key + password)`. It also stored an unsalted `SHA-256(password)` on the server. Anyone who held the ciphertext could brute-force it offline.

| # | Item | Status |
|---|------|--------|
| 1.1 | Replace the Go→WASM crypto module (2.6 MB) with the browser-native **WebCrypto API** | ✅ |
| 1.2 | Short link key: 10 random base64url characters (60 bits) in the URL fragment, never sent to the server. It is stretched with PBKDF2 (600k rounds) and a per-note salt, so brute-forcing even a stolen database is impractical | ✅ |
| 1.3 | Key and password stretched together with **PBKDF2-SHA256, 600,000 iterations** and a random per-note salt | ✅ |
| 1.4 | Separate keys via **HKDF-SHA256**: one encryption key, one *access token* | ✅ |
| 1.5 | The server stores only `SHA-256(access token)`. The token is needed to open **and** delete a note, so guessing an ID no longer burns views or deletes notes | ✅ |
| 1.6 | Wrong password → `401` **before** a view is consumed (passwords are checked server-side via the token, without revealing anything) | ✅ |
| 1.7 | Versioned ciphertext (`v2:`). Legacy v1 notes can still be opened until they expire (≤ 30 days) | ✅ |
| 1.8 | Password generator uses `crypto.getRandomValues` and 16 characters | ✅ |
| 1.9 | Note IDs: 6 characters from `crypto/rand` instead of `math/rand`. Length doesn't matter for security because the access token gates everything | ✅ |
| 1.10 | Remove the legacy v1 code path 30 days after deploying v2 | 🔜 |

## Phase 2: Data and repository hygiene

| # | Item | Status |
|---|------|--------|
| 2.1 | Stop tracking `database/sqlite.db`; ignore `database/*.db*` | ✅ |
| 2.2 | Stop tracking `frontend/node_modules` and `frontend/.svelte-kit` (≈ 4,000 files) | ✅ |
| 2.3 | Add `.dockerignore` files so local `node_modules` never leaks into images | ✅ |
| 2.4 | Purge the database from git history (`git filter-repo`). This rewrites history, so the owner must run it by hand | 🔜 |

## Phase 3: Backend correctness

| # | Item | Status |
|---|------|--------|
| 3.1 | **Atomic view consumption**: read, check, decrement and delete run in one transaction on one connection, so a 1-view note can't be read twice | ✅ |
| 3.2 | Delete the note **immediately** when its last view is used | ✅ |
| 3.3 | Expiry comparisons use bound UTC parameters instead of `CURRENT_TIMESTAMP` string comparison | ✅ |
| 3.4 | Server-side validation: views 1–100, expiration 5 min – 30 days, payload size cap that fits 10,000 Persian characters | ✅ |
| 3.5 | Real errors are no longer reported as "ID collision"; only `UNIQUE` constraint errors are retried | ✅ |
| 3.6 | `DB_PATH`, `PORT` and rate limits are configurable through the environment | ✅ |
| 3.7 | SQLite: WAL mode, `busy_timeout`, a single writer connection, an index on `expires_at` | ✅ |
| 3.8 | Per-IP rate limiting (create / read), with the client IP forwarded by the frontend | ✅ |
| 3.9 | Graceful shutdown, body-size limit, cleanup runs at startup and every 10 minutes | ✅ |
| 3.10 | Go tests for the repository and the HTTP API | ✅ |
| 3.11 | `go mod tidy` (direct dependencies were all marked `// indirect`) | ✅ |

## Phase 4: Frontend platform

| # | Item | Status |
|---|------|--------|
| 4.1 | One shared proxy helper. Backend status codes and JSON bodies pass through unchanged (no more 400→500 or HTML error pages) | ✅ |
| 4.2 | Security headers: CSP, `Referrer-Policy: no-referrer`, `X-Frame-Options`, `Permissions-Policy`, `no-store` on notes and the API | ✅ |
| 4.3 | Locale resolved **on the server** (cookie → `Accept-Language`), so `<html lang dir>` is correct on first paint with no RTL flash | ✅ |
| 4.4 | Fonts self-hosted (Inter + Vazirmatn) instead of Google Fonts, for privacy and reachability from Iran | ✅ |
| 4.5 | Unit tests for the crypto module (round trips, legacy v1 compatibility, Persian payload size) | ✅ |
| 4.6 | Slimmer Docker image (no WASM stage, `npm ci`, no `node_modules` in the runtime image), plus healthchecks | ✅ |

## Phase 5: UI / UX overhaul

| # | Item | Status |
|---|------|--------|
| 5.1 | Shared layout: sticky header (logo, navigation, language switch) and a full **footer** (links, source code, site URL, credits) | ✅ |
| 5.2 | Dark mode that follows the system setting | ✅ |
| 5.3 | Home: a clearer hero section, trust badges, options shown as a compact always-visible row, show/hide password, Ctrl/⌘+Enter to submit | ✅ |
| 5.4 | Success screen: copy and native share, a summary (views / expiry / password), a "shown only once" warning | ✅ |
| 5.5 | "How it works" steps and a feature grid replace the wall of text; FAQ section | ✅ |
| 5.6 | Viewer: clear states (loading, missing key, not found, ready, password, revealed). Inline password errors that don't burn the note. Shows views remaining | ✅ |
| 5.7 | Two-step inline delete confirmation instead of blocking `confirm()` dialogs | ✅ |
| 5.8 | About and Privacy pages translated (EN / FA), with the privacy text updated to match the new design | ✅ |
| 5.9 | Stats page restyled; custom error / 404 page | ✅ |
| 5.10 | Accessibility: labels, focus rings, `aria-live` toasts, toasts positioned correctly for RTL | ✅ |
| 5.11 | SEO: fixed `robots.txt` (it was missing `User-agent`), canonical `https://www.safenote.ir/`, sitemap, a real `og-image.png` | ✅ |

## Phase 5b: Mobile optimization

Audited with a headless phone viewport (320, 375 and 390 px; English and Persian) measuring overflow, input font sizes and tap-target sizes.

| # | Item | Status |
|---|------|--------|
| 5b.1 | Form fields are 16px on phones, so iOS Safari no longer zooms in when a field is focused | ✅ |
| 5b.2 | All controls are at least 44px tall on phones (buttons, header icons, password toggle, FAQ rows, footer links, toast close) | ✅ |
| 5b.3 | Compact hero on phones: the note box is above the fold even at 320×568 | ✅ |
| 5b.4 | "Create secure link" stays pinned to the bottom of the screen while the form is in view | ✅ |
| 5b.5 | The share link wraps, so the key after `#` is always visible; Copy and Share are full-width side by side | ✅ |
| 5b.6 | Mobile menu: dimmed backdrop, 48px rows with icons, scroll lock, closes on Escape or a tap outside | ✅ |
| 5b.7 | Toasts sit at the top on phones so they never cover action buttons; redundant success toasts removed | ✅ |
| 5b.8 | Safe-area insets (notch / home indicator), `viewport-fit=cover`, `100dvh`, no double-tap delay | ✅ |
| 5b.9 | Two-column footer on phones; tighter section spacing; settings side by side; 2-up stats cards | ✅ |
| 5b.10 | Mobile keyboards: no autocapitalize or autocorrect on passwords, `enterkeyhint`, no auto-selection popups | ✅ |

## Phase 5c: Minimal redesign

The first redesign was too busy. The UI now follows a minimal, single-column layout.

| # | Item | Status |
|---|------|--------|
| 5c.1 | Home is just a heading, the note box, a collapsed "Options" line (with a summary of the current settings) and one button | ✅ |
| 5c.2 | "How it works" and the FAQ moved to About; the feature grid, badges, trust chips and gradients are gone | ✅ |
| 5c.3 | Text-only header (name and language switch) and a one-line footer (About · Privacy · Stats · GitHub) | ✅ |
| 5c.4 | Monochrome palette with borders instead of cards and shadows; one consistent column width on every page | ✅ |
| 5c.5 | Success, viewer, stats and error pages follow the same pattern; secondary actions are plain text buttons | ✅ |
| 5c.6 | Mobile guarantees from 5b kept (16px fields, 44px targets, safe areas); unused components and icons removed | ✅ |

## Phase 5d: Short links and a fresher look

| # | Item | Status |
|---|------|--------|
| 5d.1 | SMS-friendly links: `https://safenote.ir/<6-char id>#<10-char key>` (37 characters, was 59). Old `/n/<id>` links keep working. The bare domain 301-redirects to www and the `#key` survives the redirect | ✅ |
| 5d.2 | Brand color, logo and a soft background glow brought back while keeping the single-column layout | ✅ |
| 5d.3 | Note box is a card with option pills (views · expiry · password) that open the settings | ✅ |
| 5d.4 | Gradient accent in the title, pop-in success check, icons on viewer states | ✅ |

## Phase 6: Documentation

| # | Item | Status |
|---|------|--------|
| 6.1 | README redesigned for the GitHub page, with the live URL, security model, architecture and self-hosting guide | ✅ |
| 6.2 | GitHub repository description and website set to https://www.safenote.ir/ | ✅ |

## Follow-ups (not in this pass)

- 🔜 QR code for the share link (needs a small library and a CSP review).
- 🔜 Run containers as a non-root user (the database volume's ownership must be migrated first).
- 🔜 End-to-end browser tests (Playwright) for the create → open → burn flow.
- 🔜 Optional "burn after reading" notice for the sender (read receipts), without storing identities.
