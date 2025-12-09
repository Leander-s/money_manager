package model

type MoneyEntry struct {
	ID        int64
	Balance   float64
	Budget    float64
	ratio     float64
	CreatedAt string
	UserID    int64
	UserEmail string
}

func (db *Database) InsertMoneyEntry(entry *MoneyEntry) (int64, error) {
	var id int64
	err := db.DB.QueryRow(
		"INSERT INTO money (balance, budget, ratio, user_id, user_email) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		entry.Balance, entry.Budget, entry.ratio, entry.UserID, entry.UserEmail,
	).Scan(&id)
	entry.ID = id
	return id, err
}

func (db *Database) GetUserMoney(userID int64) ([]*MoneyEntry, error) {
	rows, err := db.DB.Query("SELECT id, balance, budget, ratio, created_at, user_id, user_email FROM money WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*MoneyEntry
	for rows.Next() {
		entry := &MoneyEntry{}
		if err := rows.Scan(&entry.ID, &entry.Balance, &entry.Budget, &entry.ratio, &entry.CreatedAt, &entry.UserID, &entry.UserEmail); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (db *Database) UpdateMoneyEntry(entry *MoneyEntry) error {
	_, err := db.DB.Exec(
		"UPDATE money SET balance = $1, budget = $2, ratio = $3 WHERE id = $4",
		entry.Balance, entry.Budget, entry.ratio, entry.ID,
	)
	return err
}

func (db *Database) DeleteMoneyEntry(id int64) error {
	_, err := db.DB.Exec(
		"DELETE FROM money WHERE id = $1",
		id,
	)
	return err
}
