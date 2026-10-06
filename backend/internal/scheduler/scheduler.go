package scheduler

import (
	"context"
	"log"
	"safenote/internal/repositories"
	"time"
)

type Scheduler struct {
	Repo     *repositories.NoteRepository
	Interval time.Duration
}

func NewScheduler(repo *repositories.NoteRepository) *Scheduler {
	return &Scheduler{Repo: repo, Interval: 10 * time.Minute}
}

// Start purges expired notes immediately and then on every tick until ctx is
// cancelled.
func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.Interval)
		defer ticker.Stop()
		for {
			s.purge()
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}

func (s *Scheduler) purge() {
	count, err := s.Repo.DeleteExpired(time.Now())
	if err != nil {
		log.Printf("Error deleting expired notes: %v", err)
	} else if count > 0 {
		log.Printf("Deleted %d expired notes", count)
	}
}
