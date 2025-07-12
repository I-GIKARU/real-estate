package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PasswordReset represents a password reset token
type PasswordReset struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Token     string    `json:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// ForgotPasswordRequest represents the request to initiate password reset
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request to reset password
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// ChangePasswordRequest represents the request to change password (for logged-in users)
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// IsExpired checks if the reset token has expired
func (pr *PasswordReset) IsExpired() bool {
	return time.Now().After(pr.ExpiresAt)
}

// IsUsed checks if the reset token has been used
func (pr *PasswordReset) IsUsed() bool {
	return pr.UsedAt != nil
}

// MarkAsUsed marks the reset token as used
func (pr *PasswordReset) MarkAsUsed() {
	now := time.Now()
	pr.UsedAt = &now
}

// BeforeCreate GORM hook to set ID
func (pr *PasswordReset) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for PasswordReset model
func (PasswordReset) TableName() string {
	return "password_resets"
}

// PasswordResetRepository handles database operations for password resets
type PasswordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create creates a new password reset token
func (r *PasswordResetRepository) Create(reset *PasswordReset) error {
	return r.db.Create(reset).Error
}

// GetByToken retrieves a password reset by token
func (r *PasswordResetRepository) GetByToken(token string) (*PasswordReset, error) {
	var reset PasswordReset
	err := r.db.Where("token = ?", token).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

// GetByTokenWithUser retrieves a password reset by token with user data
func (r *PasswordResetRepository) GetByTokenWithUser(token string) (*PasswordReset, error) {
	var reset PasswordReset
	err := r.db.Preload("User").Where("token = ?", token).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

// Update updates a password reset token
func (r *PasswordResetRepository) Update(reset *PasswordReset) error {
	return r.db.Save(reset).Error
}

// DeleteExpiredTokens deletes all expired tokens
func (r *PasswordResetRepository) DeleteExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&PasswordReset{}).Error
}

// DeleteByUserID deletes all password reset tokens for a user
func (r *PasswordResetRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&PasswordReset{}).Error
}

// DeleteByToken deletes a password reset token
func (r *PasswordResetRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&PasswordReset{}).Error
}
