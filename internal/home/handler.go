package home

import (
	"github.com/gofiber/fiber/v2"
)

type HomeHandler struct {
	router fiber.Router
}

func NewHomeHandler(router fiber.Router) *HomeHandler {
	h := &HomeHandler{router: router}
	apiRouterGroup := h.router.Group("/api") // Указываем префикс /api для всех маршрутов в этой группе
	apiRouterGroup.Get("/", h.home)
	return h
}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	// log.Info().
	// 	Bool("is_home", true).
	// 	Int("status_code", fiber.StatusOK).
	// 	Str("ip", c.IP()).
	// 	Str("user_agent", c.Get("User-Agent")).
	// 	Msg("GET /api/ - get home") // в zerolog используется паттерн билдера для построения сообщения
	// ------------------------------------------------------------------------------------------------
	// tmpl := template.Must(template.ParseFiles("./html/page.html"))
	// var tmp bytes.Buffer
	// if err := tmpl.Execute(&tmp, fiber.Map{
	// 	"Message": "Hello, World!",
	// }); err != nil {
	// 	return fiber.NewError(fiber.StatusInternalServerError, "failed to execute template")
	// }
	// c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	// return c.Send(tmp.Bytes())
	names := []string{
		"John",
		"Jane",
		"Jim",
		"Jill",
	}

	users := []struct {
		Name  string
		Age   int
		Email string
	}{
		{Name: "John", Age: 30, Email: "john@example.com"},
		{Name: "Jane", Age: 25, Email: "jane@example.com"},
		{Name: "Jim", Age: 35, Email: "jim@example.com"},
		{Name: "Jill", Age: 40, Email: "jill@example.com"},
	}

	return c.Render("page", fiber.Map{
		"Names": names,
		"Users": users,
	})

}
