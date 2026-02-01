package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/google/uuid"
)

type UserStore interface {
	// User-related methods
	InsertUserDB(userForInsert *UserForInsert) (User, error)
	SelectAllUsersDB() ([]*User, error)
	SelectUserByEmailDB(email string) (*User, error)
	SelectUserByIDDB(id *uuid.UUID) (*User, error)
	UpdateUserDB(user *UserForUpdate) error
	DeleteUserDB(id *uuid.UUID) error 
}

type RoleStore interface {
	// Role-related methods
	GetUserRolesDB(userID *uuid.UUID) ([]Role, error) 
	AssignRoleToUserDB(userID *uuid.UUID, role string) error 
	RemoveRoleFromUserDB(userID *uuid.UUID, role string) error 
	CheckUserRoleDB(userID *uuid.UUID, role string) (bool, error) 
}

type TokenStore interface {
	// Token-related methods
	InsertToken(Token *Token) error 
	GetToken(token *uuid.UUID) (*Token, error) 
	ListTokens() ([]*Token, error) 
	DeleteToken(token *uuid.UUID) error 
	DeleteTokensByUserID(userID *uuid.UUID) error 
	DeleteExpiredTokens() error 
}

type UserRoleStore interface {
	UserStore
	RoleStore
}

type AuthStore interface {
	// Authentication-related methods would go here
	TokenStore
	UserRoleStore
}

type MoneyStore interface {
	// Money-related methods would go here
	InsertMoneyDB(entry *MoneyEntry) (uuid.UUID, error) 
	SelectMoneyByIDDB(id *uuid.UUID) (*MoneyEntry, error)
	SelectUserMoneyDB(userID *uuid.UUID) ([]*MoneyEntry, error) 
	SelectUserMoneyByCountDB(userID *uuid.UUID, count int64) ([]*MoneyEntry, error) 
	UpdateMoneyBatchDB(entries []*MoneyEntry) error
	UpdateMoneyDB(entry *MoneyEntry) error 
	DeleteMoneyDB(id *uuid.UUID) error
}

type DatabaseInterface interface {
	AuthStore
	MoneyStore

	Close() error
}

type Database struct {
	DB *sql.DB
}

func OpenDB(dataSourceName string) (Database, error) {
	result := Database{DB: nil}
	db, err := sql.Open("pgx", dataSourceName)
	fmt.Println("Database source name:", dataSourceName)
	if err != nil {
		return result, fmt.Errorf("failed to open database: %w", err)
	}
	result.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	i := 0
	for {
		i++	
		if err := db.PingContext(ctx); err != nil {
			fmt.Println("Waiting for database to be ready...")
		} else {
			fmt.Println("Database is ready!")
			return result, nil
		}

		select {
		case <-ctx.Done():
			result.DB.Close()
			result.DB = nil
			return result, fmt.Errorf("failed to ping database the %d. time: %w", i, err)
		case <-ticker.C:
		}
	}
}

func (database *Database) Close() error {
	if database.DB != nil {
		return database.DB.Close()
	}
	return nil
}
