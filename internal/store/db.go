package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a connection pool to the given database URL.
// TODO: think about what "success" actually means here. Does opening
// a pool guarantee the database is reachable, or does pgxpool.New
// just validate the connection string and defer actual connecting?
// (Worth looking this up rather than assuming — it matters for whether
// you need to Ping() afterward.)
func Connect(url string) (*pgxpool.Pool, error) {
	return nil,nil
}