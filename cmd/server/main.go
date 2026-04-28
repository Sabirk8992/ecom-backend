package main

import (
	"github.com/Sabirk8992/ecom-backend/internal/config"
	"github.com/Sabirk8992/ecom-backend/internal/server"
)

func main() {
	cfg := config.Load()
	server.Run(cfg)
}
