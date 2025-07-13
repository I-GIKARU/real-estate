//go:build !noswagger
// +build !noswagger

package main

import (
	"fmt"

	"real-estate-backend/docs"
	"real-estate-backend/internal/config"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func setupSwagger(router *gin.Engine, cfg *config.Config) {
	// Swagger documentation
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
