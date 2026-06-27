package service

import (
	"fmt"
	"log"
	"time"

	"lotto-journal/api/internal/client"
	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
)

// bangkokLoc is the fixed UTC+7 timezone used for all draw-date calculations.
// Using time.FixedZone avoids loading the OS timezone database, which may be
// absent in minimal container images (e.g. scratch, distroless).
var bangkokLoc = time.FixedZone("ICT", 7*60*60)

type DrawService struct {
	repo   *repository.DrawRepository
	client *client.LotteryClient
}

func NewDrawService(repo *repository.DrawRepository, client *client.LotteryClient) *DrawService {
	return &DrawService{repo: repo, client: client}
}


// NextDrawDate returns the nearest upcoming draw date (1st or 16th of the month)
// expressed in Bangkok time. Candidates are evaluated in order; the first one
// that is >= today (Bangkok) is returned.
func NextDrawDate(now time.Time) time.Time {
	bkk := now.In(bangkokLoc)
	today := time.Date(bkk.Year(), bkk.Month(), bkk.Day(), 0, 0, 0, 0, bangkokLoc)

	candidates := []time.Time{
		time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, bangkokLoc),
		time.Date(today.Year(), today.Month(), 16, 0, 0, 0, 0, bangkokLoc),
		// time.Date with Month+1 overflows safely (e.g. December+1 = January next year)
		time.Date(today.Year(), today.Month()+1, 1, 0, 0, 0, 0, bangkokLoc),
	}

	for _, d := range candidates {
		if !d.Before(today) {
			return d
		}
	}

	// Unreachable given the three candidates, but keeps the compiler happy.
	return candidates[2]
}

// FindOrCreateUpcoming finds or creates the draws record for the next draw date.
//
// It uses a database-first approach: querying the local draws table first.
// If not found (e.g. on first startup), it falls back to mathematically calculating
// the next draw date as an emergency fallback, creating it in the DB.
func (s *DrawService) FindOrCreateUpcoming() (*models.Draw, error) {
	bkkNow := time.Now().In(bangkokLoc)
	todayUTC := time.Date(bkkNow.Year(), bkkNow.Month(), bkkNow.Day(), 0, 0, 0, 0, time.UTC)

	draw, err := s.repo.FindNextDraw(todayUTC)
	if err == nil {
		return draw, nil
	}

	// Emergency fallback
	fallbackDate := NextDrawDate(time.Now())
	return s.repo.FindOrCreate(fallbackDate)
}

// SyncDrawSchedule fetches the schedule from GLO API for the current year,
// deduplicates draw dates, and bulk updates our draws table.
func (s *DrawService) SyncDrawSchedule() error {
	year := time.Now().In(bangkokLoc).Year()
	dates, err := s.client.FetchDrawSchedule(year)
	if err != nil {
		return fmt.Errorf("fetch schedule from GLO: %w", err)
	}

	log.Printf("[draw_service] syncing schedule for year %d, got %d unique dates", year, len(dates))
	for _, d := range dates {
		_, err := s.repo.FindOrCreate(d)
		if err != nil {
			log.Printf("[draw_service] failed to save synced date %s: %v", d.Format("2006-01-02"), err)
		}
	}
	return nil
}
