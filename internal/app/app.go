// Package app wires the taskq runtime: config -> Redis client -> (future)
// producer + worker pool. The current implementation only verifies Redis
// connectivity; queue logic is the next milestone (see README roadmap).
package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Chetas-Patil/taskq/config"
	"github.com/go-redis/redis/v8"
)

// pingTimeout bounds the initial Redis health-check so a misconfigured
// Addr/Password fails fast at boot instead of hanging an operator.
const pingTimeout = 5 * time.Second

// App holds runtime dependencies. It is intentionally non-singleton so multiple
// instances (e.g. in tests) can coexist without sharing global state.
type App struct {
	Redis *redis.Client
}

// New constructs an App, connects to Redis, and verifies the connection with
// PING. Returns a wrapped error if the connection or PING fails.
func New(ctx context.Context, cfg *config.Config) (*App, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping %s: %w", cfg.Redis.Addr, err)
	}

	return &App{Redis: rdb}, nil
}

// Close releases the Redis connection pool. Safe to call multiple times.
func (a *App) Close() error {
	if a == nil || a.Redis == nil {
		return nil
	}
	if err := a.Redis.Close(); err != nil {
		return fmt.Errorf("close redis: %w", err)
	}
	return nil
}
