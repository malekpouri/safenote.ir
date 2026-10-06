package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"safenote/internal/api/controllers"
	"safenote/internal/repositories"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	testToken = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQ"
	testSalt  = "c2FsdHNhbHRzYWx0c2FsdA"
)

type harness struct {
	t     *testing.T
	app   *fiber.App
	notes *controllers.NoteController
	repo  *repositories.NoteRepository
	now   time.Time
}

func newHarness(t *testing.T, opts Options) *harness {
	t.Helper()
	repo, err := repositories.NewNoteRepository(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { repo.Close() })
	app, notes := NewApp(repo, opts)
	h := &harness{t: t, app: app, notes: notes, repo: repo, now: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)}
	notes.Now = func() time.Time { return h.now }
	return h
}

func (h *harness) do(method, path string, body any, headers ...string) (int, map[string]any) {
	h.t.Helper()
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, r)
	req.Header.Set("Content-Type", "application/json")
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	resp, err := h.app.Test(req, -1)
	if err != nil {
		h.t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func (h *harness) create(extra map[string]any) string {
	h.t.Helper()
	body := map[string]any{"encrypted_data": "v2:ciphertext", "access_token": testToken, "salt": testSalt}
	for k, v := range extra {
		body[k] = v
	}
	status, out := h.do("POST", "/api/notes", body)
	if status != http.StatusCreated {
		h.t.Fatalf("create: status %d, body %v", status, out)
	}
	id, _ := out["id"].(string)
	if len(id) != controllers.IDLength {
		h.t.Fatalf("unexpected id %q", id)
	}
	return id
}

func TestCreateValidation(t *testing.T) {
	h := newHarness(t, Options{})
	cases := []struct {
		name string
		body map[string]any
		want int
	}{
		{"missing data", map[string]any{"access_token": testToken}, 400},
		{"legacy format", map[string]any{"encrypted_data": "abc", "access_token": testToken}, 400},
		{"missing token", map[string]any{"encrypted_data": "v2:x", "salt": testSalt}, 400},
		{"missing salt", map[string]any{"encrypted_data": "v2:x", "access_token": testToken}, 400},
		{"too many views", map[string]any{"encrypted_data": "v2:x", "access_token": testToken, "salt": testSalt, "views_remaining": 101}, 400},
		{"negative views", map[string]any{"encrypted_data": "v2:x", "access_token": testToken, "salt": testSalt, "views_remaining": -1}, 400},
		{"expiration too long", map[string]any{"encrypted_data": "v2:x", "access_token": testToken, "salt": testSalt, "expiration": 30*1440 + 1}, 400},
		{"too large", map[string]any{"encrypted_data": "v2:" + strings.Repeat("a", controllers.MaxEncryptedLength), "access_token": testToken, "salt": testSalt}, 413},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if status, out := h.do("POST", "/api/notes", tc.body); status != tc.want {
				t.Fatalf("got %d (%v), want %d", status, out, tc.want)
			}
		})
	}
}

func TestOpenRequiresTokenAndBurns(t *testing.T) {
	h := newHarness(t, Options{})
	id := h.create(nil)

	status, meta := h.do("GET", "/api/notes/"+id, nil)
	if status != 200 || meta["legacy"] != false || meta["salt"] != testSalt {
		t.Fatalf("meta: %d %v", status, meta)
	}

	// Wrong token must not consume the only view.
	if status, _ := h.do("POST", "/api/notes/"+id+"/open", map[string]any{"access_token": strings.Repeat("x", 43)}); status != 401 {
		t.Fatalf("wrong token: got %d", status)
	}
	if status, _ := h.do("POST", "/api/notes/"+id+"/open", nil); status != 401 {
		t.Fatalf("no token: got %d", status)
	}

	status, out := h.do("POST", "/api/notes/"+id+"/open", map[string]any{"access_token": testToken})
	if status != 200 || out["encrypted_data"] != "v2:ciphertext" || out["views_remaining"].(float64) != 0 {
		t.Fatalf("open: %d %v", status, out)
	}

	// Burned: row is gone immediately, not just hidden.
	var count int
	h.repo.DB.QueryRow(`SELECT COUNT(*) FROM notes`).Scan(&count)
	if count != 0 {
		t.Fatalf("expected note row deleted, found %d", count)
	}
	if status, _ := h.do("POST", "/api/notes/"+id+"/open", map[string]any{"access_token": testToken}); status != 404 {
		t.Fatalf("second open: got %d", status)
	}
}

func TestMultipleViewsAndExpiry(t *testing.T) {
	h := newHarness(t, Options{})
	id := h.create(map[string]any{"views_remaining": 3, "expiration": 60})

	for want := 2; want >= 1; want-- {
		_, out := h.do("POST", "/api/notes/"+id+"/open", map[string]any{"access_token": testToken})
		if int(out["views_remaining"].(float64)) != want {
			t.Fatalf("views_remaining = %v, want %d", out["views_remaining"], want)
		}
	}

	h.now = h.now.Add(61 * time.Minute)
	if status, _ := h.do("GET", "/api/notes/"+id, nil); status != 404 {
		t.Fatalf("expired meta: got %d", status)
	}
}

func TestConcurrentOpenOnlyOnce(t *testing.T) {
	h := newHarness(t, Options{})
	id := h.create(nil)

	var wg sync.WaitGroup
	var mu sync.Mutex
	successes := 0
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			status, _ := h.do("POST", "/api/notes/"+id+"/open", map[string]any{"access_token": testToken})
			if status == 200 {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if successes != 1 {
		t.Fatalf("1-view note opened %d times", successes)
	}
}

func TestDeleteRequiresToken(t *testing.T) {
	h := newHarness(t, Options{})
	id := h.create(nil)

	if status, _ := h.do("DELETE", "/api/notes/"+id, nil); status != 401 {
		t.Fatalf("delete without token: got %d", status)
	}
	if status, _ := h.do("DELETE", "/api/notes/"+id, map[string]any{"access_token": testToken}); status != 200 {
		t.Fatalf("delete with token: got %d", status)
	}
	if status, _ := h.do("GET", "/api/notes/"+id, nil); status != 404 {
		t.Fatalf("after delete: got %d", status)
	}
}

func TestLegacyNotes(t *testing.T) {
	h := newHarness(t, Options{})
	_, err := h.repo.DB.Exec(`INSERT INTO notes (id, encrypted_data, password_hash, is_password_protected, views_remaining, expires_at, created_at)
		VALUES ('Leg001', 'legacy-b64', 'pwhash', 1, 1, ?, ?)`, h.now.Add(time.Hour), h.now)
	if err != nil {
		t.Fatal(err)
	}

	_, meta := h.do("GET", "/api/notes/Leg001", nil)
	if meta["legacy"] != true {
		t.Fatalf("expected legacy note, got %v", meta)
	}
	if status, _ := h.do("DELETE", "/api/notes/Leg001", map[string]any{"password_hash": "wrong"}); status != 401 {
		t.Fatalf("legacy delete wrong hash: got %d", status)
	}
	status, out := h.do("POST", "/api/notes/Leg001/open", nil)
	if status != 200 || out["encrypted_data"] != "legacy-b64" {
		t.Fatalf("legacy open: %d %v", status, out)
	}
}

func TestStatsAndCleanup(t *testing.T) {
	h := newHarness(t, Options{})
	h.create(nil)
	h.create(map[string]any{"expiration": 10})

	_, stats := h.do("GET", "/api/admin/stats", nil)
	if stats["active_notes"].(float64) != 2 || stats["total_notes_created"].(float64) != 2 {
		t.Fatalf("stats: %v", stats)
	}

	n, err := h.repo.DeleteExpired(h.now.Add(11 * time.Minute))
	if err != nil || n != 1 {
		t.Fatalf("DeleteExpired = %d, %v; want 1", n, err)
	}
}

func TestRateLimit(t *testing.T) {
	h := newHarness(t, Options{ReadRateLimit: 2})
	codes := []int{}
	for i := 0; i < 3; i++ {
		status, _ := h.do("GET", "/api/notes/doesnotexist", nil, "X-Real-IP", "1.2.3.4")
		codes = append(codes, status)
	}
	if codes[2] != http.StatusTooManyRequests {
		t.Fatalf("expected 429 on third request, got %v", codes)
	}
	if status, _ := h.do("GET", "/api/notes/doesnotexist", nil, "X-Real-IP", "5.6.7.8"); status != 404 {
		t.Fatalf("other IP should not be limited, got %d", status)
	}
}

func TestInvalidIDRejected(t *testing.T) {
	h := newHarness(t, Options{})
	if status, _ := h.do("GET", "/api/notes/bad-id!", nil); status != 404 {
		t.Fatalf("got %d", status)
	}
}
