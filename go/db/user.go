package database

import (
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	CreatedAt     string    `json:"created_at"`
	EmailVerified bool      `json:"email_verified"`
}

type UserForInsert struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
}

type UserForUpdate struct {
	ID            uuid.UUID `json:"id"`
	Username      *string    `json:"username"`
	Password      *string    `json:"password"`
	Email         *string    `json:"email"`
	EmailVerified *bool      `json:"email_verified"`
}

type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
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

func (db *Database) SelectUserByIDDB(id *uuid.UUID) (*User, error) {
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

func (db *Database) DeleteUserDB(id *uuid.UUID) error {
	_, err := db.DB.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	return err
}

func (db *Database) GetUserRolesDB(userID *uuid.UUID) ([]Role, error) {
	rows, err := db.DB.Query(
		`SELECT r.id, r.name
		 FROM roles r
		 JOIN user_roles ur ON r.id = ur.role_id
		 WHERE ur.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		role := Role{}
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (db *Database) AssignRoleToUserDB(userID *uuid.UUID, role string) error {
	_, err := db.DB.Exec(
		`INSERT INTO user_roles (user_id, role_id)
		 SELECT $1, r.id FROM roles r WHERE r.name = $2`,
		userID, role,
	)
	return err
}

func (db *Database) RemoveRoleFromUserDB(userID *uuid.UUID, role string) error {
	_, err := db.DB.Exec(
		`DELETE FROM user_roles
		 WHERE user_id = $1 AND role_id = (SELECT id FROM roles WHERE name = $2)`,
		userID, role,
	)
	return err
}

func (db *Database) CheckUserRoleDB(userID *uuid.UUID, role string) (bool, error) {
	var count int
	err := db.DB.QueryRow(
		`SELECT COUNT(*)
		 FROM user_roles ur
		 JOIN roles r ON ur.role_id = r.id
		 WHERE ur.user_id = $1 AND r.name = $2`,
		userID, role,
	).Scan(&count)
	return count > 0, err
}
