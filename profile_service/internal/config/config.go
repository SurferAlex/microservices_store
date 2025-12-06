package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	ServerPort     string
	AuthServiceURL string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:         getEnv("DBHost", "localhost"),
		DBPort:         getEnv("DBPort", "5432"),
		DBName:         getEnv("DBName", "profile_service"),
		DBUser:         getEnv("DBUser", "postgres"),
		DBPassword:     getEnv("DBPassword", "qweasdzxc"),
		AuthServiceURL: getEnv("AuthServiceURL", "http://localhost:8080"),
	}
}

func (c *Config) GetDBConnectionString() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		c.DBUser, c.DBPassword, c.DBName, c.DBHost, c.DBPort)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
