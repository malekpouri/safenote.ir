package controllers

import (
	"safenote/internal/models"
	"safenote/internal/repositories"
	"safenote/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type NoteController struct {
	Repo *repositories.NoteRepository
}

func NewNoteController(repo *repositories.NoteRepository) *NoteController {
	return &NoteController{Repo: repo}
}

type CreateNoteRequest struct {
	EncryptedData       string `json:"encrypted_data"`
	IsPasswordProtected bool   `json:"is_password_protected"`
	ViewsRemaining      int    `json:"views_remaining"`
	Expiration          int    `json:"expiration"` // in minutes
}

func (c *NoteController) CreateNote(ctx *fiber.Ctx) error {
	var req CreateNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.EncryptedData == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Encrypted data is required"})
	}

	if req.ViewsRemaining <= 0 {
		req.ViewsRemaining = 1 // Default to 1 view
	}

	if req.Expiration <= 0 {
		req.Expiration = 1440 // Default to 24 hours
	}

	// Retry loop for collision handling (simple version)
	var id string
	for i := 0; i < 3; i++ {
		id = utils.GenerateShortID(6)
		// Check if exists (omitted for brevity, assuming collision is rare or DB constraint handles it)
		// Ideally Repo.Create should return error on duplicate key, then we retry.
		// For now, let's just try to create.
		expiresAt := time.Now().Add(time.Duration(req.Expiration) * time.Minute)

		note := &models.Note{
			ID:                  id,
			EncryptedData:       req.EncryptedData,
			IsPasswordProtected: req.IsPasswordProtected,
			ViewsRemaining:      req.ViewsRemaining,
			ExpiresAt:           expiresAt,
			CreatedAt:           time.Now(),
		}

		if err := c.Repo.Create(note); err == nil {
			return ctx.JSON(fiber.Map{
				"id": id,
			})
		}
		// If error, likely collision, retry
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate unique ID"})
}

func (c *NoteController) GetNote(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	note, err := c.Repo.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Note not found"})
	}

	if time.Now().After(note.ExpiresAt) {
		c.Repo.Delete(id) // Logical delete or actual delete
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Note expired"})
	}

	if note.ViewsRemaining <= 0 {
		c.Repo.Delete(id)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Note not found"})
	}

	// Decrement views
	if err := c.Repo.DecrementViews(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update view count"})
	}

	return ctx.JSON(fiber.Map{
		"encrypted_data":        note.EncryptedData,
		"is_password_protected": note.IsPasswordProtected,
		"created_at":            note.CreatedAt,
		"expires_at":            note.ExpiresAt,
	})
}
