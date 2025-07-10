package main

import (
	"fmt"
	"log"

	"kenyan-real-estate-backend/internal/config"
	"kenyan-real-estate-backend/internal/database"
	"kenyan-real-estate-backend/internal/handlers"
	"kenyan-real-estate-backend/internal/middleware"
	"kenyan-real-estate-backend/internal/models"
	"kenyan-real-estate-backend/internal/services"
	"kenyan-real-estate-backend/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(&cfg.JWT)

	// Initialize repositories
	userRepo := models.NewUserRepository(database.GetDB())
	propertyRepo := models.NewPropertyRepository(database.GetDB())
	countyRepo := models.NewCountyRepository(database.GetDB())
	subCountyRepo := models.NewSubCountyRepository(database.GetDB())
	propertyImageRepo := models.NewPropertyImageRepository(database.GetDB())
	// rentalApplicationRepo := models.NewRentalApplicationRepository(database.GetDB())
	leaseRepo := models.NewLeaseRepository(database.GetDB())
	paymentRepo := models.NewPaymentRepository(database.GetDB())

	// Initialize services
	mpesaService := services.NewMPesaService(&cfg.MPesa)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo, jwtManager)
	propertyHandler := handlers.NewPropertyHandler(propertyRepo, propertyImageRepo)
	locationHandler := handlers.NewLocationHandler(countyRepo, subCountyRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentRepo, leaseRepo, mpesaService)
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
		public.GET("/counties/:county_id/sub-counties", locationHandler.GetSubCounties)
		public.GET("/sub-counties/:id", locationHandler.GetSubCounty)

		// Kenyan-specific features
		public.GET("/amenities", kenyanFeaturesHandler.GetAmenities)
		public.GET("/property-types", kenyanFeaturesHandler.GetPropertyTypes)
		public.GET("/utilities", kenyanFeaturesHandler.GetUtilities)
		public.GET("/rental-terms", kenyanFeaturesHandler.GetRentalTerms)
		public.GET("/popular-areas", kenyanFeaturesHandler.GetPopularAreas)
		public.POST("/validate-phone", kenyanFeaturesHandler.ValidatePhoneNumber)
		public.POST("/format-currency", kenyanFeaturesHandler.FormatCurrency)

		// M-Pesa callback (public endpoint for Safaricom)
		public.POST("/payments/mpesa/callback", paymentHandler.HandleMPesaCallback)
	}

	// Protected routes (authentication required)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		// User profile
		protected.GET("/profile", userHandler.GetProfile)
		protected.PUT("/profile", userHandler.UpdateProfile)

		// Property management (landlord only)
		landlordRoutes := protected.Group("/")
		landlordRoutes.Use(middleware.RequireUserType("landlord"))
		{
			landlordRoutes.POST("/properties", propertyHandler.CreateProperty)
			landlordRoutes.PUT("/properties/:id", propertyHandler.UpdateProperty)
			landlordRoutes.DELETE("/properties/:id", propertyHandler.DeleteProperty)
			landlordRoutes.GET("/my-properties", propertyHandler.GetMyProperties)
			landlordRoutes.POST("/properties/:id/images", propertyHandler.AddPropertyImage)
		}

		// Tenant routes
		tenantRoutes := protected.Group("/")
		tenantRoutes.Use(middleware.RequireUserType("tenant"))
		{
			// tenantRoutes.POST("/applications", applicationHandler.CreateApplication)
			// tenantRoutes.GET("/my-applications", applicationHandler.GetMyApplications)
			// tenantRoutes.GET("/my-leases", leaseHandler.GetMyLeases)
			
			// Payment routes for tenants
			tenantRoutes.POST("/payments/initiate", paymentHandler.InitiateRentPayment)
			tenantRoutes.GET("/payments/status/:checkout_request_id", paymentHandler.QueryPaymentStatus)
		}

		// Payment routes accessible by both landlords and tenants
		protected.GET("/payments/lease/:lease_id", paymentHandler.GetPaymentsByLease)

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

