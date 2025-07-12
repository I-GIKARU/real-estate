package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"real-estate-backend/internal/models"
	"real-estate-backend/internal/services"
)

// PasswordResetHandler handles password reset operations
type PasswordResetHandler struct {
	userRepo          *models.UserRepository
	passwordResetRepo *models.PasswordResetRepository
	emailService      *services.EmailService
}

// NewPasswordResetHandler creates a new password reset handler
func NewPasswordResetHandler(userRepo *models.UserRepository, passwordResetRepo *models.PasswordResetRepository, emailService *services.EmailService) *PasswordResetHandler {
	return &PasswordResetHandler{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		emailService:      emailService,
	}
}

// ForgotPassword handles password reset requests
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} map[string]interface{} "Password reset email sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/forgot-password [post]
func (h *PasswordResetHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Don't reveal if email exists or not for security
			c.JSON(http.StatusOK, gin.H{
				"message": "If your email is registered, you will receive a password reset link",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	// Generate reset token
	token, err := generateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	// Delete any existing tokens for this user
	err = h.passwordResetRepo.DeleteByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	// Create new password reset token (expires in 1 hour)
	passwordReset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err = h.passwordResetRepo.Create(passwordReset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reset token"})
		return
	}

	// Send password reset email
	err = h.emailService.SendPasswordResetEmail(user.Email, user.FirstName, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent successfully",
	})
}

// ResetPassword handles password reset with token
// @Summary Reset password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} map[string]interface{} "Password reset successful"
// @Failure 400 {object} map[string]interface{} "Invalid request or token"
// @Failure 404 {object} map[string]interface{} "Token not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/reset-password [post]
func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get password reset token
	passwordReset, err := h.passwordResetRepo.GetByTokenWithUser(req.Token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired reset token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate reset token"})
		return
	}

	// Check if token is expired
	if passwordReset.IsExpired() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reset token has expired"})
		return
	}

	// Check if token is already used
	if passwordReset.IsUsed() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reset token has already been used"})
		return
	}

	// Update user password
	user := &passwordReset.User
	err = user.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	err = h.userRepo.Update(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Mark token as used
	passwordReset.MarkAsUsed()
	err = h.passwordResetRepo.Update(passwordReset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reset token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successful",
	})
}

// ChangePassword handles password change for logged-in users
// @Summary Change password
// @Description Change password for authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "Change password request"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/change-password [post]
// @Security BearerAuth
func (h *PasswordResetHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Verify current password
	if !user.CheckPassword(req.CurrentPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Update password
	err = user.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	err = h.userRepo.Update(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// ValidateResetToken validates a reset token
// @Summary Validate reset token
// @Description Validate if a reset token is valid and not expired
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Reset token"
// @Success 200 {object} map[string]interface{} "Token is valid"
// @Failure 400 {object} map[string]interface{} "Invalid or expired token"
// @Failure 404 {object} map[string]interface{} "Token not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/validate-reset-token [get]
func (h *PasswordResetHandler) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	// Get password reset token
	passwordReset, err := h.passwordResetRepo.GetByToken(token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid reset token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
		return
	}

	// Check if token is expired
	if passwordReset.IsExpired() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reset token has expired"})
		return
	}

	// Check if token is already used
	if passwordReset.IsUsed() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reset token has already been used"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
		"expires_at": passwordReset.ExpiresAt,
	})
}

// CleanupExpiredTokens removes expired tokens (can be called periodically)
func (h *PasswordResetHandler) CleanupExpiredTokens() error {
	return h.passwordResetRepo.DeleteExpiredTokens()
}

// generateResetToken generates a secure random token for password reset
func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetResetPasswordForm serves the HTML form for password reset
// @Summary Get password reset form
// @Description Serve HTML form for password reset
// @Tags web
// @Accept html
// @Produce html
// @Param token query string true "Reset token"
// @Success 200 {string} string "HTML form"
// @Failure 400 {string} string "Invalid token"
// @Router /web/reset-password [get]
func (h *PasswordResetHandler) GetResetPasswordForm(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(getErrorHTML("Invalid reset token")))
		return
	}

	// Validate token
	passwordReset, err := h.passwordResetRepo.GetByToken(token)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(getErrorHTML("Invalid or expired reset token")))
		return
	}

	if passwordReset.IsExpired() || passwordReset.IsUsed() {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(getErrorHTML("Reset token has expired or already been used")))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(getResetPasswordFormHTML(token)))
}

// PostResetPasswordForm handles form submission for password reset
// @Summary Handle password reset form submission
// @Description Process password reset form submission
// @Tags web
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param token formData string true "Reset token"
// @Param password formData string true "New password"
// @Param confirm_password formData string true "Confirm new password"
// @Success 200 {string} string "Success page"
// @Failure 400 {string} string "Error page"
// @Router /web/reset-password [post]
func (h *PasswordResetHandler) PostResetPasswordForm(c *gin.Context) {
	token := c.PostForm("token")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// Validate input
	if token == "" || password == "" || confirmPassword == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(getErrorHTML("All fields are required")))
		return
	}

	if password != confirmPassword {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(getErrorHTML("Passwords do not match")))
		return
	}

	if len(password) < 8 {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(getErrorHTML("Password must be at least 8 characters long")))
		return
	}

	// Create reset request
	req := models.ResetPasswordRequest{
		Token:           token,
		Password:        password,
		ConfirmPassword: confirmPassword,
	}

	// Process reset using the same logic as API
	passwordReset, err := h.passwordResetRepo.GetByTokenWithUser(req.Token)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(getErrorHTML("Invalid or expired reset token")))
		return
	}

	if passwordReset.IsExpired() || passwordReset.IsUsed() {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(getErrorHTML("Reset token has expired or already been used")))
		return
	}

	// Update user password
	user := &passwordReset.User
	err = user.HashPassword(req.Password)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(getErrorHTML("Failed to process password reset")))
		return
	}

	err = h.userRepo.Update(user)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(getErrorHTML("Failed to update password")))
		return
	}

	// Mark token as used
	passwordReset.MarkAsUsed()
	err = h.passwordResetRepo.Update(passwordReset)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(getErrorHTML("Failed to process reset")))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(getSuccessHTML("Password reset successful! You can now login with your new password.")))
}

// HTML template functions

func getResetPasswordFormHTML(token string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password - Real Estate Platform</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 20px;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        .container {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 40px;
            width: 100%%;
            max-width: 400px;
        }
        .logo {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo h1 {
            color: #2c3e50;
            font-size: 24px;
            margin: 0;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: 600;
            color: #333;
        }
        input[type="password"] {
            width: 100%%;
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
            box-sizing: border-box;
        }
        input[type="password"]:focus {
            outline: none;
            border-color: #3498db;
            box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
        }
        .btn {
            width: 100%%;
            padding: 12px;
            background-color: #3498db;
            color: white;
            border: none;
            border-radius: 4px;
            font-size: 16px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        .btn:hover {
            background-color: #2980b9;
        }
        .requirements {
            font-size: 12px;
            color: #666;
            margin-top: 5px;
        }
        .title {
            text-align: center;
            margin-bottom: 30px;
            color: #2c3e50;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <h1>🏠 Real Estate Platform</h1>
        </div>
        <h2 class="title">Reset Your Password</h2>
        <form method="POST" action="/web/reset-password">
            <input type="hidden" name="token" value="%s">
            
            <div class="form-group">
                <label for="password">New Password</label>
                <input type="password" id="password" name="password" required>
                <div class="requirements">Password must be at least 8 characters long</div>
            </div>
            
            <div class="form-group">
                <label for="confirm_password">Confirm New Password</label>
                <input type="password" id="confirm_password" name="confirm_password" required>
            </div>
            
            <button type="submit" class="btn">Reset Password</button>
        </form>
    </div>
</body>
</html>
`, token)
}

func getErrorHTML(message string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Error - Real Estate Platform</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 20px;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        .container {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 40px;
            width: 100%%;
            max-width: 400px;
            text-align: center;
        }
        .logo {
            margin-bottom: 30px;
        }
        .logo h1 {
            color: #2c3e50;
            font-size: 24px;
            margin: 0;
        }
        .error-icon {
            font-size: 48px;
            color: #e74c3c;
            margin-bottom: 20px;
        }
        .error-message {
            color: #e74c3c;
            font-size: 18px;
            margin-bottom: 30px;
            line-height: 1.5;
        }
        .btn {
            display: inline-block;
            padding: 12px 24px;
            background-color: #3498db;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            font-size: 16px;
            transition: background-color 0.3s;
        }
        .btn:hover {
            background-color: #2980b9;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <h1>🏠 Real Estate Platform</h1>
        </div>
        <div class="error-icon">⚠️</div>
        <div class="error-message">%s</div>
        <a href="/" class="btn">Go to Homepage</a>
    </div>
</body>
</html>
`, message)
}

func getSuccessHTML(message string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Success - Real Estate Platform</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 20px;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        .container {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 40px;
            width: 100%%;
            max-width: 400px;
            text-align: center;
        }
        .logo {
            margin-bottom: 30px;
        }
        .logo h1 {
            color: #2c3e50;
            font-size: 24px;
            margin: 0;
        }
        .success-icon {
            font-size: 48px;
            color: #27ae60;
            margin-bottom: 20px;
        }
        .success-message {
            color: #27ae60;
            font-size: 18px;
            margin-bottom: 30px;
            line-height: 1.5;
        }
        .btn {
            display: inline-block;
            padding: 12px 24px;
            background-color: #27ae60;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            font-size: 16px;
            transition: background-color 0.3s;
        }
        .btn:hover {
            background-color: #219a52;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <h1>🏠 Real Estate Platform</h1>
        </div>
        <div class="success-icon">✅</div>
        <div class="success-message">%s</div>
        <a href="/" class="btn">Continue to Login</a>
    </div>
</body>
</html>
`, message)
}
