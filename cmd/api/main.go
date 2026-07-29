package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"learning-go-shop/internal/config"
	"learning-go-shop/internal/database"
	"learning-go-shop/internal/logger"
)

func main() {
	fmt.Println("this is print")
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database connection")
	}
	defer mainDB.Close()
	gin.SetMode(cfg.Server.GinMode)
	log.Info().Msg("starting server")
}
