package auth_test

import (
	"github.com/labstack/echo/v4"
	"github.com/mbsoft31/go-auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDB(t *testing.T) *auth.Store {
	config := auth.DefaultConfig()
	config.DatabaseFilePath = ":memory:"
	store, err := auth.NewStore(config)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	return store
}
func TestRegister(t *testing.T) {
	store := setupTestDB(t)

	userID, err := auth.Register(store, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if userID == 0 {
		t.Fatalf("Expected valid user ID, got %d", userID)
	}
}

func TestLogin(t *testing.T) {
	store := setupTestDB(t)

	_, err := auth.Register(store, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=testuser&password=password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	err = auth.Login(store, "testuser", "password", rec)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	cookie := rec.Result().Cookies()
	if len(cookie) == 0 || cookie[0].Name != "session_token" {
		t.Fatalf("Expected session token cookie, got %v", cookie)
	}
}

func TestAuthenticate(t *testing.T) {
	store := setupTestDB(t)

	_, err := auth.Register(store, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=testuser&password=password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	err = auth.Login(store, "testuser", "password", rec)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	cookie := rec.Result().Cookies()[0]

	user, err := auth.Authenticate(store, cookie.Value)
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	if user.Username != "testuser" {
		t.Fatalf("Expected username 'testuser', got %v", user.Username)
	}
}

func TestAuthMiddleware(t *testing.T) {
	store := setupTestDB(t)

	// Register and login the user
	_, err := auth.Register(store, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=testuser&password=password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	err = auth.Login(store, "testuser", "password", rec)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	cookie := rec.Result().Cookies()[0]

	// Create a request with the session token cookie
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()

	// Initialize Echo
	e := echo.New()

	// Define the handler
	handler := func(c echo.Context) error {
		user := c.Get("user").(*auth.User)
		if user.Username != "testuser" {
			t.Fatalf("Expected username 'testuser', got %v", user.Username)
		}
		return c.NoContent(http.StatusOK)
	}

	// Create the route with the middleware
	e.GET("/protected", handler, auth.AuthMiddleware(store))

	// Serve the request
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", rec.Code)
	}
}

func TestInitiatePasswordReset(t *testing.T) {
	store := setupTestDB(t)

	_, err := auth.Register(store, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	resetToken, err := auth.InitiatePasswordReset(store, "testuser")
	if err != nil {
		t.Fatalf("Failed to initiate password reset: %v", err)
	}

	if resetToken == "" {
		t.Fatalf("Expected valid reset token, got empty string")
	}
}

func TestResetPassword(t *testing.T) {
	store := setupTestDB(t)

	_, err := auth.Register(store, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	resetToken, err := auth.InitiatePasswordReset(store, "testuser")
	if err != nil {
		t.Fatalf("Failed to initiate password reset: %v", err)
	}

	err = auth.ResetPassword(store, resetToken, "newpassword")
	if err != nil {
		t.Fatalf("Failed to reset password: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=testuser&password=newpassword"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	err = auth.Login(store, "testuser", "newpassword", rec)
	if err != nil {
		t.Fatalf("Failed to login with new password: %v", err)
	}
}

func TestLogout(t *testing.T) {
	store := setupTestDB(t)

	_, err := auth.Register(store, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=testuser&password=password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	err = auth.Login(store, "testuser", "password", rec)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	cookie := rec.Result().Cookies()[0]

	rec = httptest.NewRecorder()
	err = auth.Logout(store, rec, cookie.Value)
	if err != nil {
		t.Fatalf("Failed to logout: %v", err)
	}

	if len(rec.Result().Cookies()) == 0 || rec.Result().Cookies()[0].Value != "" {
		t.Fatalf("Expected empty session token cookie, got %v", rec.Result().Cookies())
	}
}
