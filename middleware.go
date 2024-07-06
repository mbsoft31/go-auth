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
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized", "message": err.Error()})
			}

			sessionToken := cookie.Value
			session, err := store.GetSessionByToken(sessionToken)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized", "message": err.Error()})
			}
			user, err := store.GetUserByID(session.UserID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized", "message": err.Error()})
			}
			c.Set("user", user)
			return next(c)
		}
	}
}
