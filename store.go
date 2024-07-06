package auth

import (
	"database/sql"
	"errors"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

type Store struct {
	DB     *sql.DB
	config Config
}

func NewStore(config Config) (*Store, error) {
	db, err := sql.Open("sqlite", config.DatabaseFilePath)
	if err != nil {
		return nil, err
	}
	s := &Store{
		DB:     db,
		config: config,
	}

	err = s.initializeTables()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) initializeTables() error {
	// Check if the users table exists
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='users';`
	var tableName string
	err := s.DB.QueryRow(query).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// If the users table doesn't exist, initialize the database schema
	if tableName == "" {
		bytes, err := os.ReadFile("./migrations/initial_schema.sql")
		if err != nil {
			log.Fatalf("Failed to read migrations: %v", err)
			return err
		}
		_, err = s.DB.Exec(string(bytes))
		if err != nil {
			log.Fatalf("Failed to create tables: %v", err)
			return err
		}
	}

	return nil
}

func (s *Store) CreateUser(username, password string) (int64, error) {
	stmt, err := s.DB.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(username, password)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	row := s.DB.QueryRow("SELECT id, username, password, created_at FROM users WHERE username = ?", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByID(id int) (*User, error) {
	user := &User{}
	row := s.DB.QueryRow("SELECT id, username, password, created_at FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) CreateSession(userID int, token string) (int64, error) {
	stmt, err := s.DB.Prepare("INSERT INTO sessions (user_id, token) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(userID, token)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetSessionByToken(token string) (*Session, error) {
	session := &Session{}
	row := s.DB.QueryRow("SELECT id, user_id, token, created_at FROM sessions WHERE token = ?", token)
	err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Store) DeleteSessionByToken(token string) error {
	stmt, err := s.DB.Prepare("DELETE FROM sessions WHERE token = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(token)
	return err
}

func (s *Store) SavePasswordResetToken(userID int, token string) error {
	stmt, err := s.DB.Prepare("INSERT INTO password_resets (user_id, token) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, token)
	return err
}

func (s *Store) GetUserIDByPasswordResetToken(token string) (int, error) {
	var userID int
	row := s.DB.QueryRow("SELECT user_id FROM password_resets WHERE token = ?", token)
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *Store) UpdateUserPassword(userID int, password string) error {
	stmt, err := s.DB.Prepare("UPDATE users SET password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(password, userID)
	return err
}
