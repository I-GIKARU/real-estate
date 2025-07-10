package database

import (
	"fmt"
	"log"
	"time"

	"kenyan-real-estate-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// Connect establishes a connection to the database
func Connect(cfg *config.DatabaseConfig) error {
	var err error
	
	dsn := cfg.GetDSN()
	
	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)
	if cfg.Environment == "production" {
		gormLogger = logger.Default.LogMode(logger.Error)
	}
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}

// AutoMigrate runs database migrations
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}
	return DB.AutoMigrate(models...)
}

