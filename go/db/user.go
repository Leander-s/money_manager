package database

type User struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	CreatedAt     string `json:"created_at"`
	EmailVerified bool   `json:"email_verified"`
}

type UserForInsert struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
}

type UserForUpdate struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (db *Database) InsertUserDB(userForInsert *UserForInsert) (User, error) {
	var user User
	err := db.DB.QueryRow(
		"INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id, username, password_hash, email, created_at",
		userForInsert.Username, userForInsert.PasswordHash, userForInsert.Email,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
	return user, err
}

func (db *Database) SelectAllUsersDB() ([]*User, error) {
	rows, err := db.DB.Query("SELECT id, username, password_hash, email, created_at, email_verified FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.EmailVerified); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (db *Database) SelectUserByEmailDB(email string) (*User, error) {
	user := &User{}
	err := db.DB.QueryRow(
		"SELECT id, username, password_hash, email, created_at, email_verified FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.EmailVerified)
	return user, err
}

func (db *Database) SelectUserByIDDB(id int64) (*User, error) {
	user := &User{}
	err := db.DB.QueryRow(
		"SELECT id, username, password_hash, email, created_at, email_verified FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.EmailVerified)
	return user, err
}

func (db *Database) UpdateUserDB(user *UserForUpdate) error {
	_, err := db.DB.Exec(
		"UPDATE users SET username = $1, password_hash = $2, email = $3, email_verified = $4 WHERE id = $5",
		user.Username, user.Password, user.Email, user.EmailVerified, user.ID,
	)
	return err
}

func (db *Database) DeleteUserDB(id int64) error {
	_, err := db.DB.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	return err
}
