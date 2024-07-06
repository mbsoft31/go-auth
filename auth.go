package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
)

func hashPassword(password string) string {
	// Use bcrypt or another strong hashing algorithm instead of SHA-256
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

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
	hashedPassword := hashPassword(password)
	return s.CreateUser(username, hashedPassword)
}

func Login(s *Store, username, password string, w http.ResponseWriter) error {
	user, err := s.GetUserByUsername(username)
	if err != nil || hashPassword(password) != user.Password {
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
	hashedPassword := hashPassword(newPassword)
	return s.UpdateUserPassword(userID, hashedPassword)
}
