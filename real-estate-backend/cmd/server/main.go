package main

import (
	"fmt"
	"log"

	"real-estate-backend/docs"
	"real-estate-backend/internal/config"
	"real-estate-backend/internal/database"
	"real-estate-backend/internal/handlers"
	"real-estate-backend/internal/middleware"
	"real-estate-backend/internal/models"
	"real-estate-backend/internal/services"
	"real-estate-backend/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title           Real Estate API
// @version         1.0
// @description     A comprehensive API for managing real estate properties
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Run auto-migration
	log.Println("Running database migrations...")
	if err := database.AutoMigrate(
		&models.User{},
		&models.County{},
		&models.SubCounty{},
		&models.Property{},
		&models.PropertyImage{},
		&models.EmailVerification{},
		&models.PasswordReset{},
		// Add other models here as needed
	); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(&cfg.JWT)

	// Initialize repositories
	userRepo := models.NewUserRepository(database.GetDB())
	propertyRepo := models.NewPropertyRepository(database.GetDB())
	countyRepo := models.NewCountyRepository(database.GetDB())
	subCountyRepo := models.NewSubCountyRepository(database.GetDB())
	propertyImageRepo := models.NewPropertyImageRepository(database.GetDB())
	emailVerificationRepo := models.NewEmailVerificationRepository(database.GetDB())
	passwordResetRepo := models.NewPasswordResetRepository(database.GetDB())
	// rentalApplicationRepo := models.NewRentalApplicationRepository(database.GetDB())
	// leaseRepo := models.NewLeaseRepository(database.GetDB())
	// paymentRepo := models.NewPaymentRepository(database.GetDB())

	// Initialize services
	// mpesaService := services.NewMPesaService(&cfg.MPesa)
	cloudinaryService, err := services.NewCloudinaryService(&cfg.Cloudinary)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary service:", err)
	}
	emailService := services.NewEmailService(&cfg.Email)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo, jwtManager, emailVerificationRepo, emailService)
	propertyHandler := handlers.NewPropertyHandler(propertyRepo, propertyImageRepo, cloudinaryService, &cfg.Upload)
	locationHandler := handlers.NewLocationHandler(countyRepo, subCountyRepo)
	emailVerificationHandler := handlers.NewEmailVerificationHandler(userRepo, emailVerificationRepo, emailService)
	passwordResetHandler := handlers.NewPasswordResetHandler(userRepo, passwordResetRepo, emailService)
	// paymentHandler := handlers.NewPaymentHandler(paymentRepo, leaseRepo, mpesaService)

	// Set up Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Swagger documentation
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "real-estate-backend",
		})
	})

	// Web forms for password reset (not under /api/v1)
	web := router.Group("/web")
	{
		web.GET("/reset-password", passwordResetHandler.GetResetPasswordForm)
		web.POST("/reset-password", passwordResetHandler.PostResetPasswordForm)
	}

	// API routes
	api := router.Group("/api/v1")

	// Public routes (no authentication required)
	public := api.Group("/")
	{
		// User authentication
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)

		// Public property listings
		public.GET("/properties", propertyHandler.GetPublicProperties)
		public.GET("/properties/:id", propertyHandler.GetProperty)

		// Location data
		public.GET("/counties", locationHandler.GetCounties)
		public.GET("/counties/:id", locationHandler.GetCounty)
		public.GET("/counties/:id/sub-counties", locationHandler.GetSubCounties) // updated
		public.GET("/sub-counties/:id", locationHandler.GetSubCounty)

		// Email verification (public)
		public.POST("/verify-email", emailVerificationHandler.VerifyEmail)
		public.GET("/verify-email", emailVerificationHandler.VerifyEmailGET)

		// Password reset (public)
		public.POST("/auth/forgot-password", passwordResetHandler.ForgotPassword)
		public.POST("/auth/reset-password", passwordResetHandler.ResetPassword)
		public.GET("/auth/validate-reset-token", passwordResetHandler.ValidateResetToken)

		// M-Pesa callback (public endpoint for Safaricom)
		// public.POST("/payments/mpesa/callback", paymentHandler.HandleMPesaCallback)
	}

	// Protected routes (authentication required)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		// User profile
		protected.GET("/profile", userHandler.GetProfile)

		// Email verification (protected)
		protected.POST("/send-verification-email", emailVerificationHandler.SendVerificationEmail)
		protected.GET("/verification-status", emailVerificationHandler.GetVerificationStatus)

		// Password management (protected)
		protected.POST("/auth/change-password", passwordResetHandler.ChangePassword)

		// Admin routes - admin access required
		adminRoutes := protected.Group("/admin")
		adminRoutes.Use(middleware.RequireUserType("admin"))
		{
			adminRoutes.GET("/pending-agents", userHandler.GetPendingAgents)
			adminRoutes.POST("/approve-agent/:agentId", userHandler.ApproveAgent)
			adminRoutes.GET("/agents", userHandler.GetAllAgents)
		}

		// Property management (agent only) - requires email verification and admin approval
		agentRoutes := protected.Group("/")
		agentRoutes.Use(middleware.RequireUserType("agent"))
		agentRoutes.Use(middleware.RequireVerifiedEmail(userRepo))
		agentRoutes.Use(middleware.RequireApprovedAgent(userRepo))
		{
			agentRoutes.POST("/properties", propertyHandler.CreateProperty)
			agentRoutes.PUT("/properties/:id", propertyHandler.UpdateProperty)
			agentRoutes.DELETE("/properties/:id", propertyHandler.DeleteProperty)
			agentRoutes.GET("/my-properties", propertyHandler.GetMyProperties)
			agentRoutes.POST("/properties/:id/images", propertyHandler.AddPropertyImage)
			agentRoutes.DELETE("/properties/:id/images/:image_id", propertyHandler.DeletePropertyImage)
		}

		// Tenant routes - requires email verification for applications and payments
		tenantRoutes := protected.Group("/")
		tenantRoutes.Use(middleware.RequireUserType("tenant"))
		tenantRoutes.Use(middleware.RequireVerifiedEmail(userRepo))
		{
			// tenantRoutes.POST("/applications", applicationHandler.CreateApplication)
			// tenantRoutes.GET("/my-applications", applicationHandler.GetMyApplications)
			// tenantRoutes.GET("/my-leases", leaseHandler.GetMyLeases)

			// Payment routes for tenants
			// tenantRoutes.POST("/payments/initiate", paymentHandler.InitiateRentPayment)
			// tenantRoutes.GET("/payments/status/:checkout_request_id", paymentHandler.QueryPaymentStatus)
		}

		// Payment routes accessible by both landlords and tenants
		// protected.GET("/payments/lease/:lease_id", paymentHandler.GetPaymentsByLease)

		// Routes accessible by both landlords and tenants
		// protected.POST("/properties/:id/favorite", favoriteHandler.AddFavorite)
		// protected.DELETE("/properties/:id/favorite", favoriteHandler.RemoveFavorite)
		// protected.GET("/favorites", favoriteHandler.GetFavorites)
	}

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
