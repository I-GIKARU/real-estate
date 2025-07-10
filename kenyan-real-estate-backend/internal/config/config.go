package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Upload    UploadConfig
	Cloudinary CloudinaryConfig
	MPesa     MPesaConfig
	Email     EmailConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string
	Port int
	Env  string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	DBName      string
	SSLMode     string
	Environment string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret     string
	ExpiryHours int
}

// UploadConfig holds file upload configuration
type UploadConfig struct {
	MaxFileSize  int64
	AllowedTypes []string
	UploadDir    string
}

// CloudinaryConfig holds Cloudinary service configuration
type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
	Folder    string
}

// MPesaConfig holds M-Pesa integration configuration
type MPesaConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	Environment    string // sandbox or production
	PassKey        string
	ShortCode      string
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	FromEmail    string
	SupportEmail string
	BaseURL      string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnvAsInt("DB_PORT", 5432),
			User:        getEnv("DB_USER", "postgres"),
			Password:    getEnv("DB_PASSWORD", ""),
			DBName:      getEnv("DB_NAME", "kenyan_real_estate"),
			SSLMode:     getEnv("DB_SSL_MODE", "disable"),
			Environment: getEnv("APP_ENV", "development"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		Upload: UploadConfig{
			MaxFileSize:  getEnvAsInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB
			AllowedTypes: []string{"image/jpeg", "image/png", "image/webp", "image/jpg"},
			UploadDir:    getEnv("UPLOAD_DIR", "./uploads"),
		},
		Cloudinary: CloudinaryConfig{
			CloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
			APIKey:    getEnv("CLOUDINARY_API_KEY", ""),
			APISecret: getEnv("CLOUDINARY_API_SECRET", ""),
			Folder:    getEnv("CLOUDINARY_FOLDER", "real-estate-properties"),
		},
		MPesa: MPesaConfig{
			ConsumerKey:    getEnv("MPESA_CONSUMER_KEY", ""),
			ConsumerSecret: getEnv("MPESA_CONSUMER_SECRET", ""),
			Environment:    getEnv("MPESA_ENVIRONMENT", "sandbox"),
			PassKey:        getEnv("MPESA_PASS_KEY", ""),
			ShortCode:      getEnv("MPESA_SHORT_CODE", ""),
		},
		Email: EmailConfig{
			Host:         getEnv("EMAIL_HOST", "smtp.gmail.com"),
			Port:         getEnvAsInt("EMAIL_PORT", 587),
			Username:     getEnv("EMAIL_USERNAME", ""),
			Password:     getEnv("EMAIL_PASSWORD", ""),
			FromEmail:    getEnv("EMAIL_FROM", "noreply@kenyanrealestate.com"),
			SupportEmail: getEnv("EMAIL_SUPPORT", "support@kenyanrealestate.com"),
			BaseURL:      getEnv("BASE_URL", "http://localhost:3000"),
		},
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

