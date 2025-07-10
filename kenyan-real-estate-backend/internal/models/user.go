package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserType represents the type of user
type UserType string

const (
	UserTypeLandlord UserType = "landlord"
	UserTypeTenant   UserType = "tenant"
	UserTypeAgent    UserType = "agent"
)

// User represents a user in the system
type User struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email           string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash    string     `json:"-" gorm:"not null"`
	FirstName       string     `json:"first_name" gorm:"not null"`
	LastName        string     `json:"last_name" gorm:"not null"`
	PhoneNumber     string     `json:"phone_number" gorm:"uniqueIndex;not null"`
	UserType        UserType   `json:"user_type" gorm:"not null;type:varchar(20)"`
	ProfileImageURL *string    `json:"profile_image_url,omitempty"`
	IsVerified      bool       `json:"is_verified" gorm:"default:false"`
	IsActive        bool       `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Properties []Property `json:"properties,omitempty" gorm:"foreignKey:LandlordID"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=8"`
	FirstName   string   `json:"first_name" binding:"required"`
	LastName    string   `json:"last_name" binding:"required"`
	PhoneNumber string   `json:"phone_number" binding:"required"`
	UserType    UserType `json:"user_type" binding:"required"`
	IDNumber    *string  `json:"id_number,omitempty"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the user response (without sensitive data)
type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	PhoneNumber     string    `json:"phone_number"`
	UserType        UserType  `json:"user_type"`
	ProfileImageURL *string   `json:"profile_image_url,omitempty"`
	IsVerified      bool      `json:"is_verified"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		PhoneNumber:     u.PhoneNumber,
		UserType:        u.UserType,
		ProfileImageURL: u.ProfileImageURL,
		IsVerified:      u.IsVerified,
		IsActive:        u.IsActive,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

// HashPassword hashes the user's password
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// BeforeCreate GORM hook to set ID and hash password
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*User, error) {
	var user User
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&User{}, id).Error
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("email = ? AND is_active = ?", email, true).Count(&count).Error
	return count > 0, err
}

// PhoneExists checks if a phone number already exists
func (r *UserRepository) PhoneExists(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("phone_number = ? AND is_active = ?", phone, true).Count(&count).Error
	return count > 0, err
}

