package main

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	postgresStorage "github.com/gofiber/storage/postgres/v3"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/config"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/internal/home"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/internal/sitemap"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/internal/vacancy"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/logger"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/middleware"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/postgres"
	"github.com/rs/zerolog/log"
)

func main() {

	// Инициализация логгера по умолчанию для логирования процесса загрузки (до получения конфига для логгера)
	if err := logger.InitWithDefaults(); err != nil {
		panic(err)
	}

	config.LoadEnvFile()

	logConfig, err := config.NewLogConfig()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get log config")
		panic(err)
	}

	logOutputType, err := logger.ParseOutputType(logConfig.Output)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Invalid output type, using stdout")
		panic(err)
	}

	logFormatType, err := logger.ParseFormatType(logConfig.Format)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Invalid format type, using json")
		panic(err)
	}

	// Инициализация логгера по конфигу, который далее будет использоваться во всем приложении
	if err := logger.Init(logger.Config{
		Level:  logConfig.Level,
		Output: logOutputType,
		Format: logFormatType,
	}); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to initialize logger")
		panic(err)
	}

	app := fiber.New()
	app.Static("/static", "./static")
	app.Static("/robots.txt", "./static/robots.txt")

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))
	app.Use(recover.New())

	_ = config.NewDatabaseConfig()

	databaseConfig := config.NewDatabaseConfig()
	pgpool, err := postgres.NewPool(&postgres.Config{
		URL: databaseConfig.Url,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to create database pool")
		panic(err)
	}
	defer pgpool.Close()

	store := session.New(session.Config{
		Storage: postgresStorage.New(postgresStorage.Config{
			DB:    pgpool,
			Table: "sessions",
		}),
	})
	app.Use(middleware.AuthMiddleware(store))
	// -- Repositories --
	vacancyRepo := vacancy.NewPostgresRepository(vacancy.RepositoryDI{
		DB: pgpool,
	})

	// -- Handlers --
	home.NewHomeHandler(home.HandlerDI{
		Router:     app,
		Repository: vacancyRepo,
		Store:      store,
	})
	vacancy.NewHandler(vacancy.HandlerDI{
		Router:     app,
		Repository: vacancyRepo,
	})

	sitemap.NewHandler(app)
	app.Listen(":3000")
}
