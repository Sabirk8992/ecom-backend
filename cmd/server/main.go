package main

import (
	"github.com/Sabirk8992/ecom-backend/internal/config"
	"github.com/Sabirk8992/ecom-backend/internal/db"
	"github.com/Sabirk8992/ecom-backend/internal/server"
)

func main() {
	cfg := config.Load()
	dbConn := db.Connect(cfg)
	defer dbConn.Close()
	server.Run(cfg, dbConn)
}