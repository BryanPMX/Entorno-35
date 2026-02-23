package services

import (
	"log"
	"time"

	"github.com/entorno35/backend/internal/core/ports"
)

// PendingRegistrationCleanupService periodically deletes expired, incomplete registration drafts.
type PendingRegistrationCleanupService struct {
	repo      ports.PendingRegistrationRepository
	interval  time.Duration
	retention time.Duration
}

// NewPendingRegistrationCleanupService creates a cleanup worker service.
func NewPendingRegistrationCleanupService(repo ports.PendingRegistrationRepository, interval, retention time.Duration) *PendingRegistrationCleanupService {
	return &PendingRegistrationCleanupService{
		repo:      repo,
		interval:  interval,
		retention: retention,
	}
}

// Start launches a background cleanup ticker and returns a stop function.
func (s *PendingRegistrationCleanupService) Start() func() {
	if s == nil || s.repo == nil || s.interval <= 0 || s.retention < 0 {
		return func() {}
	}

	stopCh := make(chan struct{})
	ticker := time.NewTicker(s.interval)

	run := func() {
		cutoff := time.Now().UTC().Add(-s.retention)
		deleted, err := s.repo.DeleteExpiredIncompleteBefore(cutoff)
		if err != nil {
			log.Printf("Warning: pending registration cleanup failed: %v", err)
			return
		}
		if deleted > 0 {
			log.Printf("Pending registration cleanup removed %d expired draft(s)", deleted)
		}
	}

	go func() {
		// Run once shortly after startup to keep the table clean without waiting a full interval.
		run()
		for {
			select {
			case <-ticker.C:
				run()
			case <-stopCh:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(stopCh)
	}
}
