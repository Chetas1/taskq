package main

import (
	"github.com/Chetas1/taskq/config"
	"github.com/Chetas1/taskq/internal/app"
)

func main() {
	cfg := config.LoadConfig()
	app.Init(cfg)
}
