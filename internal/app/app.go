package app

import (
	"fmt"
	"sync"

	"github.com/Chetas1/taskq/config"
	"github.com/go-redis/redis"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func Init(cfg *config.Config) {
	initRedis(cfg)
}

func initRedis(cfg *config.Config) {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		err := rdb.Ping().Err()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print("Connected to Redis!")
	})
}
