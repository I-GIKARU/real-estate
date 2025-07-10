package main

import (
	"fmt"
	"log"

	// "kenyan-real-estate-backend/docs"
	"kenyan-real-estate-backend/internal/config"
	"kenyan-real-estate-backend/internal/database"
	"kenyan-real-estate-backend/internal/handlers"
	"kenyan-real-estate-backend/internal/middleware"
	"kenyan-real-estate-backend/internal/models"
	"kenyan-real-estate-backend/internal/services"
	"kenyan-real-estate-backend/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// swaggerFiles "github.com/swaggo/files"
)

// @title           Kenyan Real Estate API
// @version         1.0
// @description     A comprehensive API for managing real estate properties in Kenya
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
	userHandler := handlers.NewUserHandler(userRepo, jwtManager)
	propertyHandler := handlers.NewPropertyHandler(propertyRepo, propertyImageRepo, cloudinaryService, &cfg.Upload)
	locationHandler := handlers.NewLocationHandler(countyRepo, subCountyRepo)
	emailVerificationHandler := handlers.NewEmailVerificationHandler(userRepo, emailVerificationRepo, emailService)
	// paymentHandler := handlers.NewPaymentHandler(paymentRepo, leaseRepo, mpesaService)
	kenyanFeaturesHandler := handlers.NewKenyanFeaturesHandler()

	// Set up Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Swagger documentation (disabled for Docker build)
	// docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "kenyan-real-estate-backend",
		})
	})

	// API routes
	api := router.Group("/api/v1")

	// Public routes (no authentication required)
	public := api.Group("/")
	{
		// User authentication
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
		public.POST("/refresh-token", userHandler.RefreshToken)

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
		public.GET("/verify-email", emailVerificationHandler.VerifyEmailByToken)

		// Kenyan-specific features
		public.GET("/amenities", kenyanFeaturesHandler.GetAmenities)
		public.GET("/property-types", kenyanFeaturesHandler.GetPropertyTypes)
		public.GET("/utilities", kenyanFeaturesHandler.GetUtilities)
		public.GET("/rental-terms", kenyanFeaturesHandler.GetRentalTerms)
		public.GET("/popular-areas", kenyanFeaturesHandler.GetPopularAreas)
		public.POST("/validate-phone", kenyanFeaturesHandler.ValidatePhoneNumber)
		public.POST("/format-currency", kenyanFeaturesHandler.FormatCurrency)

		// M-Pesa callback (public endpoint for Safaricom)
		// public.POST("/payments/mpesa/callback", paymentHandler.HandleMPesaCallback)
	}

	// Protected routes (authentication required)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		// User profile
		protected.GET("/profile", userHandler.GetProfile)
		protected.PUT("/profile", userHandler.UpdateProfile)

		// Email verification (protected)
		protected.POST("/send-verification-email", emailVerificationHandler.SendVerificationEmail)
		protected.GET("/verification-status", emailVerificationHandler.GetVerificationStatus)

		// Property management (landlord only)
		landlordRoutes := protected.Group("/")
		landlordRoutes.Use(middleware.RequireUserType("landlord"))
		{
			landlordRoutes.POST("/properties", propertyHandler.CreateProperty)
			landlordRoutes.PUT("/properties/:id", propertyHandler.UpdateProperty)
			landlordRoutes.DELETE("/properties/:id", propertyHandler.DeleteProperty)
			landlordRoutes.GET("/my-properties", propertyHandler.GetMyProperties)
			landlordRoutes.POST("/properties/:id/images", propertyHandler.AddPropertyImage)
			landlordRoutes.DELETE("/properties/:id/images/:image_id", propertyHandler.DeletePropertyImage)
		}

		// Tenant routes
		tenantRoutes := protected.Group("/")
		tenantRoutes.Use(middleware.RequireUserType("tenant"))
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
