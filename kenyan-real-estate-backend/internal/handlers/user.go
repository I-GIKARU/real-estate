package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"kenyan-real-estate-backend/pkg/auth"
	"kenyan-real-estate-backend/internal/models"
	"kenyan-real-estate-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userRepo              *models.UserRepository
	jwtManager            *auth.JWTManager
	emailVerificationRepo *models.EmailVerificationRepository
	emailService          *services.EmailService
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	userRepo *models.UserRepository,
	jwtManager *auth.JWTManager,
	emailVerificationRepo *models.EmailVerificationRepository,
	emailService *services.EmailService,
) *UserHandler {
	return &UserHandler{
		userRepo:              userRepo,
		jwtManager:            jwtManager,
		emailVerificationRepo: emailVerificationRepo,
		emailService:          emailService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account with email, password, and profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User registration data"
// @Success 201 {object} object{message=string,user=models.UserResponse,token=string} "User created successfully"
// @Failure 400 {object} object{error=string,details=string} "Invalid request data"
// @Failure 409 {object} object{error=string} "Email or phone already exists"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Check if email already exists
	emailExists, err := h.userRepo.EmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check email existence",
		})
		return
	}
	if emailExists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already exists",
		})
		return
	}

	// Check if phone number already exists
	phoneExists, err := h.userRepo.PhoneExists(req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check phone number existence",
		})
		return
	}
	if phoneExists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Phone number already exists",
		})
		return
	}

	// Create user
	user := &models.User{
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		UserType:    req.UserType,
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Save user to database
	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email, string(user.UserType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Automatically send verification email
	go h.sendVerificationEmailAsync(user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Please check your email to verify your account.",
		"user":    user.ToResponse(),
		"token":   token,
		"email_verification_required": true,
	})
}

// sendVerificationEmailAsync sends a verification email asynchronously
// to the newly registered user.
func (h *UserHandler) sendVerificationEmailAsync(user *models.User) {
	// Check for rate limiting (max 1 email per 5 minutes)
	existingVerification, err := h.emailVerificationRepo.GetByUserID(user.ID)
	if err == nil && !existingVerification.IsUsed {
		timeSinceLastEmail := time.Since(existingVerification.CreatedAt)
		if timeSinceLastEmail < 5*time.Minute {
			return
		}
	}

	// Generate verification token
	token := services.GenerateSecureToken()

	// Create verification record
	verification := &models.EmailVerification{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
		IsUsed:    false,
	}

	h.emailVerificationRepo.Create(verification)

	// Send verification email
	fullName := user.FirstName + " " + user.LastName
	h.emailService.SendVerificationEmail(user.Email, fullName, token)
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User login credentials"
// @Success 200 {object} object{message=string,user=models.UserResponse,token=string} "Login successful"
// @Failure 400 {object} object{error=string,details=string} "Invalid request data"
// @Failure 401 {object} object{error=string} "Invalid credentials or account deactivated"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Account is deactivated",
		})
		return
	}

	// Check if user's email is verified
	if !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Email verification required. Please verify your email before logging in.",
			"verification_required": true,
			"user_email": user.Email,
		})
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email, string(user.UserType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user.ToResponse(),
		"token":   token,
	})
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object{user=models.UserResponse} "User profile"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 404 {object} object{error=string} "User not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}

// GetPendingAgents handles getting all agents waiting for approval (admin only)
// @Summary Get pending agents
// @Description Get all agents waiting for admin approval
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object{agents=[]models.UserResponse} "Pending agents"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "Forbidden - Admin access required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /admin/pending-agents [get]
func (h *UserHandler) GetPendingAgents(c *gin.Context) {
	// Check if user is admin
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Admin access required",
		})
		return
	}

	// Get pending agents
	agents, err := h.userRepo.GetPendingAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get pending agents",
		})
		return
	}

	// Convert to response format
	var agentResponses []models.UserResponse
	for _, agent := range agents {
		agentResponses = append(agentResponses, *agent.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"agents": agentResponses,
	})
}

// ApproveAgent handles approving an agent (admin only)
// @Summary Approve an agent
// @Description Approve an agent to allow property management
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param agentId path string true "Agent ID"
// @Success 200 {object} object{message=string,agent=models.UserResponse} "Agent approved successfully"
// @Failure 400 {object} object{error=string} "Invalid agent ID"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "Forbidden - Admin access required"
// @Failure 404 {object} object{error=string} "Agent not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /admin/approve-agent/{agentId} [post]
func (h *UserHandler) ApproveAgent(c *gin.Context) {
	// Check if user is admin
	adminID, isAdmin := h.getAdminID(c)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Admin access required",
		})
		return
	}

	// Get agent ID from URL
	agentIDStr := c.Param("agentId")
	agentID, err := uuid.Parse(agentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid agent ID",
		})
		return
	}

	// Approve the agent
	err = h.userRepo.ApproveAgent(agentID, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to approve agent",
			"details": err.Error(),
		})
		return
	}

	// Get updated agent data
	agent, err := h.userRepo.GetByID(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Agent not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Agent approved successfully",
		"agent":   agent.ToResponse(),
	})
}

// GetAllAgents handles getting all agents (admin only)
// @Summary Get all agents
// @Description Get all agents (approved and pending)
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object{agents=[]models.UserResponse} "All agents"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "Forbidden - Admin access required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /admin/agents [get]
func (h *UserHandler) GetAllAgents(c *gin.Context) {
	// Check if user is admin
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Admin access required",
		})
		return
	}

	// Get all agents
	agents, err := h.userRepo.GetAllAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get agents",
		})
		return
	}

	// Convert to response format
	var agentResponses []models.UserResponse
	for _, agent := range agents {
		agentResponses = append(agentResponses, *agent.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"agents": agentResponses,
	})
}

// Helper functions

// isAdmin checks if the current user is an admin
func (h *UserHandler) isAdmin(c *gin.Context) bool {
	userType, exists := c.Get("user_type")
	if !exists {
		return false
	}
	return userType == string(models.UserTypeAdmin)
}

// getAdminID gets the admin ID if the current user is an admin
func (h *UserHandler) getAdminID(c *gin.Context) (uuid.UUID, bool) {
	if !h.isAdmin(c) {
		return uuid.Nil, false
	}
	
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}
	
	return userUUID, true
}
