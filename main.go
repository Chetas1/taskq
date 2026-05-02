// Command taskq boots the task-queue runtime: load config, connect to Redis,
// then block until SIGINT/SIGTERM and shut down cleanly. Producer and worker
// loops will be added behind app.App once the queue protocol stabilizes.
package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chetas1/taskq/config"
	"github.com/Chetas1/taskq/internal/app"
)

func main() {
	if err := run(); err != nil {
		log.Printf("fatal: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a, err := app.New(ctx, cfg)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := a.Close(); cerr != nil && !errors.Is(cerr, context.Canceled) {
			log.Printf("close app: %v", cerr)
		}
	}()

	log.Printf("taskq started; redis=%s db=%d", cfg.Redis.Addr, cfg.Redis.DB)
	<-ctx.Done()
	log.Printf("shutdown signal received")
	return nil
}
