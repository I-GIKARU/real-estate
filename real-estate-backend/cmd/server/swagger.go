package main

import (
	"fmt"
	"log"

	"real-estate-backend/docs"
	"real-estate-backend/internal/config"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func setupSwagger(router *gin.Engine, cfg *config.Config) {
	log.Println("Setting up Swagger documentation...")
	
	// Configure swagger host based on environment
	var swaggerHost string
	if cfg.Server.Env == "production" {
		// For production, use the actual deployed URL
		swaggerHost = "real-estate-backend-840370620772.us-central1.run.app"
	} else {
		// For development, use localhost
		swaggerHost = fmt.Sprintf("localhost:%d", cfg.Server.Port)
	}
	
	docs.SwaggerInfo.Host = swaggerHost
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	log.Printf("Swagger documentation enabled at /swagger/index.html (host: %s)", swaggerHost)
}
