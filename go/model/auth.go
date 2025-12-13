package model

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	Token  uuid.UUID `json:"token"`
	UserID int64     `json:"userID"`
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

func (db *Database) GetToken(token uuid.UUID) (*Token, error) {
	Token := &Token{}
	err := db.DB.QueryRow(
		"SELECT user_id, expires_at FROM tokens WHERE token = $1",
		token,
	).Scan(&Token.UserID, &Token.Expiry)
	Token.Token = token
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

func (db *Database) DeleteToken(token uuid.UUID) error {
	_, err := db.DB.Exec(
		"DELETE FROM tokens WHERE token = $1",
		token,
	)
	return err
}
