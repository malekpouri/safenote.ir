package controllers

import (
	"safenote/internal/models"
	"safenote/internal/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NoteController struct {
	Repo *repositories.NoteRepository
}

func NewNoteController(repo *repositories.NoteRepository) *NoteController {
	return &NoteController{Repo: repo}
}

type CreateNoteRequest struct {
	EncryptedData  string `json:"encrypted_data"`
	ViewsRemaining int    `json:"views_remaining"`
	Expiration     int    `json:"expiration"` // in minutes
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

	id := uuid.New().String() // Simple UUID for now, can be Base62 later
	expiresAt := time.Now().Add(time.Duration(req.Expiration) * time.Minute)

	note := &models.Note{
		ID:             id,
		EncryptedData:  req.EncryptedData,
		ViewsRemaining: req.ViewsRemaining,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now(),
	}

	if err := c.Repo.Create(note); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save note"})
	}

	return ctx.JSON(fiber.Map{
		"id": id,
	})
}
