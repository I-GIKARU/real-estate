package main

import (
	"fmt"
	"log"
	"os"

	"real-estate-backend/docs"
	"real-estate-backend/internal/config"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func setupSwagger(router *gin.Engine, cfg *config.Config) {
	log.Println("Setting up Swagger documentation...")
	
	// Configure swagger host based on environment variable
	swaggerHost := os.Getenv("SWAGGER_HOST")
	if swaggerHost == "" {
		log.Println("Warning: SWAGGER_HOST not set. Using default host configuration.")
		swaggerHost = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	}
	docs.SwaggerInfo.Host = swaggerHost
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	log.Printf("Swagger documentation enabled at /swagger/index.html (host: %s)", swaggerHost)
}
