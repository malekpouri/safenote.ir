package scheduler

import (
	"log"
	"safenote/internal/repositories"
	"time"
)

type Scheduler struct {
	Repo *repositories.NoteRepository
}

func NewScheduler(repo *repositories.NoteRepository) *Scheduler {
	return &Scheduler{Repo: repo}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			count, err := s.Repo.DeleteExpired()
			if err != nil {
				log.Printf("Error deleting expired notes: %v", err)
			} else if count > 0 {
				log.Printf("Deleted %d expired notes", count)
			}
		}
	}()
}
