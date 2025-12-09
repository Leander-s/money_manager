package model

import (
	"github.com/google/uuid"
)

type Token struct {
	Token     uuid.UUID
	UserID    int64
	Expiry    int64
}

func (db *Database) InsertToken(Token *Token) error {
	err := db.DB.QueryRow(
		"INSERT INTO tokens (Token, user_id, Expiry) VALUES ($1, $2, $3, $4) RETURNING id",
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
		"SELECT user_id, Expiry FROM tokens WHERE Token = $1",
		token,
	).Scan(&Token.UserID, &Token.Expiry)
	Token.Token = token 
	return Token, err
}

func (db *Database) DeleteToken(token uuid.UUID) error {
	_, err := db.DB.Exec(
		"DELETE FROM tokens WHERE Token = $1",
		token,
	)
	return err
}
