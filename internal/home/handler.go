package home

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/internal/vacancy"
	templadapter "github.com/maraqja/go-fiber-templ-htmx_headhunter/pkg/templ_adapter"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/views"
	"github.com/maraqja/go-fiber-templ-htmx_headhunter/views/components"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	DefaultLimit = 2
	UserIDKey    = "user_id"
)

type HandlerDI struct {
	Router     fiber.Router
	Repository IRepository
	Store      *session.Store
}

type IRepository interface {
	GetVacancies(ctx context.Context, limit int, offset int) ([]vacancy.Vacancy, error)
	GetVacanciesCount(ctx context.Context) (int, error)
}

type HomeHandler struct {
	router     fiber.Router
	repository IRepository
	store      *session.Store
	logger     *zerolog.Logger
}

func NewHomeHandler(di HandlerDI) *HomeHandler {
	logger := log.Logger.With().Str("component", "HomeHandler").Logger()
	h := &HomeHandler{router: di.Router, repository: di.Repository, store: di.Store, logger: &logger}
	h.router.Get("/", h.home)
	h.router.Get("/login", h.login)
	h.router.Post("/api/login", h.apiLogin)
	h.router.Post("/api/logout", h.apiLogout)
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

	// Получаем сессию: читает session_id из Cookie, загружает данные из storage
	session, err := h.store.Get(c)
	if err != nil {
		return c.Redirect("/login", http.StatusUnauthorized)
	}
	// Проверяем наличие user_id в сессии
	userID, ok := session.Get(UserIDKey).(string)
	if !ok || userID == "" {
		h.logger.Info().Msg("User ID not found in session")
		// return c.Redirect("/login", http.StatusUnauthorized)
	} else {
		h.logger.Info().Str("user_id", userID).Msg("User authenticated")
	}
	var email string = ""
	email, ok = session.Get("email").(string)
	c.Locals("email", email)

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

	// Получаем сессию: читает session_id из Cookie или создаёт новую
	session, err := h.store.Get(c)
	if err != nil {
		return c.Redirect("/login", http.StatusInternalServerError)
	}
	// Сохраняем user_id в памяти объекта сессии (ещё не отправлено клиенту)
	session.Set(UserIDKey, "1") // сохраняем как string для единообразия
	// Сохраняем в storage и отправляем Set-Cookie браузеру
	if err := session.Save(); err != nil {
		return c.Redirect("/login", http.StatusInternalServerError)
	}
	userEmail := ""
	if email, ok := session.Get("email").(string); ok {
		userEmail = email
	}
	c.Locals("email", userEmail)
	return templadapter.Render(c, component, http.StatusOK)
}

func (h *HomeHandler) apiLogin(c *fiber.Ctx) error { // Для мокового теста логина
	form := LoginForm{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}
	if form.Email == "a@a.ru" && form.Password == "1" { // Мок вместо проверки в БД
		sess, err := h.store.Get(c)
		if err != nil {
			panic(err)
		}
		sess.Set("email", form.Email)
		if err := sess.Save(); err != nil {
			panic(err)
		}
		c.Response().Header.Add("Hx-Redirect", "/")
		return c.Redirect("/", http.StatusOK)
	}
	component := components.Notification("Неверный логин или пароль", components.NotificationStatusError)
	return templadapter.Render(c, component, http.StatusBadRequest)
}

func (h *HomeHandler) apiLogout(c *fiber.Ctx) error {
	session, err := h.store.Get(c)
	if err != nil {
		c.Response().Header.Add("Hx-Redirect", "/login")
		return c.Redirect("/login", http.StatusInternalServerError)
	}
	// Destroy() удаляет всю сессию (все данные + из storage)
	if err := session.Destroy(); err != nil {
		c.Response().Header.Add("Hx-Redirect", "/login")
		return c.Redirect("/login", http.StatusInternalServerError)
	}
	// Save() не нужен после Destroy(), так как Destroy() сам сохраняет изменения
	c.Response().Header.Add("Hx-Redirect", "/")
	return c.Redirect("/", http.StatusOK)
}
