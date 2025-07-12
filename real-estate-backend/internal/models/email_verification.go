package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmailVerification represents an email verification record
type EmailVerification struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Token     string    `json:"-" gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// EmailVerificationRequest represents a request to send email verification
type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmailRequest represents a request to verify email with token
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// BeforeCreate GORM hook to set ID
func (e *EmailVerification) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for EmailVerification model
func (EmailVerification) TableName() string {
	return "email_verifications"
}

// IsExpired checks if the verification token is expired
func (e *EmailVerification) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// EmailVerificationRepository handles database operations for email verifications
type EmailVerificationRepository struct {
	db *gorm.DB
}

// NewEmailVerificationRepository creates a new email verification repository
func NewEmailVerificationRepository(db *gorm.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

// Create creates a new email verification record
func (r *EmailVerificationRepository) Create(verification *EmailVerification) error {
	return r.db.Create(verification).Error
}

// GetByToken retrieves an email verification by token
func (r *EmailVerificationRepository) GetByToken(token string) (*EmailVerification, error) {
	var verification EmailVerification
	err := r.db.Preload("User").Where("token = ? AND is_used = false", token).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByUserID retrieves the latest email verification for a user
func (r *EmailVerificationRepository) GetByUserID(userID uuid.UUID) (*EmailVerification, error) {
	var verification EmailVerification
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// MarkAsUsed marks a verification as used
func (r *EmailVerificationRepository) MarkAsUsed(id uuid.UUID) error {
	return r.db.Model(&EmailVerification{}).Where("id = ?", id).Update("is_used", true).Error
}

// DeleteByUserID deletes all verification records for a user (cleanup)
func (r *EmailVerificationRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&EmailVerification{}).Error
}

// CleanupExpired deletes expired verification records
func (r *EmailVerificationRepository) CleanupExpired() error {
	return r.db.Where("expires_at < ? OR is_used = true", time.Now().Add(-24*time.Hour)).Delete(&EmailVerification{}).Error
}
