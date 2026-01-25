package database

import (
	"errors"

	"github.com/google/uuid"
)

type MoneyEntry struct {
	ID        uuid.UUID `json:"id"`
	Balance   float64   `json:"balance"`
	Budget    float64   `json:"budget"`
	Ratio     float64   `json:"ratio"`
	CreatedAt string    `json:"created_at"`
	UserID    uuid.UUID `json:"user_id"`
}

func (db *Database) InsertMoneyDB(entry *MoneyEntry) (uuid.UUID, error) {
	var id uuid.UUID
	err := db.DB.QueryRow(
		"INSERT INTO money (balance, budget, ratio, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
		entry.Balance, entry.Budget, entry.Ratio, entry.UserID,
	).Scan(&id)
	entry.ID = id
	return id, err
}

func (db *Database) SelectUserMoneyDB(userID *uuid.UUID) ([]*MoneyEntry, error) {
	if userID == nil {
		return nil, errors.New("userID is nil")
	}
	rows, err := db.DB.Query("SELECT id, balance, budget, ratio, created_at, user_id FROM money WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*MoneyEntry
	for rows.Next() {
		entry := &MoneyEntry{}
		if err := rows.Scan(&entry.ID, &entry.Balance, &entry.Budget, &entry.Ratio, &entry.CreatedAt, &entry.UserID); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (db *Database) SelectUserMoneyByCountDB(userID *uuid.UUID, count int64) ([]*MoneyEntry, error) {
	if userID == nil {
		return nil, errors.New("userID is nil")
	}
	rows, err := db.DB.Query("SELECT id, balance, budget, ratio, created_at, user_id FROM money WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2", userID, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*MoneyEntry
	for rows.Next() {
		entry := &MoneyEntry{}
		if err := rows.Scan(&entry.ID, &entry.Balance, &entry.Budget, &entry.Ratio, &entry.CreatedAt, &entry.UserID); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (db *Database) UpdateMoneyDB(entry *MoneyEntry) error {
	_, err := db.DB.Exec(
		"UPDATE money SET balance = $1, budget = $2, ratio = $3 WHERE id = $4",
		entry.Balance, entry.Budget, entry.Ratio, entry.ID,
	)
	return err
}

func (db *Database) DeleteMoneyDB(id *uuid.UUID) error {
	if id == nil {
		return errors.New("id is nil")
	}
	_, err := db.DB.Exec(
		"DELETE FROM money WHERE id = $1",
		id,
	)
	return err
}
