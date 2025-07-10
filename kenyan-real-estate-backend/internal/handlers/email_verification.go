package handlers

import (
	"database/sql"
	"fmt"
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

// VerifyEmailGET verifies a user's email using the verification token via GET request
// @Summary Verify email address via GET
// @Description Verify email address using the verification token from URL parameter
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} object{message=string,user=models.UserResponse} "Email verified successfully"
// @Failure 400 {object} object{error=string} "Invalid token or expired"
// @Failure 404 {object} object{error=string} "Token not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /verify-email [get]
func (h *EmailVerificationHandler) VerifyEmailGET(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		errorHTML := h.generateErrorHTML("Missing Token", "The verification token is missing from the request. Please check your email and click the verification link again.")
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(errorHTML))
		return
	}

	h.verifyEmailWithToken(c, token)
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

	h.verifyEmailWithToken(c, req.Token)
}

// generateErrorHTML generates HTML for error responses
func (h *EmailVerificationHandler) generateErrorHTML(title, message string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
	    <meta charset="UTF-8">
	    <title>%s</title>
	    <style>
	        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
	        .header { background-color: #dc3545; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
	        .content { background-color: #f9f9f9; padding: 30px; border-radius: 0 0 5px 5px; text-align: center; }
	        .footer { margin-top: 30px; font-size: 12px; color: #666; text-align: center; }
	    </style>
	</head>
	<body>
	    <div class="header">
	        <h1>%s</h1>
	    </div>
	    <div class="content">
	        <p>%s</p>
	    </div>
	    <div class="footer">
	        <p>&copy; 2025 Kenyan Real Estate. All rights reserved.</p>
	    </div>
	</body>
	</html>
	`, title, title, message)
}

// verifyEmailWithToken handles the common verification logic
func (h *EmailVerificationHandler) verifyEmailWithToken(c *gin.Context, token string) {
	// Get verification record by token
	verification, err := h.emailVerificationRepo.GetByToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			errorHTML := h.generateErrorHTML("Invalid Token", "The verification token is invalid or does not exist.")
			c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(errorHTML))
			return
		}
		errorHTML := h.generateErrorHTML("Server Error", "An error occurred while processing your request. Please try again later.")
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(errorHTML))
		return
	}

	// Check if token is expired
	if verification.IsExpired() {
		errorHTML := h.generateErrorHTML("Token Expired", "The verification token has expired. Please request a new verification email.")
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(errorHTML))
		return
	}

	// Check if token is already used
	if verification.IsUsed {
		errorHTML := h.generateErrorHTML("Token Already Used", "This verification token has already been used. Your email may already be verified.")
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(errorHTML))
		return
	}

	// Mark verification as used
	if err := h.emailVerificationRepo.MarkAsUsed(verification.ID); err != nil {
		errorHTML := h.generateErrorHTML("Server Error", "An error occurred while processing your verification. Please try again later.")
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(errorHTML))
		return
	}

	// Update user's verified status
	user := verification.User
	user.IsVerified = true
	if err := h.userRepo.Update(&user); err != nil {
		errorHTML := h.generateErrorHTML("Server Error", "An error occurred while updating your account. Please try again later.")
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(errorHTML))
		return
	}

	// Send welcome email
	fullName := user.FirstName + " " + user.LastName
	go func() {
		// Send welcome email in background (don't fail the request if this fails)
		h.emailService.SendWelcomeEmail(user.Email, fullName)
	}()

successHTML := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
	    <meta charset="UTF-8">
	    <title>Email Verified</title>
	    <style>
	        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
	        .header { background-color: #28a745; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
	        .content { background-color: #f9f9f9; padding: 30px; border-radius: 0 0 5px 5px; text-align: center; }
	        .footer { margin-top: 30px; font-size: 12px; color: #666; text-align: center; }
	    </style>
	</head>
	<body>
	    <div class="header">
	        <h1>Email Verified Successfully</h1>
	    </div>
	    <div class="content">
	        <p>Hello %s,</p>
	        <p>Your email has been verified successfully! You can now start using your account.</p>
	    </div>
	    <div class="footer">
	        <p>&copy; 2025 Kenyan Real Estate. All rights reserved.</p>
	    </div>
	</body>
	</html>
	`, user.FirstName + " " + user.LastName)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(successHTML))
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
