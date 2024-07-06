package auth

import (
	"time"
)

// Config holds the configuration settings for the auth package
type Config struct {
	CookieName       string
	CookiePath       string
	CookieMaxAge     int
	CookieHttpOnly   bool
	CookieSecure     bool
	CookieSameSite   int
	SessionTokenTTL  time.Duration
	PasswordResetTTL time.Duration
	HashCost         int
	DatabaseFilePath string
}

// DefaultConfig returns a Config with default settings
func DefaultConfig() Config {
	return Config{
		CookieName:       "session_token",
		CookiePath:       "/",
		CookieMaxAge:     3600,
		CookieHttpOnly:   true,
		CookieSecure:     false,
		CookieSameSite:   2, // http.SameSiteStrictMode
		SessionTokenTTL:  24 * time.Hour,
		PasswordResetTTL: 1 * time.Hour,
		HashCost:         14,
		DatabaseFilePath: "auth.db",
	}
}
