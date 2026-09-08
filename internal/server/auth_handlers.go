package server

import (
	"github.com/gin-gonic/gin"

	"learning-go-shop/internal/dto"
	"learning-go-shop/internal/util"
)

// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} util.Response{data=dto.AuthResponse} "User registered successfully"
// @Failure 400 {object} util.Response "Invalid request data or user already exists"
// @Router /auth/register [post]
func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.Register(&req)
	if err != nil {
		util.BadRequestResponse(c, "Registration failed", err)
		return
	}

	util.CreatedResponse(c, "User registered successfully", response)
}

// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200 {object} util.Response{data=dto.AuthResponse} "Login successful"
// @Failure 401 {object} util.Response "Invalid credentials"
// @Router /auth/login [post]
func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.Login(&req)
	if err != nil {
		util.UnauthorizedResponse(c, "Login failed")
		return
	}

	util.SuccessResponse(c, "Login successful", response)
}

// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} util.Response{data=dto.AuthResponse} "Token refreshed successfully"
// @Failure 401 {object} util.Response "Invalid refresh token"
// @Router /auth/refresh [post]
func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.RefreshToken(&req)
	if err != nil {
		util.UnauthorizedResponse(c, "Token refresh failed")
		return
	}

	util.SuccessResponse(c, "Token refreshed successfully", response)
}

// @Summary Get user profile
// @Description Get current authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} util.Response{data=dto.UserResponse} "Profile retrieved successfully"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 404 {object} util.Response "User not found"
// @Router /users/profile [get]
func (s *Server) getProfile(c *gin.Context) {
	userId := c.GetUint("user_id")
	profile, err := s.userService.GetProfile(userId)
	if err != nil {
		util.NotFoundResponse(c, "User not found")
		return
	}
	util.SuccessResponse(c, "Profile retrieved successfully", profile)
}

// @Summary Update user profile
// @Description Update current authenticated user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} util.Response{data=dto.UserResponse} "Profile updated successfully"
// @Failure 400 {object} util.Response "Invalid request data"
// @Failure 401 {object} util.Response "Unauthorized"
// @Router /users/profile [put]
func (s *Server) updateProfile(c *gin.Context) {
	userId := c.GetUint("user_id")
	var req dto.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	profile, err := s.userService.UpdateProfile(userId, &req)
	if err != nil {
		util.NotFoundResponse(c, "User not found")
		return
	}
	util.SuccessResponse(c, "Profile updated successfully", profile)
}
