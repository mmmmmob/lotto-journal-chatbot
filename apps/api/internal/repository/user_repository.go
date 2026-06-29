package repository

import (
	"lotto-journal/api/internal/models"

	"gorm.io/gorm"
)

const (
	TableUsers        = "users"
	UserColID         = "id"
	UserColLineUserID = "line_user_id"
	UserColStatus     = "status"
	UserColLanguage   = "language"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindByLineUserID returns the user with the given LINE user ID, or nil + gorm.ErrRecordNotFound.
func (r *UserRepository) FindByLineUserID(lineUserID string) (*models.User, error) {
	var user models.User
	result := r.db.Where(UserColLineUserID+" = ?", lineUserID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindOrCreate finds the user by LINE user ID or creates a new active user.
// Returns (user, isNewlyCreated, error).
func (r *UserRepository) FindOrCreate(lineUserID string) (*models.User, bool, error) {
	user := models.User{
		LineUserID: lineUserID,
		Status:     "active",
	}
	result := r.db.Where(UserColLineUserID+" = ?", lineUserID).FirstOrCreate(&user)
	if result.Error != nil {
		return nil, false, result.Error
	}
	// RowsAffected == 1 means the record was freshly created; 0 means it already existed.
	return &user, result.RowsAffected > 0, nil
}

// UpdateStatus sets the account_status for the user with the given LINE user ID.
// Called when a user unfollows the LINE Official Account.
func (r *UserRepository) UpdateStatus(lineUserID string, status string) error {
	return r.db.
		Model(&models.User{}).
		Where(UserColLineUserID+" = ?", lineUserID).
		Update(UserColStatus, status).
		Error
}

// UpdateLanguage sets the language for the user with the given LINE user ID.
func (r *UserRepository) UpdateLanguage(lineUserID string, language string) error {
	return r.db.
		Model(&models.User{}).
		Where(UserColLineUserID+" = ?", lineUserID).
		Update(UserColLanguage, language).
		Error
}

