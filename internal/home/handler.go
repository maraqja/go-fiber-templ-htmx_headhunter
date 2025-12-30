package home

import (
	"context"
	"fmt"
	"math"
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
	GetVacanciesCount(ctx context.Context) (int, error)
}

type HomeHandler struct {
	router     fiber.Router
	repository IRepository
}

func NewHomeHandler(di HandlerDI) *HomeHandler {
	h := &HomeHandler{router: di.Router, repository: di.Repository}
	h.router.Get("/", h.home)
	h.router.Get("/login", h.login)
	return h
}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", DefaultLimit)
	page := c.QueryInt("page", 1)

	if limit < 1 || page < 1 {
		return c.Redirect(fmt.Sprintf("/?page=1&limit=%d", DefaultLimit), http.StatusFound)
	}
	offset := (page - 1) * limit
	count, err := h.repository.GetVacanciesCount(c.Context())
	if err != nil {
		component := components.Notification(err.Error(), components.NotificationStatusError)
		return templadapter.Render(c, component, http.StatusInternalServerError)
	}
	vacancies, err := h.repository.GetVacancies(c.Context(), limit, offset)
	if err != nil {
		component := components.Notification(err.Error(), components.NotificationStatusError)
		return templadapter.Render(c, component, http.StatusInternalServerError)
	}
	component := views.Main(vacancies, int(math.Ceil(float64(count/limit))), page)
	return templadapter.Render(c, component, http.StatusOK)
}

func (h *HomeHandler) login(c *fiber.Ctx) error {
	component := views.Login()
	return templadapter.Render(c, component, http.StatusOK)
}
