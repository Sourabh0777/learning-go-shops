package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"learning-go-shop/internal/config"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
	logger zerolog.Logger
}

func New(cfg *config.Config, db *gorm.DB, logger zerolog.Logger) *Server {
	return &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}
}

func (s *Server) Router() *gin.Engine {
	router := gin.New()

	return router
}
