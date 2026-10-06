package controllers

import (
	"errors"
	"log"
	"regexp"
	"safenote/internal/models"
	"safenote/internal/repositories"
	"safenote/internal/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	// Short IDs keep links SMS-friendly. Guessing an ID gains nothing: opening
	// or deleting a note also requires the access token derived from the key.
	IDLength = 6

	// MaxEncryptedLength comfortably fits the 10,000-character UI limit even
	// for 3-byte UTF-8 text after AES-GCM and base64url overhead (~40 KB).
	MaxEncryptedLength = 60000
	MinViews           = 1
	MaxViews           = 100
	DefaultViews       = 1
	MinExpiration      = 5         // minutes
	MaxExpiration      = 30 * 1440 // 30 days
	DefaultExpiration  = 1440      // 24 hours
	EncryptionPrefix   = "v2:"
)

var (
	idPattern    = regexp.MustCompile(`^[A-Za-z0-9]{6,32}$`)
	tokenPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{22,128}$`)
	saltPattern  = regexp.MustCompile(`^[A-Za-z0-9_-]{22,64}$`)
)

type NoteController struct {
	Repo *repositories.NoteRepository
	// Now is overridable for tests.
	Now func() time.Time
}

func NewNoteController(repo *repositories.NoteRepository) *NoteController {
	return &NoteController{Repo: repo, Now: time.Now}
}

func jsonError(ctx *fiber.Ctx, status int, msg string) error {
	return ctx.Status(status).JSON(fiber.Map{"error": msg})
}

type CreateNoteRequest struct {
	EncryptedData       string `json:"encrypted_data"`
	AccessToken         string `json:"access_token"`
	Salt                string `json:"salt"`
	IsPasswordProtected bool   `json:"is_password_protected"`
	ViewsRemaining      int    `json:"views_remaining"`
	Expiration          int    `json:"expiration"` // in minutes
}

func (c *NoteController) CreateNote(ctx *fiber.Ctx) error {
	var req CreateNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return jsonError(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	switch {
	case req.EncryptedData == "":
		return jsonError(ctx, fiber.StatusBadRequest, "Encrypted data is required")
	case !strings.HasPrefix(req.EncryptedData, EncryptionPrefix):
		return jsonError(ctx, fiber.StatusBadRequest, "Unsupported encryption format")
	case len(req.EncryptedData) > MaxEncryptedLength:
		return jsonError(ctx, fiber.StatusRequestEntityTooLarge, "Note content too long")
	case !tokenPattern.MatchString(req.AccessToken):
		return jsonError(ctx, fiber.StatusBadRequest, "Invalid access token")
	case !saltPattern.MatchString(req.Salt):
		return jsonError(ctx, fiber.StatusBadRequest, "Invalid salt")
	}

	if req.ViewsRemaining == 0 {
		req.ViewsRemaining = DefaultViews
	}
	if req.Expiration == 0 {
		req.Expiration = DefaultExpiration
	}
	if req.ViewsRemaining < MinViews || req.ViewsRemaining > MaxViews {
		return jsonError(ctx, fiber.StatusBadRequest, "Views must be between 1 and 100")
	}
	if req.Expiration < MinExpiration || req.Expiration > MaxExpiration {
		return jsonError(ctx, fiber.StatusBadRequest, "Expiration must be between 5 minutes and 30 days")
	}

	now := c.Now().UTC()
	note := &models.Note{
		EncryptedData:       req.EncryptedData,
		AccessTokenHash:     utils.HashToken(req.AccessToken),
		KdfSalt:             req.Salt,
		IsPasswordProtected: req.IsPasswordProtected,
		ViewsRemaining:      req.ViewsRemaining,
		ExpiresAt:           now.Add(time.Duration(req.Expiration) * time.Minute),
		CreatedAt:           now,
	}

	for i := 0; i < 5; i++ {
		note.ID = utils.GenerateShortID(IDLength)
		err := c.Repo.Create(note)
		if err == nil {
			return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
				"id":         note.ID,
				"expires_at": note.ExpiresAt,
			})
		}
		if !errors.Is(err, repositories.ErrDuplicateID) {
			log.Printf("create note: %v", err)
			return jsonError(ctx, fiber.StatusInternalServerError, "Failed to store note")
		}
	}
	return jsonError(ctx, fiber.StatusInternalServerError, "Failed to generate unique ID")
}

// GetNoteMeta returns non-secret metadata without consuming a view.
func (c *NoteController) GetNoteMeta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !idPattern.MatchString(id) {
		return jsonError(ctx, fiber.StatusNotFound, "Note not found")
	}
	note, err := c.Repo.GetByID(id, c.Now())
	if err != nil {
		return c.repoError(ctx, err)
	}
	return ctx.JSON(fiber.Map{
		"is_password_protected": note.IsPasswordProtected,
		"legacy":                note.IsLegacy(),
		"salt":                  note.KdfSalt,
		"expires_at":            note.ExpiresAt,
	})
}

type tokenRequest struct {
	AccessToken string `json:"access_token"`
	// PasswordHash authorizes deletion of legacy (v1) password-protected notes.
	PasswordHash string `json:"password_hash"`
}

func parseTokenRequest(ctx *fiber.Ctx) tokenRequest {
	var req tokenRequest
	if len(ctx.Body()) > 0 {
		_ = ctx.BodyParser(&req)
	}
	return req
}

func authorizeAccess(token string) func(*models.Note) error {
	return func(n *models.Note) error {
		if n.IsLegacy() {
			return nil
		}
		if token == "" || !utils.ConstantTimeEqual(utils.HashToken(token), n.AccessTokenHash) {
			return repositories.ErrUnauthorized
		}
		return nil
	}
}

// OpenNote consumes one view and returns the ciphertext. For v2 notes the
// caller must present the access token derived from the link key (and
// password), so a wrong password is rejected before a view is spent.
func (c *NoteController) OpenNote(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !idPattern.MatchString(id) {
		return jsonError(ctx, fiber.StatusNotFound, "Note not found")
	}
	req := parseTokenRequest(ctx)
	note, err := c.Repo.Consume(id, c.Now(), authorizeAccess(req.AccessToken))
	if err != nil {
		return c.repoError(ctx, err)
	}
	return ctx.JSON(fiber.Map{
		"encrypted_data":        note.EncryptedData,
		"is_password_protected": note.IsPasswordProtected,
		"views_remaining":       note.ViewsRemaining,
		"created_at":            note.CreatedAt,
		"expires_at":            note.ExpiresAt,
	})
}

func (c *NoteController) DeleteNote(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !idPattern.MatchString(id) {
		return jsonError(ctx, fiber.StatusNotFound, "Note not found")
	}
	req := parseTokenRequest(ctx)
	err := c.Repo.DeleteAuthorized(id, c.Now(), func(n *models.Note) error {
		if n.IsLegacy() {
			if n.IsPasswordProtected && !utils.ConstantTimeEqual(req.PasswordHash, n.PasswordHash) {
				return repositories.ErrUnauthorized
			}
			return nil
		}
		return authorizeAccess(req.AccessToken)(n)
	})
	if err != nil {
		return c.repoError(ctx, err)
	}
	return ctx.JSON(fiber.Map{"message": "Note deleted successfully"})
}

func (c *NoteController) GetStats(ctx *fiber.Ctx) error {
	stats, err := c.Repo.GetStats(c.Now())
	if err != nil {
		log.Printf("stats: %v", err)
		return jsonError(ctx, fiber.StatusInternalServerError, "Failed to get stats")
	}
	return ctx.JSON(stats)
}

func (c *NoteController) repoError(ctx *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, repositories.ErrNotFound):
		return jsonError(ctx, fiber.StatusNotFound, "Note not found or expired")
	case errors.Is(err, repositories.ErrUnauthorized):
		return jsonError(ctx, fiber.StatusUnauthorized, "Invalid key or password")
	default:
		log.Printf("repository error: %v", err)
		return jsonError(ctx, fiber.StatusInternalServerError, "Internal server error")
	}
}
