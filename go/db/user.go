package database

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (db *Database) InsertUser(user *User) (int64, error) {
	var id int64
	err := db.DB.QueryRow(
		"INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Password, user.Email,
	).Scan(&id)
	user.ID = id
	return id, err
}

func (db *Database) GetAllUsers() ([]*User, error) {
	rows, err := db.DB.Query("SELECT id, username, password_hash, email, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (db *Database) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := db.DB.QueryRow(
		"SELECT id, username, password_hash, email, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
	return user, err
}

func (db *Database) GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := db.DB.QueryRow(
		"SELECT id, username, password_hash, email, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
	return user, err
}

func (db *Database) UpdateUser(user *User) error {
	_, err := db.DB.Exec(
		"UPDATE users SET username = $1, password_hash = $2, email = $3 WHERE id = $4",
		user.Username, user.Password, user.Email, user.ID,
	)
	return err
}

func (db *Database) DeleteUser(id int64) error {
	_, err := db.DB.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	return err
}
