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

type IRepository interface {
	CreateVacancy(form VacancyCreateForm)
}

type HandlerDI struct {
	Router     fiber.Router
	Repository IRepository
}

type Handler struct {
	router     fiber.Router
	repository IRepository
}

func NewHandler(di HandlerDI) *Handler {
	h := &Handler{
		router:     di.Router,
		repository: di.Repository,
	}
	vacancyGroup := h.router.Group("/vacancy")
	vacancyGroup.Post("/", h.createVacancy)
	return h
}

func (h *Handler) createVacancy(c *fiber.Ctx) error {
	form := VacancyCreateForm{
		Email:    c.FormValue("email"),
		Role:     c.FormValue("role"),
		Company:  c.FormValue("company"),
		Salary:   c.FormValue("salary"),
		Type:     c.FormValue("type"),
		Location: c.FormValue("location"),
	}
	errors := validate.Validate(
		&validators.EmailIsPresent{Name: "email", Field: form.Email},
		&validators.StringIsPresent{Name: "role", Field: form.Role},
		&validators.StringIsPresent{Name: "company", Field: form.Company},
		&validators.StringIsPresent{Name: "salary", Field: form.Salary},
		&validators.StringIsPresent{Name: "type", Field: form.Type},
		&validators.StringIsPresent{Name: "location", Field: form.Location},
	)
	time.Sleep(1 * time.Second) // для дебага лоадера
	var component templ.Component
	if len(errors.Errors) > 0 {
		component = components.Notification(validator.FormatErrors(errors), components.NotificationStatusError)
		return templadapter.Render(c, component)
	}
	component = components.Notification("Vacancy created", components.NotificationStatusSuccess)
	return templadapter.Render(c, component)
}
