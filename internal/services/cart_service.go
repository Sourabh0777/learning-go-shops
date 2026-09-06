package services

import (
	"gorm.io/gorm"

	"learning-go-shop/internal/dto"
	"learning-go-shop/internal/models"
)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{db: db}
}

func (s *CartService) GetCart(userID uint) (*dto.CartItemResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Product.Category").
		Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return s.convertToCartResponse(&cart), nil
}

func (s *CartService) convertToCartResponse(cart *models.Cart) *dto.CartResponse {

	cartItems := make([]dto.CartItemResponse, len(cart.CartItems)) // memory allocation
	var total float64

	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
			Quantity:  cart.CartItems[i].Quantity,
			Subtotal:  subtotal,
			CreatedAt: cart.CartItems[i].CreatedAt,
			UpdatedAt: cart.CartItems[i].UpdatedAt,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
