package service

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronScheduler struct {
	drawService    *DrawService
	resultService  *ResultService
	cron           *cron.Cron
	syncSchedule   string
	verifySchedule string
}

func NewCronScheduler(
	drawService *DrawService,
	resultService *ResultService,
	syncSchedule string,
	verifySchedule string,
) *CronScheduler {
	return &CronScheduler{
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
		if err := s.drawService.SyncDrawSchedule(ctx); err != nil {
			log.Printf("[scheduler] Startup draw schedule sync failed: %v", err)
		}
		log.Println("[scheduler] Executing startup results catch-up check...")
		s.checkResultsStartup(ctx)
	}()

	// 2. Initialize robfig/cron with Bangkok location and default panic recovery chain
	s.cron = cron.New(
		cron.WithLocation(bangkokLoc),
		cron.WithChain(cron.Recover(cron.DefaultLogger)),
	)

	// Job 1: Daily Draw Schedule Sync
	_, err := s.cron.AddFunc(s.syncSchedule, func() {
		log.Println("[scheduler] Triggering daily draw schedule sync...")
		if err := s.drawService.SyncDrawSchedule(context.Background()); err != nil {
			log.Printf("[scheduler] Daily schedule sync failed: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("[scheduler] Failed to schedule daily sync job: %v", err)
	}

	// Job 2: Check Draw Results
	_, err = s.cron.AddFunc(s.verifySchedule, func() {
		s.checkResults(context.Background())
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

func (s *CronScheduler) checkResultsStartup(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[scheduler] panic in checkResultsStartup: %v\n%s", r, debug.Stack())
		}
	}()

	bkkNow := time.Now().In(bangkokLoc)
	todayUTC := time.Date(bkkNow.Year(), bkkNow.Month(), bkkNow.Day(), 0, 0, 0, 0, time.UTC)

	// Check if there are any unverified draws on or before today
	_, err := s.drawService.repo.FindLatestUnverified(todayUTC)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("[scheduler] No unverified draws to catch up.")
			return // Nothing to catch up!
		}
		log.Printf("[scheduler] Error querying unverified draws for startup: %v", err)
		return
	}

	log.Printf("[scheduler] Unverified draws detected. Checking latest GLO results...")
	if err := s.resultService.VerifyLatestDrawResults(ctx); err != nil {
		log.Printf("[scheduler] Draw results verification failed: %v", err)
	}
}

func (s *CronScheduler) checkResults(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[scheduler] panic in checkResults: %v\n%s", r, debug.Stack())
		}
	}()

	bkkNow := time.Now().In(bangkokLoc)
	todayUTC := time.Date(bkkNow.Year(), bkkNow.Month(), bkkNow.Day(), 0, 0, 0, 0, time.UTC)

	// Only call GLO if today is a scheduled draw day and is unverified
	draw, err := s.drawService.repo.FindByDate(todayUTC)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return // Today is not a draw day, do nothing.
		}
		log.Printf("[scheduler] Error querying today's draw: %v", err)
		return
	}

	if draw.IsVerified {
		return // Today's draw is already verified, do nothing.
	}

	log.Printf("[scheduler] Today is a draw day and is unverified. Checking latest GLO results...")
	if err := s.resultService.VerifyLatestDrawResults(ctx); err != nil {
		log.Printf("[scheduler] Draw results verification failed: %v", err)
	}
}
