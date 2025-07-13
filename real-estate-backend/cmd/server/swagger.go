package main

import (
	"log"

	"real-estate-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func setupSwagger(router *gin.Engine, cfg *config.Config) {
	// Swagger is disabled in Docker builds
	log.Println("Swagger setup skipped - disabled in Docker builds")
}
