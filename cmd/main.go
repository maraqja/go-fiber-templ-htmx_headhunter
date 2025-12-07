package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/config"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/internal/home"
)

func main() {
	app := fiber.New()
	app.Use(recover.New())

	config.LoadEnvFile()

	databaseConfig := config.NewDatabaseConfig()
	log.Println("Database: ", databaseConfig)

	home.NewHomeHandler(app)

	app.Listen(":3000")
}
