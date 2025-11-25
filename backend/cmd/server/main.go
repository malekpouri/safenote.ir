package main

import (
	"log"
	"safenote/internal/api/controllers"
	"safenote/internal/repositories"
	"safenote/internal/scheduler"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Initialize Database
	dbPath := "/data/sqlite.db"
	repo, err := repositories.NewNoteRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Controllers
	noteController := controllers.NewNoteController(repo)

	// Initialize Scheduler
	sched := scheduler.NewScheduler(repo)
	sched.Start()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	api := app.Group("/api")
	api.Post("/notes", noteController.CreateNote)
	api.Get("/notes/:id", noteController.GetNote)
	api.Delete("/notes/:id", noteController.DeleteNote)

	log.Fatal(app.Listen(":8080"))
}
