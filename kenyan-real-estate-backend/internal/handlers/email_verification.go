package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"kenyan-real-estate-backend/internal/models"
	"kenyan-real-estate-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EmailVerificationHandler handles email verification HTTP requests
type EmailVerificationHandler struct {
	userRepo                *models.UserRepository
	emailVerificationRepo   *models.EmailVerificationRepository
	emailService            *services.EmailService
}

// NewEmailVerificationHandler creates a new email verification handler
func NewEmailVerificationHandler(
	userRepo *models.UserRepository,
	emailVerificationRepo *models.EmailVerificationRepository,
	emailService *services.EmailService,
) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
		emailService:          emailService,
	}
}

// SendVerificationEmail sends a verification email to the user
// @Summary Send email verification
// @Description Send a verification email to the authenticated user
// @Tags Email Verification
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object{message=string} "Verification email sent"
// @Failure 400 {object} object{error=string} "User already verified"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 429 {object} object{error=string} "Too many requests"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /send-verification-email [post]
func (h *EmailVerificationHandler) SendVerificationEmail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Get user
	user, err := h.userRepo.GetByID(userUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	// Check if user is already verified
	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is already verified",
		})
		return
	}

	// Check for rate limiting (max 1 email per 5 minutes)
	existingVerification, err := h.emailVerificationRepo.GetByUserID(userUUID)
	if err == nil && !existingVerification.IsUsed {
		timeSinceLastEmail := time.Since(existingVerification.CreatedAt)
		if timeSinceLastEmail < 5*time.Minute {
			remainingTime := 5*time.Minute - timeSinceLastEmail
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Please wait before requesting another verification email",
				"retry_after_seconds": int(remainingTime.Seconds()),
			})
			return
		}
	}

	// Generate verification token
	token := services.GenerateSecureToken()
	
	// Create verification record
	verification := &models.EmailVerification{
		UserID:    userUUID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
		IsUsed:    false,
	}

	if err := h.emailVerificationRepo.Create(verification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create verification record",
		})
		return
	}

	// Send verification email
	fullName := user.FirstName + " " + user.LastName
	if err := h.emailService.SendVerificationEmail(user.Email, fullName, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send verification email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
	})
}

// VerifyEmail verifies a user's email using the verification token
// @Summary Verify email address
// @Description Verify email address using the verification token
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param verification body models.VerifyEmailRequest true "Verification token"
// @Success 200 {object} object{message=string,user=models.UserResponse} "Email verified successfully"
// @Failure 400 {object} object{error=string} "Invalid token or expired"
// @Failure 404 {object} object{error=string} "Token not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /verify-email [post]
func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
	var req models.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Get verification record by token
	verification, err := h.emailVerificationRepo.GetByToken(req.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Invalid verification token",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get verification record",
		})
		return
	}

	// Check if token is expired
	if verification.IsExpired() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verification token has expired",
		})
		return
	}

	// Check if token is already used
	if verification.IsUsed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verification token has already been used",
		})
		return
	}

	// Mark verification as used
	if err := h.emailVerificationRepo.MarkAsUsed(verification.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to mark verification as used",
		})
		return
	}

	// Update user's verified status
	user := verification.User
	user.IsVerified = true
	if err := h.userRepo.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user verification status",
		})
		return
	}

	// Send welcome email
	fullName := user.FirstName + " " + user.LastName
	go func() {
		// Send welcome email in background (don't fail the request if this fails)
		h.emailService.SendWelcomeEmail(user.Email, fullName)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
		"user":    user.ToResponse(),
	})
}

// VerifyEmailByToken verifies a user's email using a token from query parameter (GET request)
// @Summary Verify email address via GET
// @Description Verify email address using the verification token from query parameter
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} object{message=string,user=models.UserResponse} "Email verified successfully"
// @Failure 400 {object} object{error=string} "Invalid token or expired"
// @Failure 404 {object} object{error=string} "Token not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /verify-email [get]
func (h *EmailVerificationHandler) VerifyEmailByToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is required",
		})
		return
	}

	// Get verification record by token
	verification, err := h.emailVerificationRepo.GetByToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Invalid verification token",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get verification record",
		})
		return
	}

	// Check if token is expired
	if verification.IsExpired() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verification token has expired",
		})
		return
	}

	// Check if token is already used
	if verification.IsUsed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verification token has already been used",
		})
		return
	}

	// Mark verification as used
	if err := h.emailVerificationRepo.MarkAsUsed(verification.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to mark verification as used",
		})
		return
	}

	// Update user's verified status
	user := verification.User
	user.IsVerified = true
	if err := h.userRepo.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user verification status",
		})
		return
	}

	// Send welcome email
	fullName := user.FirstName + " " + user.LastName
	go func() {
		// Send welcome email in background (don't fail the request if this fails)
		h.emailService.SendWelcomeEmail(user.Email, fullName)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
		"user":    user.ToResponse(),
	})
}

// GetVerificationStatus gets the email verification status for the current user
// @Summary Get verification status
// @Description Get the email verification status for the authenticated user
// @Tags Email Verification
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object{is_verified=bool,pending_verification=bool,can_resend=bool} "Verification status"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /verification-status [get]
func (h *EmailVerificationHandler) GetVerificationStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Get user
	user, err := h.userRepo.GetByID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	response := gin.H{
		"is_verified": user.IsVerified,
		"pending_verification": false,
		"can_resend": !user.IsVerified,
	}

	// If not verified, check for pending verifications
	if !user.IsVerified {
		verification, err := h.emailVerificationRepo.GetByUserID(userUUID)
		if err == nil && !verification.IsUsed && !verification.IsExpired() {
			response["pending_verification"] = true
			
			// Check if user can resend (5 minute cooldown)
			timeSinceLastEmail := time.Since(verification.CreatedAt)
			if timeSinceLastEmail < 5*time.Minute {
				response["can_resend"] = false
				response["retry_after_seconds"] = int((5*time.Minute - timeSinceLastEmail).Seconds())
			}
		}
	}

	c.JSON(http.StatusOK, response)
}
