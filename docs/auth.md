
# Auth

The `auth` package manages user authentication and session handling.

## Key Functions

### **CreateUser**
Creates a new user in the database.
```go
func (s *Store) CreateUser(username, password string) (int64, error)
```

### **GetUserByUsername**
Retrieves a user by their username.
```go
func (s *Store) GetUserByUsername(username string) (*User, error)
```

### **CreateSession**
Creates a session for a user.
```go
func (s *Store) CreateSession(userID int, token string) (int64, error)
```

### **SavePasswordResetToken**
Stores a password reset token for a user.
```go
func (s *Store) SavePasswordResetToken(userID int, token string) error
```

### **UpdateUserPassword**
Updates a user's password.
```go
func (s *Store) UpdateUserPassword(userID int, password string) error
```

[Back to README](./index.md)