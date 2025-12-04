package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func openDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dataSourceName)
	fmt.Println("Database source name:", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

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
			return db, nil
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to ping database the %d. time: %w", i, err)
		case <-ticker.C:
		}
	}
}
