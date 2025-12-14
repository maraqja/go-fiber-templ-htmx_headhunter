package vacancy

import (
	"time"

	"github.com/a-h/templ"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofiber/fiber/v2"
	templadapter "github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/templ_adapter"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/validator"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/views/components"
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
	form := VacancyCreateForm{
		Email: c.FormValue("email"),
	}
	errors := validate.Validate(
		&validators.EmailIsPresent{Name: "email", Field: form.Email},
	)
	time.Sleep(2 * time.Second)
	var component templ.Component
	if len(errors.Errors) > 0 {
		component = components.Notification(validator.FormatErrors(errors), components.NotificationStatusError)
		return templadapter.Render(c, component)
	}
	component = components.Notification("Vacancy created", components.NotificationStatusSuccess)
	return templadapter.Render(c, component)
}
