package service

import (
	"time"

	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
)

// bangkokLoc is the fixed UTC+7 timezone used for all draw-date calculations.
// Using time.FixedZone avoids loading the OS timezone database, which may be
// absent in minimal container images (e.g. scratch, distroless).
var bangkokLoc = time.FixedZone("ICT", 7*60*60)

type DrawService struct {
	repo *repository.DrawRepository
}

func NewDrawService(repo *repository.DrawRepository) *DrawService {
	return &DrawService{repo: repo}
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
func (s *DrawService) FindOrCreateUpcoming() (*models.Draw, error) {
	upcoming := NextDrawDate(time.Now())
	return s.repo.FindOrCreate(upcoming)
}
