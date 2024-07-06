package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthMiddleware(store *Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(store.config.CookieName)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			sessionToken := cookie.Value
			session, err := store.GetSessionByToken(sessionToken)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}
			user, err := store.GetUserByID(session.UserID)
			if err != nil {
				return err
			}
			c.Set("user", user)
			return next(c)
		}
	}
}
