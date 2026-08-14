package store

import (
	"context"
	"time"
	"fmt"


	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a connection pool to the given database URL.
// TODO: think about what "success" actually means here. Does opening
// a pool guarantee the database is reachable, or does pgxpool.New
// just validate the connection string and defer actual connecting?
// (Worth looking this up rather than assuming — it matters for whether
// you need to Ping() afterward.)
func Connect(url string) (*pgxpool.Pool, error) {
	// create the connection time limit 
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.Close()

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	
	return db, nil
}