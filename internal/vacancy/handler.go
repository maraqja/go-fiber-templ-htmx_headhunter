package vacancy

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	fmt.Println(email)
	log.Logger.Info().Str("email", email)
	return c.SendString("Vacancy created")
}
