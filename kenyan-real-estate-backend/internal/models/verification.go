package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VerificationType represents the type of verification
type VerificationType string

const (
	VerificationTypeEmail    VerificationType = "email"
	VerificationTypePhone    VerificationType = "phone"
	VerificationTypeDocument VerificationType = "document"
	VerificationTypeManual   VerificationType = "manual"
)

// VerificationStatus represents the status of verification
type VerificationStatus string

const (
	VerificationStatusPending   VerificationStatus = "pending"
	VerificationStatusVerified  VerificationStatus = "verified"
	VerificationStatusExpired   VerificationStatus = "expired"
	VerificationStatusRejected  VerificationStatus = "rejected"
)

// UserVerification represents a verification record for a user
type UserVerification struct {
	ID           uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID          `json:"user_id" gorm:"type:uuid;not null"`
	Type         VerificationType   `json:"type" gorm:"not null;type:varchar(20)"`
	Status       VerificationStatus `json:"status" gorm:"not null;type:varchar(20);default:'pending'"`
	Code         *string            `json:"-" gorm:"type:varchar(10)"` // OTP/verification code
	Token        *string            `json:"-" gorm:"type:varchar(255)"` // Email verification token
	ExpiresAt    *time.Time         `json:"expires_at"`
	VerifiedAt   *time.Time         `json:"verified_at"`
	VerifiedBy   *uuid.UUID         `json:"verified_by" gorm:"type:uuid"` // Admin who verified
	DocumentURL  *string            `json:"document_url"`
	DocumentType *string            `json:"document_type"`
	Notes        *string            `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt     `json:"-" gorm:"index"`

	// Relationships
	User       User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	VerifierUser *User `json:"verifier_user,omitempty" gorm:"foreignKey:VerifiedBy"`
}

// EmailVerificationRequest represents a request to send email verification

// PhoneVerificationRequest represents a request to send phone verification
type PhoneVerificationRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// VerifyCodeRequest represents a request to verify a code
type VerifyCodeRequest struct {
	Code string `json:"code" binding:"required,min=4,max=10"`
	Type VerificationType `json:"type" binding:"required"`
}

// DocumentVerificationRequest represents a request for document verification
type DocumentVerificationRequest struct {
	DocumentType string `json:"document_type" binding:"required"`
	DocumentURL  string `json:"document_url" binding:"required"`
	Notes        string `json:"notes,omitempty"`
}

// AdminVerificationRequest represents an admin action on verification
type AdminVerificationRequest struct {
	Status VerificationStatus `json:"status" binding:"required"`
	Notes  string            `json:"notes,omitempty"`
}

// VerificationResponse represents a verification response
type VerificationResponse struct {
	ID           uuid.UUID          `json:"id"`
	UserID       uuid.UUID          `json:"user_id"`
	Type         VerificationType   `json:"type"`
	Status       VerificationStatus `json:"status"`
	ExpiresAt    *time.Time         `json:"expires_at"`
	VerifiedAt   *time.Time         `json:"verified_at"`
	DocumentType *string            `json:"document_type,omitempty"`
	Notes        *string            `json:"notes,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// ToResponse converts UserVerification to VerificationResponse
func (v *UserVerification) ToResponse() *VerificationResponse {
	return &VerificationResponse{
		ID:           v.ID,
		UserID:       v.UserID,
		Type:         v.Type,
		Status:       v.Status,
		ExpiresAt:    v.ExpiresAt,
		VerifiedAt:   v.VerifiedAt,
		DocumentType: v.DocumentType,
		Notes:        v.Notes,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}

// BeforeCreate GORM hook to set ID
func (v *UserVerification) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for UserVerification model
func (UserVerification) TableName() string {
	return "user_verifications"
}

// IsExpired checks if verification is expired
func (v *UserVerification) IsExpired() bool {
	if v.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*v.ExpiresAt)
}

// UserVerificationRepository handles database operations for user verifications
type UserVerificationRepository struct {
	db *gorm.DB
}

// NewUserVerificationRepository creates a new user verification repository
func NewUserVerificationRepository(db *gorm.DB) *UserVerificationRepository {
	return &UserVerificationRepository{db: db}
}

// Create creates a new verification record
func (r *UserVerificationRepository) Create(verification *UserVerification) error {
	return r.db.Create(verification).Error
}

// GetByID retrieves a verification by ID
func (r *UserVerificationRepository) GetByID(id uuid.UUID) (*UserVerification, error) {
	var verification UserVerification
	err := r.db.Preload("User").Preload("VerifierUser").Where("id = ?", id).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByUserAndType retrieves the latest verification by user and type
func (r *UserVerificationRepository) GetByUserAndType(userID uuid.UUID, verificationType VerificationType) (*UserVerification, error) {
	var verification UserVerification
	err := r.db.Where("user_id = ? AND type = ?", userID, verificationType).
		Order("created_at DESC").
		First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByToken retrieves a verification by token
func (r *UserVerificationRepository) GetByToken(token string) (*UserVerification, error) {
	var verification UserVerification
	err := r.db.Preload("User").Where("token = ? AND status = ?", token, VerificationStatusPending).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByUserAndCode retrieves a verification by user and code
func (r *UserVerificationRepository) GetByUserAndCode(userID uuid.UUID, code string, verificationType VerificationType) (*UserVerification, error) {
	var verification UserVerification
	err := r.db.Where("user_id = ? AND code = ? AND type = ? AND status = ?", 
		userID, code, verificationType, VerificationStatusPending).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// Update updates a verification record
func (r *UserVerificationRepository) Update(verification *UserVerification) error {
	return r.db.Save(verification).Error
}

// GetPendingVerifications retrieves all pending verifications for admin review
func (r *UserVerificationRepository) GetPendingVerifications(verificationType *VerificationType, limit, offset int) ([]*UserVerification, error) {
	var verifications []*UserVerification
	query := r.db.Preload("User").Where("status = ?", VerificationStatusPending)
	
	if verificationType != nil {
		query = query.Where("type = ?", *verificationType)
	}
	
	err := query.Order("created_at ASC").Limit(limit).Offset(offset).Find(&verifications).Error
	return verifications, err
}

// GetUserVerifications retrieves all verifications for a user
func (r *UserVerificationRepository) GetUserVerifications(userID uuid.UUID) ([]*UserVerification, error) {
	var verifications []*UserVerification
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&verifications).Error
	return verifications, err
}

// ExpirePendingVerifications marks expired verifications as expired
func (r *UserVerificationRepository) ExpirePendingVerifications() error {
	return r.db.Model(&UserVerification{}).
		Where("status = ? AND expires_at < ?", VerificationStatusPending, time.Now()).
		Update("status", VerificationStatusExpired).Error
}
