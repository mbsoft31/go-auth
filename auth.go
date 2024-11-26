package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// HashFunc defines the signature for a hashing function.
// It accepts a byte slice as input and returns a hashed byte slice.
type HashFunc func(data []byte) []byte

// HashPasswordParams holds the parameters required to hash a password.
type HashPasswordParams struct {
	Password string   // The password to be hashed.
	HashFunc HashFunc // The hashing function to use. If nil, a default function is used.
}

func BcryptHash(data []byte) []byte {
	hash, _ := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)
	return hash
}

func Sha256Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// DefaultHashFunc provides a secure default hashing implementation using SHA-256.
var DefaultHashFunc HashFunc = Sha256Hash

// hashPassword hashes a password using the provided parameters.
// If no hashing function is provided, it uses the default SHA-256 implementation.
func hashPassword(params HashPasswordParams) (string, error) {
	// Validate the input
	if params.Password == "" {
		return "", ErrEmptyPassword
	}

	// Use the default hashing function if none is provided
	if params.HashFunc == nil {
		params.HashFunc = DefaultHashFunc
	}

	// Hash the password
	hashedData := params.HashFunc([]byte(params.Password))

	// Return the hexadecimal string representation of the hash
	return hex.EncodeToString(hashedData), nil
}

// ErrEmptyPassword is returned when the provided password is empty.
var ErrEmptyPassword = errors.New("password cannot be empty")

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func generatePasswordResetToken() (string, error) {
	// Generate a unique token for password reset
	return generateToken()
}

func Register(s *Store, username, password string) (int64, error) {
	hashedPassword, err := hashPassword(HashPasswordParams{Password: password})
	if err != nil {
		return 0, err
	}
	return s.CreateUser(username, hashedPassword)
}

func Login(s *Store, username, password string, w http.ResponseWriter) error {
	hashedPassword, err := hashPassword(HashPasswordParams{Password: password})
	if err != nil {
		return err
	}
	user, err := s.GetUserByUsername(username)
	if err != nil || hashedPassword != user.Password {
		return errors.New("invalid username or password")
	}
	token, err := generateToken()
	if err != nil {
		return err
	}
	_, err = s.CreateSession(user.ID, token)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Use Secure attribute in production
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

func Logout(s *Store, w http.ResponseWriter, token string) error {
	// Delete session from the database
	err := s.DeleteSessionByToken(token)
	if err != nil {
		return err
	}

	// Clear the session token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true, // Use Secure attribute in production
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func Authenticate(s *Store, token string) (*User, error) {
	session, err := s.GetSessionByToken(token)
	if err != nil {
		return nil, errors.New("invalid session")
	}
	user, err := s.GetUserByID(session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func InitiatePasswordReset(s *Store, username string) (string, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	resetToken, err := generatePasswordResetToken()
	if err != nil {
		return "", err
	}
	err = s.SavePasswordResetToken(user.ID, resetToken)
	if err != nil {
		return "", err
	}
	return resetToken, nil
}

func ResetPassword(s *Store, resetToken, newPassword string) error {
	userID, err := s.GetUserIDByPasswordResetToken(resetToken)
	if err != nil {
		return err
	}
	hashedPassword, err := hashPassword(HashPasswordParams{Password: newPassword})
	if err != nil {
		return err
	}
	return s.UpdateUserPassword(userID, hashedPassword)
}
