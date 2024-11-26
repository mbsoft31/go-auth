# Middleware

Middleware functions for authentication and session validation using the Echo framework.

## AuthMiddleware
Validates a session token from the cookie and adds the user to the context.

### Usage
```go
e.Use(AuthMiddleware(store))
```

### Function Signature
```go
func AuthMiddleware(store *Store) echo.MiddlewareFunc
```

### Example
```go
e.GET("/profile", func(c echo.Context) error {
    user := c.Get("user").(*User)
    return c.JSON(http.StatusOK, user)
})
```

[Back to README](./index.md)
