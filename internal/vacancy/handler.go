package vacancy

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	templadapter "github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/templ_adapter"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/views/components"
	"github.com/rs/zerolog/log"
)

type VacancyHandler struct {
	router fiber.Router
}

func NewVacancyHandler(router fiber.Router) *VacancyHandler {
	h := &VacancyHandler{router: router}
	vacancyGroup := h.router.Group("/vacancy")
	vacancyGroup.Post("/", h.createVacancy)
	return h
}

func (h *VacancyHandler) createVacancy(c *fiber.Ctx) error {
	email := c.FormValue("email")
	log.Logger.Info().Str("email", email).Msg("Creating vacancy")
	var component templ.Component
	if email == "" {
		component = components.Notification("Email is required", components.NotificationStatusError)
		return templadapter.Render(c, component)
	}
	component = components.Notification("Vacancy created", components.NotificationStatusSuccess)
	return templadapter.Render(c, component) // Возвращаем html, который будет отображен в div с id="vacancy-result" с помощью hx-swap="innerHTML"
}
