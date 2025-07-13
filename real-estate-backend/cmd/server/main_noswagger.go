//go:build noswagger
// +build noswagger

package main

import (
	"real-estate-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func setupSwagger(router *gin.Engine, cfg *config.Config) {
	// Swagger is disabled in this build
}
