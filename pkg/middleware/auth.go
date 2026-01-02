package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const (
	UserIDKey = "user_id"
	EmailKey  = "email"
)

func AuthMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			c.Response().Header.Add("Hx-Redirect", "/login")
			return c.Redirect("/login", http.StatusUnauthorized)
		}

		var userID string = ""
		var email string = ""
		if userIDFromSession, ok := session.Get(UserIDKey).(string); ok {
			userID = userIDFromSession
		}
		c.Locals(UserIDKey, userID)
		if emailFromSession, ok := session.Get(EmailKey).(string); ok {
			email = emailFromSession
		}
		c.Locals(EmailKey, email)
		return c.Next()
	}
}
