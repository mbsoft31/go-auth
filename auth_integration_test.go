package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDB(t *testing.T) *Store {
	config := DefaultConfig()
	config.DatabaseFilePath = ":memory:"
	store, err := NewStore(config)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	return store
}

func TestAuthIntegration(t *testing.T) {
	e := echo.New()
	store := setupTestDB(t)

	// Register route
	e.POST("/register", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		_, err := Register(store, username, password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Registration failed"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Registration successful"})
	})

	// Login route
	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		err := Login(store, username, password, c.Response())
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
	})

	// Protected route
	e.GET("/protected", func(c echo.Context) error {
		user := c.Get("user").(*User)
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello " + user.Username})
	}, AuthMiddleware(store))

	// Register user
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("username=testuser&password=password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"Registration successful"}`, rec.Body.String())

	// Login user
	req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=testuser&password=password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"Login successful"}`, rec.Body.String())

	// Get the session cookie
	cookie := rec.Result().Cookies()[0]
	
	// Access protected route
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"Hello testuser"}`, rec.Body.String())
}
