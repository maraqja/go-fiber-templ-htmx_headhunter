package home

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type HomeHandler struct {
	router fiber.Router
}

func NewHomeHandler(router fiber.Router) *HomeHandler {
	h := &HomeHandler{router: router}
	apiRouterGroup := h.router.Group("/api") // Указываем префикс /api для всех маршрутов в этой группе
	apiRouterGroup.Get("/", h.home)
	return h
}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	// panic("failed to get home")
	// return fiber.NewError(fiber.StatusBadRequest, "failed to get home")

	log.Info().
		Bool("is_home", true).
		Int("status_code", fiber.StatusOK).
		Str("ip", c.IP()).
		Str("user_agent", c.Get("User-Agent")).
		Msg("GET /api/ - get home") // в zerolog используется паттерн билдера для построения сообщения
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}
