## Go Authentication Package

### Overview

This package provides a simple and reusable authentication system for Go web applications. It uses SQLite to store user information and sessions, similar to Laravel's auth system. The package supports registration, login, password reset, and session management functionalities.

### Features

- User registration and login
- Session management using cookies
- Password hashing and verification
- Password reset functionality
- Middleware for user authentication
- Configurable settings for customization

### Installation

To install the package, use the following command:

```bash
go get github.com/mbsoft31/go-auth
```

### Usage

#### Configuration

You can customize the authentication package by providing a configuration when initializing the store. If no configuration is provided, default settings will be used.

```go
package main

import (
    "log"
    "github.com/mbsoft31/go-auth"
)

func main() {
    config := auth.Config{
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

    store, err := auth.NewStore(config)
    if err != nil {
        log.Fatalf("Failed to create store: %v", err)
    }

    // Now you can use the store for user registration, login, etc.
}
```

#### Example with Echo Framework

Here's an example of how to use the authentication package with the Echo framework:

```go
package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/mbsoft31/go-auth"
)

func main() {
    config := auth.DefaultConfig()
    store, err := auth.NewStore(config)
    if err != nil {
        panic(err)
    }
    
    e := echo.New()
    
    // Register route
    e.POST("/register", func(c echo.Context) error {
        username := c.FormValue("username")
        password := c.FormValue("password")
        _, err := auth.Register(store, username, password)
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "Registration failed"})
        }
        return c.JSON(http.StatusOK, map[string]string{"message": "Registration successful"})
    })
    
    // Login route
    e.POST("/login", func(c echo.Context) error {
        username := c.FormValue("username")
        password := c.FormValue("password")
        err := auth.Login(store, username, password, c.Response())
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
        }
        return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
    })
    
    // Protected route
    e.GET("/protected", func(c echo.Context) error {
        user := c.Get("user").(*auth.User)
        return c.JSON(http.StatusOK, map[string]string{"message": "Hello " + user.Username})
    }, auth.AuthMiddleware(store))
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

#### Initialization

Create a new store with the provided configuration:

```go
store, err := auth.NewStore(auth.DefaultConfig())
if err != nil {
    log.Fatalf("Failed to create store: %v", err)
}
```

#### User Registration

Register a new user with a username and password:

```go
userID, err := auth.Register(store, "username", "password")
if err != nil {
    log.Fatalf("Failed to register user: %v", err)
}
```

#### User Login

Authenticate a user and create a session token:

```go
err := auth.Login(store, "username", "password", responseWriter)
if err != nil {
    log.Fatalf("Failed to login user: %v", err)
}
```

#### User Logout

Log out a user and clear the session token cookie:

```go
err := auth.Logout(store, responseWriter, "session_token")
if err != nil {
    log.Fatalf("Failed to logout user: %v", err)
}
```

#### Password Reset

Initiate a password reset:

```go
resetToken, err := auth.InitiatePasswordReset(store, "username")
if err != nil {
    log.Fatalf("Failed to initiate password reset: %v", err)
}
```

Reset the user's password using the reset token:

```go
err := auth.ResetPassword(store, "reset_token", "new_password")
if err != nil {
    log.Fatalf("Failed to reset password: %v", err)
}
```

### Configuration Options

The `Config` struct holds all the configuration settings for the auth package:

- `CookieName`: Name of the session token cookie.
- `CookiePath`: Path for the session token cookie.
- `CookieMaxAge`: Max age for the session token cookie.
- `CookieHttpOnly`: HttpOnly attribute for the session token cookie.
- `CookieSecure`: Secure attribute for the session token cookie.
- `CookieSameSite`: SameSite attribute for the session token cookie.
- `SessionTokenTTL`: Time-to-live for session tokens.
- `PasswordResetTTL`: Time-to-live for password reset tokens.
- `HashCost`: Cost for password hashing.
- `DatabaseFilePath`: Path to the SQLite database file.

### Testing

To run tests for the package, use the following command:

```bash
go test ./...
```

### Migrations

The `initial_schema.sql` file should contain the SQL statements to create the necessary tables:

```sql
-- initial_schema.sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS password_resets (
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
```

### License

This package is licensed under the MIT License.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

### Author

This package was created by Mouadh Bekhouche.

---# go-auth
