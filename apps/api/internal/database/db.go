package database

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		slog.Error("Error configuring database", "error", err)
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		slog.Error("Error pinging database", "error", err)
		pool.Close()
		return nil, err
	}
	return pool, nil
}

// NewRedisClient accepts either:
//   - "host:port"                          (local / docker-compose)
//   - "redis://[:password@]host:port"      (Railway, standard URL)
//   - "rediss://[:password@]host:port"     (Railway with TLS)
func NewRedisClient(rawURL string) (*redis.Client, error) {
	var opts *redis.Options
	var err error

	if strings.HasPrefix(rawURL, "redis://") || strings.HasPrefix(rawURL, "rediss://") {
		opts, err = redis.ParseURL(rawURL)
		if err != nil {
			return nil, fmt.Errorf("invalid redis URL %q: %w", rawURL, err)
		}
	} else {
		// bare host:port
		opts = &redis.Options{Addr: rawURL}
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connect failed: %w", err)
	}
	return client, nil
}
