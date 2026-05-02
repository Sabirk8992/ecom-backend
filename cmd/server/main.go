package main

import (
	"github.com/Sabirk8992/ecom-backend/internal/config"
	"github.com/Sabirk8992/ecom-backend/internal/db"
	"github.com/Sabirk8992/ecom-backend/internal/logger"
	"github.com/Sabirk8992/ecom-backend/internal/metrics"
	"github.com/Sabirk8992/ecom-backend/internal/server"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.AppEnv)
	metrics.Init()
	dbConn := db.Connect(cfg)
	defer dbConn.Close()
	server.Run(cfg, dbConn)
}
