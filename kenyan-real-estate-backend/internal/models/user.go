package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	ID              uuid.UUID  `json:"id" db:"id"`
	Email           string     `json:"email" db:"email"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	FirstName       string     `json:"first_name" db:"first_name"`
	LastName        string     `json:"last_name" db:"last_name"`
	PhoneNumber     string     `json:"phone_number" db:"phone_number"`
	UserType        UserType   `json:"user_type" db:"user_type"`
	IDNumber        *string    `json:"id_number,omitempty" db:"id_number"`
	ProfileImageURL *string    `json:"profile_image_url,omitempty" db:"profile_image_url"`
	IsVerified      bool       `json:"is_verified" db:"is_verified"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
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
	IDNumber        *string   `json:"id_number,omitempty"`
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
		IDNumber:        u.IDNumber,
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

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, phone_number, user_type, id_number, is_verified, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	user.ID = uuid.New()
	user.IsVerified = false
	user.IsActive = true
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.UserType,
		user.IDNumber,
		user.IsVerified,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone_number, user_type, id_number, profile_image_url, is_verified, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.UserType,
		&user.IDNumber,
		&user.ProfileImageURL,
		&user.IsVerified,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone_number, user_type, id_number, profile_image_url, is_verified, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.UserType,
		&user.IDNumber,
		&user.ProfileImageURL,
		&user.IsVerified,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3, phone_number = $4, id_number = $5, profile_image_url = $6, updated_at = $7
		WHERE id = $1`

	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.IDNumber,
		user.ProfileImageURL,
		user.UpdatedAt,
	)

	return err
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = $2 WHERE id = $1`
	_, err := r.db.Exec(query, id, time.Now())
	return err
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1 AND is_active = true`
	err := r.db.QueryRow(query, email).Scan(&count)
	return count > 0, err
}

// PhoneExists checks if a phone number already exists
func (r *UserRepository) PhoneExists(phone string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE phone_number = $1 AND is_active = true`
	err := r.db.QueryRow(query, phone).Scan(&count)
	return count > 0, err
}

