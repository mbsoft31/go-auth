package auth

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func AuthMiddleware(store *Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(store.config.CookieName)
			if err != nil {
				log.Printf("Unauthorized access attempt: %v", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "Unauthorized",
					"message": "Invalid session. Please log in again.",
				})
			}

			sessionToken := cookie.Value
			session, err := store.GetSessionByToken(sessionToken)
			if err != nil {
				log.Printf("Unauthorized access attempt: %v", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "Unauthorized",
					"message": "Invalid session. Please log in again.",
				})
			}
			user, err := store.GetUserByID(session.UserID)
			if err != nil {
				log.Printf("Unauthorized access attempt: %v", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "Unauthorized",
					"message": "Invalid session. Please log in again.",
				})
			}
			c.Set("user", user)
			return next(c)
		}
	}
}
