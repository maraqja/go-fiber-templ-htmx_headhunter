package main

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/config"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/internal/home"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/logger"
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

	// Инициализация HTML шаблонизатора: указываем директорию с шаблонами и их расширение
	// html.New принимает путь к папке с шаблонами и расширение файлов шаблонов
	htmlEngine := html.New("./html", ".html")

	// Создание Fiber приложения с конфигурацией шаблонизатора
	// Views позволяет использовать метод c.Render() в обработчиках для рендеринга HTML шаблонов
	app := fiber.New(fiber.Config{
		Views: htmlEngine,
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))
	app.Use(recover.New())
	_ = config.NewDatabaseConfig()

	home.NewHomeHandler(app)

	app.Listen(":3000")
}
