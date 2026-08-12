package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"learning-go-shop/internal/config"
	"learning-go-shop/internal/database"
	"learning-go-shop/internal/logger"
	"learning-go-shop/internal/server"
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
	srv := server.New(cfg, db, log)
	router := srv.SetupRoutes()

	httpServer := http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("Starting http Server")
		// would block forever because ListenAndServe() waits continuously for incoming requests.
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to start http server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal := <-quit
	fmt.Println("signal", signal)
	log.Info().Msg("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = httpServer.Shutdown(ctx)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to shutdown http server")
	}
	log.Info().Msg("shutting down database")
}
