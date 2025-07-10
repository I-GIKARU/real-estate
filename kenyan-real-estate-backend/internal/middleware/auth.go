package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"kenyan-real-estate-backend/internal/models"
	"kenyan-real-estate-backend/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must start with Bearer",
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}

// RequireUserType creates a middleware that requires specific user types
func RequireUserType(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User type not found in context",
			})
			c.Abort()
			return
		}

		userTypeStr, ok := userType.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user type format",
			})
			c.Abort()
			return
		}

		// Check if user type is allowed
		for _, allowedType := range allowedTypes {
			if userTypeStr == allowedType {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
		c.Abort()
	}
}

// UserRepositoryInterface defines the interface for user repository
type UserRepositoryInterface interface {
	GetByID(id uuid.UUID) (*models.User, error)
}

// RequireVerifiedEmail creates a middleware that requires users to have verified emails
func RequireVerifiedEmail(userRepo UserRepositoryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in context",
			})
			c.Abort()
			return
		}

		// Convert userID to UUID
		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID format",
			})
			c.Abort()
			return
		}

		// Get user from database to check verification status
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		// Check if user's email is verified
		if !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Email verification required. Please verify your email before accessing this resource.",
				"verification_required": true,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

