package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

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
