package server

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"learning-go-shop/internal/dto"
	"learning-go-shop/internal/util"
)

// @Summary Get user's cart
// @Description Retrieve current user's shopping cart with all items
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} util.Response{data=dto.CartResponse} "Cart retrieved successfully"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 404 {object} util.Response "Cart not found"
// @Router /cart [get]
func (s *Server) getCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := s.cartService.GetCart(userID)
	if err != nil {
		util.NotFoundResponse(c, "Cart not found")
		return
	}

	util.SuccessResponse(c, "Cart retrieved successfully", cart)
}

// @Summary Add item to cart
// @Description Add a product to the user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddToCartRequest true "Item to add to cart"
// @Success 200 {object} util.Response{data=dto.CartResponse} "Item added to cart successfully"
// @Failure 400 {object} util.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} util.Response "Unauthorized"
// @Router /cart/items [post]
func (s *Server) addToCart(c *gin.Context) {

	userID := c.GetUint("user_id")

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	cart, err := s.cartService.AddToCart(userID, &req)
	if err != nil {
		util.BadRequestResponse(c, "Failed to add item to cart", err)
		return
	}

	util.SuccessResponse(c, "Item added to cart successfully", cart)
}

// @Summary Update cart item quantity
// @Description Update the quantity of an item in the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Param request body dto.UpdateCartItemRequest true "New quantity"
// @Success 200 {object} util.Response{data=dto.CartResponse} "Cart item updated successfully"
// @Failure 400 {object} util.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} util.Response "Unauthorized"
// @Router /cart/items/{id} [put]
func (s *Server) updateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid cart item ID", err)
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	cart, err := s.cartService.UpdateCartItem(userID, uint(id), &req)
	if err != nil {
		util.BadRequestResponse(c, "Failed to update cart item", err)
		return
	}

	util.SuccessResponse(c, "Cart item updated successfully", cart)
}

// @Summary Remove item from cart
// @Description Remove an item from the user's shopping cart
// @Tags Cart
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Success 200 {object} util.Response "Item removed from cart successfully"
// @Failure 400 {object} util.Response "Invalid cart item ID"
// @Failure 401 {object} util.Response "Unauthorized"
// @Router /cart/items/{id} [delete]
func (s *Server) removeFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid cart item ID", err)
		return
	}

	if err := s.cartService.RemoveFromCart(userID, uint(id)); err != nil {
		util.InternalServerErrorResponse(c, "Failed to remove item from cart", err)
		return
	}

	util.SuccessResponse(c, "Item removed from cart successfully", nil)
}
