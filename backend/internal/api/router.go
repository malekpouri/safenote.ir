package api

import (
	"safenote/internal/api/controllers"
	"safenote/internal/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Options struct {
	// Requests per minute per client IP. Zero disables the limiter.
	CreateRateLimit int
	ReadRateLimit   int
}

// clientIP prefers the X-Real-IP header set by the SvelteKit proxy. The
// backend is only reachable from the internal network, so the header can be
// trusted.
func clientIP(c *fiber.Ctx) string {
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return c.IP()
}

func rateLimit(perMinute int, name string) fiber.Handler {
	if perMinute <= 0 {
		return func(c *fiber.Ctx) error { return c.Next() }
	}
	return limiter.New(limiter.Config{
		Max:        perMinute,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return name + ":" + clientIP(c)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too many requests, please try again later"})
		},
	})
}

func NewApp(repo *repositories.NoteRepository, opts Options) (*fiber.App, *controllers.NoteController) {
	app := fiber.New(fiber.Config{
		AppName:               "safenote",
		BodyLimit:             128 * 1024,
		DisableStartupMessage: true,
	})
	app.Use(recover.New())

	notes := controllers.NewNoteController(repo)

	app.Get("/health", func(c *fiber.Ctx) error {
		if err := repo.DB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "db unavailable"})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	read := rateLimit(opts.ReadRateLimit, "read")
	api := app.Group("/api")
	api.Post("/notes", rateLimit(opts.CreateRateLimit, "create"), notes.CreateNote)
	api.Get("/notes/:id", read, notes.GetNoteMeta)
	api.Post("/notes/:id/open", read, notes.OpenNote)
	api.Delete("/notes/:id", read, notes.DeleteNote)
	api.Get("/admin/stats", read, notes.GetStats)

	return app, notes
}
