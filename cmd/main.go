package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/config"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/internal/home"
)

func main() {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	log.SetLevel(log.LevelTrace)

	config.LoadEnvFile()

	logConfig, err := config.NewLogConfig()
	if err != nil {
		log.Error("Failed to get log config", err)
		return
	}
	log.Info("Log config: ", logConfig)
	log.SetLevel(log.Level(logConfig.Level))

	databaseConfig := config.NewDatabaseConfig()
	log.Info("Database: ", databaseConfig)

	home.NewHomeHandler(app)

	app.Listen(":3000")
}
