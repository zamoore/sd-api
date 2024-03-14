package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DBHost            string
	DBPort            int
	DBUser            string
	DBPassword        string
	DBName            string
	Auth0Domain       string
	Auth0ClientID     string
	Auth0ClientSecret string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	return Config{
		Port:              getEnv("PORT", "8080"),
		DBHost:            os.Getenv("HOST"),
		DBPort:            mustParseInt(os.Getenv("DB_PORT")),
		DBUser:            os.Getenv("USER"),
		DBPassword:        os.Getenv("PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		Auth0Domain:       os.Getenv("AUTH0_DOMAIN"),
		Auth0ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		Auth0ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func mustParseInt(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid integer value for %s: %v", s, err)
	}
	return value
}
