package home

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/internal/vacancy"
	templadapter "github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/templ_adapter"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/views"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/views/components"
)

type HandlerDI struct {
	Router     fiber.Router
	Repository IRepository
}

type IRepository interface {
	GetVacancies(ctx context.Context) ([]vacancy.Vacancy, error)
}

type HomeHandler struct {
	router     fiber.Router
	repository IRepository
}

func NewHomeHandler(di HandlerDI) *HomeHandler {
	h := &HomeHandler{router: di.Router, repository: di.Repository}
	h.router.Get("/", h.home)
	return h
}

func (h *HomeHandler) home(c *fiber.Ctx) error {

	vacancies, err := h.repository.GetVacancies(c.Context())
	if err != nil {
		component := components.Notification(err.Error(), components.NotificationStatusError)
		return templadapter.Render(c, component, http.StatusInternalServerError)
	}
	// component := components.VacancyCards(vacancies)
	// return templadapter.Render(c, component, http.StatusOK)
	component := views.Main(vacancies)
	return templadapter.Render(c, component, http.StatusOK)
}
