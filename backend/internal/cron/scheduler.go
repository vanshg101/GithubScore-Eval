package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

// Scheduler wraps robfig/cron to manage scheduled jobs.
type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler creates a new cron scheduler.
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithLogger(cron.DefaultLogger)),
	}
}

// AddJob registers a job with the given cron expression.
// The expression follows standard crontab format: "minute hour day month weekday".
func (s *Scheduler) AddJob(schedule string, job func()) error {
	_, err := s.cron.AddFunc(schedule, job)
	if err != nil {
		return err
	}
	log.Printf("[CRON] Job registered with schedule: %s", schedule)
	return nil
}

// Start begins the scheduler in a background goroutine.
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("[CRON] Scheduler started")
}

// Stop gracefully stops the scheduler, waiting for running jobs to finish.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("[CRON] Scheduler stopped")
}
