package templadapter

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Render(c *fiber.Ctx, component templ.Component, code int) error {
	// return component.Render(c.Context(), c.Response().BodyWriter()) // не установит заголовки, не обработает ошибки рендеринга итд
	return adaptor.HTTPHandler(templ.Handler(component, templ.WithStatus(code)))(c)
}
