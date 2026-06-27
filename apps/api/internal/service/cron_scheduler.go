package service

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
)

type CronScheduler struct {
	db             *gorm.DB
	drawService    *DrawService
	resultService  *ResultService
	cron           *cron.Cron
	syncSchedule   string
	verifySchedule string
}

func NewCronScheduler(
	db *gorm.DB,
	drawService *DrawService,
	resultService *ResultService,
	syncSchedule string,
	verifySchedule string,
) *CronScheduler {
	return &CronScheduler{
		db:             db,
		drawService:    drawService,
		resultService:  resultService,
		syncSchedule:   syncSchedule,
		verifySchedule: verifySchedule,
	}
}

// Start launches the background scheduler loop using github.com/robfig/cron/v3.
func (s *CronScheduler) Start(ctx context.Context) {
	log.Println("[scheduler] Starting in-process cron scheduler background task...")

	// 1. Run GLO schedule sync immediately on boot
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[scheduler] panic in startup sync: %v\n%s", r, debug.Stack())
			}
		}()
		log.Println("[scheduler] Executing startup draw schedule sync...")
		if err := s.drawService.SyncDrawSchedule(); err != nil {
			log.Printf("[scheduler] Startup draw schedule sync failed: %v", err)
		}
	}()

	// 2. Initialize robfig/cron with Bangkok location and default panic recovery chain
	s.cron = cron.New(
		cron.WithLocation(bangkokLoc),
		cron.WithChain(cron.Recover(cron.DefaultLogger)),
	)

	// Job 1: Daily Draw Schedule Sync
	_, err := s.cron.AddFunc(s.syncSchedule, func() {
		log.Println("[scheduler] Triggering daily draw schedule sync...")
		if err := s.drawService.SyncDrawSchedule(); err != nil {
			log.Printf("[scheduler] Daily schedule sync failed: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("[scheduler] Failed to schedule daily sync job: %v", err)
	}

	// Job 2: Check Draw Results
	_, err = s.cron.AddFunc(s.verifySchedule, func() {
		s.checkResults()
	})
	if err != nil {
		log.Fatalf("[scheduler] Failed to schedule result checking job: %v", err)
	}

	// Start cron scheduler
	s.cron.Start()

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("[scheduler] Stopping background scheduler...")
	s.cron.Stop()
	log.Println("[scheduler] Background scheduler stopped.")
}

func (s *CronScheduler) checkResults() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[scheduler] panic in checkResults: %v\n%s", r, debug.Stack())
		}
	}()

	bkkNow := time.Now().In(bangkokLoc)
	todayUTC := time.Date(bkkNow.Year(), bkkNow.Month(), bkkNow.Day(), 0, 0, 0, 0, time.UTC)
	todayStr := todayUTC.Format("2006-01-02")

	var draws []models.Draw
	err := s.db.Where(repository.DrawColDrawDate+" = ?", todayStr).Limit(1).Find(&draws).Error
	if err != nil {
		log.Printf("[scheduler] Error querying today's draw: %v", err)
		return
	}

	if len(draws) == 0 {
		return // Not a draw day
	}

	draw := draws[0]
	if draw.IsVerified {
		return // Results already verified
	}

	log.Printf("[scheduler] Draw day detected (%s) and results unverified. Checking GLO results...", todayStr)
	if err := s.resultService.VerifyDrawResults(todayUTC); err != nil {
		log.Printf("[scheduler] Draw results verification failed: %v", err)
	}
}
