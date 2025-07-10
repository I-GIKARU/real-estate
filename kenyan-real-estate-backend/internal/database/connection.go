package database

import (
	"database/sql"
	"fmt"
	"log"

	"kenyan-real-estate-backend/internal/config"

	_ "github.com/lib/pq"
)

// DB holds the database connection
var DB *sql.DB

// Connect establishes a connection to the database
func Connect(cfg *config.DatabaseConfig) error {
	var err error
	
	dsn := cfg.GetDSN()
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	log.Println("Database connection established successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}

