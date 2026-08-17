package server

import (
	"strings"

	"github.com/gin-gonic/gin"

	"learning-go-shop/internal/models"
	"learning-go-shop/internal/util"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Authorization: Bearer JWT
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.UnauthorizedResponse(c, "Authorization header required ")
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			util.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := util.ValidateToken(tokenParts[1], s.config.JWT.Secret)
		if err != nil {
			util.UnauthorizedResponse(c, "Invalid token")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			util.ForbiddenResponse(c, "Forbidden")
			c.Abort()
			return
		}

		if role != string(models.UserRoleAdmin) {
			util.ForbiddenResponse(c, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
