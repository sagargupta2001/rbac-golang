package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort      string
	DatabaseURL     string
	JWTSecret       string
	JWTExpirationInHours int64
}

// LoadConfig loads configuration from .env file
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8080"
	}
    
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is not set")
	}

	jwtExpHours, err := strconv.ParseInt(os.Getenv("JWT_EXPIRATION_HOURS"), 10, 64)
	if err != nil {
		jwtExpHours = 72
	}

	return &Config{
		ServerPort:      serverPort,
		DatabaseURL:     databaseURL,
		JWTSecret:       jwtSecret,
		JWTExpirationInHours: jwtExpHours,
	}, nil
}