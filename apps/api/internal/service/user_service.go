package service

import (
	"lotto-journal/api/internal/models"
	"lotto-journal/api/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// FindOrCreate finds a user by LINE user ID or creates a new active user.
// Returns (user, isNewlyCreated, error).
func (s *UserService) FindOrCreate(lineUserID string) (*models.User, bool, error) {
	return s.repo.FindOrCreate(lineUserID)
}

// Deactivate marks the user as inactive (called on LINE unfollow event).
func (s *UserService) Deactivate(lineUserID string) error {
	return s.repo.UpdateStatus(lineUserID, "inactive")
}
