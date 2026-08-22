package server

import (
	"github.com/gin-gonic/gin"

	"learning-go-shop/internal/dto"
	"learning-go-shop/internal/services"
	"learning-go-shop/internal/util"
)

// We deal only with DTO here
// Not going to see model
// Extract the data and then bind it into the request

// Models layer is handeling its things
// DTo layer is handeling its things
// Services layer is handeling its things
// THen comes the service layer

func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := services.NewAuthService(s.db, s.config)
	response, err := authService.Register(&req)
	if err != nil {
		util.BadRequestResponse(c, "Registration failed", err)
		return
	}

	util.CreatedResponse(c, "User registered successfully", response)
}
func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := services.NewAuthService(s.db, s.config)
	response, err := authService.Login(&req)
	if err != nil {
		util.UnauthorizedResponse(c, "Login failed")
		return
	}

	util.SuccessResponse(c, "Login successful", response)
}
func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := services.NewAuthService(s.db, s.config)

	response, err := authService.RefreshToken(&req)
	if err != nil {
		util.UnauthorizedResponse(c, "Token refresh failed")
		return
	}

	util.SuccessResponse(c, "Token refreshed successfully", response)
}

// func (s *Server) logout(c *gin.Context) {
// 	var req dto.RefreshTokenRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		util.BadRequestResponse(c, "Invalid request data", err)
// 		return
// 	}
// 	authService := services.NewAuthService(s.db, s.config)
// 	_, err := authService.Logout(&req)
// 	if err != nil {
// 		util.InternalServerErrorResponse(c, "Logout failed", err)
// 		return
// 	}

// 	util.SuccessResponse(c, "Logout successful", nil)
// }

func (s *Server) getProfile(c *gin.Context) {
	userId := c.GetUint("user_id")
	authService := services.NewUserService(s.db)
	profile, err := authService.GetProfile(userId)
	if err != nil {
		util.NotFoundResponse(c, "User not found")
		return
	}
	util.SuccessResponse(c, "Profile retrieved successfully", profile)
}

func (s *Server) updateProfile(c *gin.Context) {
	userId := c.GetUint("user_id")
	var req dto.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := services.NewUserService(s.db)
	profile, err := authService.UpdateProfile(userId, &req)
	if err != nil {
		util.NotFoundResponse(c, "User not found")
		return
	}
	util.SuccessResponse(c, "Profile updated successfully", profile)
}
