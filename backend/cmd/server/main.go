package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"safenote/internal/api"
	"safenote/internal/repositories"
	"safenote/internal/scheduler"
	"strconv"
	"syscall"
	"time"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("%s must be an integer, got %q", key, v)
	}
	return n
}

func main() {
	dbPath := env("DB_PATH", "/data/sqlite.db")
	repo, err := repositories.NewNoteRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database at %s: %v", dbPath, err)
	}

	app, _ := api.NewApp(repo, api.Options{
		CreateRateLimit: envInt("RATE_LIMIT_CREATE", 20),
		ReadRateLimit:   envInt("RATE_LIMIT_READ", 60),
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	scheduler.NewScheduler(repo).Start(ctx)

	addr := ":" + env("PORT", "8080")
	go func() {
		log.Printf("SafeNote backend listening on %s (db: %s)", addr, dbPath)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("shutdown: %v", err)
	}
	if err := repo.Close(); err != nil {
		log.Printf("closing database: %v", err)
	}
}
