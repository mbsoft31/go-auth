# Config

The `Config` struct holds configuration settings for cookies, sessions, and database storage.

## Default Configuration

Use `DefaultConfig` for recommended settings:
```go
func DefaultConfig() Config
```

### Fields
- `CookieName`: Name of the session cookie
- `CookieMaxAge`: Cookie expiration time in seconds
- `SessionTokenTTL`: Session validity duration
- `DatabaseFilePath`: Path to the SQLite database

### Example
```go
config := DefaultConfig()
config.CookieSecure = true // Use secure cookies in production
```

## Custom Configuration

Override defaults as needed:
```go
config := DefaultConfig()
config.CookieName = "my_session_cookie"
```

[Back to README](./index.md)