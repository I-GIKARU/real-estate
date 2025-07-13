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
	
	// Configure swagger
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	log.Println("Swagger documentation enabled at /swagger/index.html")
}
