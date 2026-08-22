package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"learning-go-shop/internal/config"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
	logger *zerolog.Logger
}

func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) *Server {
	return &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())
	router.GET("/health", s.healthCheck)
	api := router.Group("/api/v1")
	{

		auth := api.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			// auth.POST("/logout", s.register)1

			{
				protected := auth.Group("/")
				protected.Use(s.authMiddleware())
				{
					users := protected.Group("users")
					{
						userRoutes := users
						userRoutes.GET("/profile", s.getProfile)
						userRoutes.POST("/profile", s.updateProfile)
					}

				}
			}

		}
	}
	// asdsdasd
	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Method", "GET,POST,PUT,DELETE,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Context-Type,Authorization")
		if ctx.Request.Method == "QPTIONS" {
			ctx.AbortWithStatus(204)
		}
	}
}
