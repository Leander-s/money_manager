package database

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Token  uuid.UUID `json:"token"`
	UserID uuid.UUID `json:"userID"`
	Expiry time.Time `json:"expiry"`
}

func (db *Database) InsertToken(Token *Token) error {
	err := db.DB.QueryRow(
		"INSERT INTO tokens (token, user_id, expires_at) VALUES ($1, $2, $3)",
		Token.Token, Token.UserID, Token.Expiry,
	).Err()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetToken(token *uuid.UUID) (*Token, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	}
	Token := &Token{}
	err := db.DB.QueryRow(
		"SELECT user_id, expires_at FROM tokens WHERE token = $1",
		token,
	).Scan(&Token.UserID, &Token.Expiry)
	Token.Token = *token
	return Token, err
}

func (db *Database) ListTokens() ([]*Token, error) {
	rows, err := db.DB.Query("SELECT token, user_id, expires_at FROM tokens")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*Token
	for rows.Next() {
		token := &Token{}
		if err := rows.Scan(&token.Token, &token.UserID, &token.Expiry); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, rows.Err()
}

func (db *Database) DeleteToken(token *uuid.UUID) error {
	if token == nil {
		return errors.New("token is nil")
	}
	_, err := db.DB.Exec(
		"DELETE FROM tokens WHERE token = $1",
		token,
	)
	return err
}

func (db *Database) DeleteTokensByUserID(userID *uuid.UUID) error {
	if userID == nil {
		return errors.New("userID is nil")
	}
	_, err := db.DB.Exec(
		"DELETE FROM tokens WHERE user_id = $1",
		userID,
	)
	return err
}

func (db *Database) DeleteExpiredTokens() error {
	_, err := db.DB.Exec(
		"DELETE FROM tokens WHERE expires_at < $1",
		time.Now(),
	)
	return err
}
