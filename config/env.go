package config

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func LoadEnvFile() {
	if err := godotenv.Load(); err != nil {
		log.Warn("Failed to load .env file: ", err)
		return
	}
	log.Info("Loaded .env file")
}

func getStringEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Error("Failed to convert environment variable to int", err)
			return 0, err
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Error("Failed to convert environment variable to bool", err)
			return false, err
		}
		return boolValue, nil
	}
	return defaultValue, nil
}

type DatabaseConfig struct {
	url string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		url: getStringEnv("DATABASE_URL", ""),
	}
}

type LogConfig struct {
	Level int
}

func NewLogConfig() (*LogConfig, error) {
	logLevelFallback := int(log.LevelTrace)
	level, err := getIntEnv("LOG_LEVEL", logLevelFallback)
	if err != nil {
		log.Error("Failed to get log level", err)
		return nil, err
	}
	return &LogConfig{Level: level}, nil
}
