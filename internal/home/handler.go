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

const (
	DefaultLimit = 2
)

type HandlerDI struct {
	Router     fiber.Router
	Repository IRepository
}

type IRepository interface {
	GetVacancies(ctx context.Context, limit int, offset int) ([]vacancy.Vacancy, error)
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
	limit := c.QueryInt("limit", DefaultLimit)
	page := c.QueryInt("page", 1)
	offset := (page - 1) * limit
	vacancies, err := h.repository.GetVacancies(c.Context(), limit, offset)
	if err != nil {
		component := components.Notification(err.Error(), components.NotificationStatusError)
		return templadapter.Render(c, component, http.StatusInternalServerError)
	}
	component := views.Main(vacancies)
	return templadapter.Render(c, component, http.StatusOK)
}
