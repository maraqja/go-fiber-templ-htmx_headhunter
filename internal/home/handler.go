package home

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	templadapter "github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/templ_adapter"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/views"
)

type HomeHandler struct {
	router fiber.Router
}

func NewHomeHandler(router fiber.Router) *HomeHandler {
	h := &HomeHandler{router: router}
	h.router.Get("/", h.home)
	return h
}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	component := views.Main()
	return templadapter.Render(c, component, http.StatusOK)
}
