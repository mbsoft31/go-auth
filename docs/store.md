# Store

The `Store` struct is responsible for interacting with the SQLite database.

## Initialization

### `NewStore`
Initializes the store and creates required tables.
```go
func NewStore(config Config) (*Store, error)
```

## Tables
- `users`: Stores user information
- `sessions`: Manages session tokens
- `password_resets`: Handles password reset tokens

## Common Methods

### **initializeTables**
Creates tables if they do not exist.

### **CreateUser**
Inserts a new user into the `users` table.

### **GetSessionByToken**
Retrieves session details using the token.

[Back to README](./index.md)