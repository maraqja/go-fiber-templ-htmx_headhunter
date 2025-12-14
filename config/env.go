package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func LoadEnvFile() {
	if err := godotenv.Load(); err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to load .env file")
		return
	}
	log.Info().Msg("Loaded .env file")
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
			log.Error().
				Err(err).
				Str("key", key).
				Msg("Failed to convert environment variable to int")
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
			log.Error().
				Err(err).
				Str("key", key).
				Msg("Failed to convert environment variable to bool")
			return false, err
		}
		return boolValue, nil
	}
	return defaultValue, nil
}

type DatabaseConfig struct {
	Url string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Url: getStringEnv("DATABASE_URL", ""),
	}
}

type LogConfig struct {
	Level  int
	Output string
	Format string
}

func NewLogConfig() (*LogConfig, error) {
	level, err := getIntEnv("LOG_LEVEL", 0)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get log level")
		return nil, err
	}
	output := getStringEnv("LOG_OUTPUT", "stdout")
	format := getStringEnv("LOG_FORMAT", "json")
	return &LogConfig{
		Level:  level,
		Output: output,
		Format: format,
	}, nil
}
