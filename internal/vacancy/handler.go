package vacancy

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofiber/fiber/v2"
	templadapter "github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/templ_adapter"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/validator"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/views/components"
)

type IRepository interface {
	CreateVacancy(ctx context.Context, form VacancyCreateForm) error
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
	var component templ.Component
	if len(errors.Errors) > 0 {
		component = components.Notification(validator.FormatErrors(errors), components.NotificationStatusError)
		return templadapter.Render(c, component, http.StatusBadRequest)
	}
	err := h.repository.CreateVacancy(c.Context(), form)
	if err != nil {
		component = components.Notification(err.Error(), components.NotificationStatusError)
		return templadapter.Render(c, component, http.StatusInternalServerError)
	}
	component = components.Notification("Vacancy created", components.NotificationStatusSuccess)
	return templadapter.Render(c, component, http.StatusOK)
}
