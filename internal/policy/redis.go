package policy


import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Connect opens a Redis client and verifies it's reachable before
// handing it back — same reasoning as store.Connect for Postgres.
func Connect(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: addr})

	if err := client.Ping(context.Background()).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return client, nil
}