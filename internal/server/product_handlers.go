package server

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"learning-go-shop/internal/dto"
	"learning-go-shop/internal/services"
	"learning-go-shop/internal/util"
)

// @Summary Create a new category
// @Description Create a new product category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} util.Response{data=dto.CategoryResponse} "Category created successfully"
// @Failure 400 {object} util.Response "Invalid request data"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 403 {object} util.Response "Admin access required"
// @Router /categories [post]
func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	productService := services.NewProductService(s.db)
	category, err := productService.CreateCategory(&req)
	if err != nil {
		util.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}
	util.CreatedResponse(c, "Category created successfully", category)
}

// @Summary Get all categories
// @Description Retrieve all active categories
// @Tags Categories
// @Produce json
// @Success 200 {object} util.Response{data=[]dto.CategoryResponse} "Categories retrieved successfully"
// @Failure 500 {object} util.Response "Internal server error"
// @Router /categories [get]
func (s *Server) getCategories(c *gin.Context) {
	productService := services.NewProductService(s.db)

	categories, err := productService.GetCategories()
	if err != nil {
		util.InternalServerErrorResponse(c, "Failed to fetch categories", err)
		return
	}

	util.SuccessResponse(c, "Categories retrieved successfully", categories)
}

// @Summary Update a category
// @Description Update an existing category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} util.Response{data=dto.CategoryResponse} "Category updated successfully"
// @Failure 400 {object} util.Response "Invalid request data"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 403 {object} util.Response "Admin access required"
// @Router /categories/{id} [put]
func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid category ID", err)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	productService := services.NewProductService(s.db)

	category, err := productService.UpdateCategory(uint(id), &req)
	if err != nil {
		util.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}

	util.SuccessResponse(c, "Category updated successfully", category)
}

// @Summary Delete a category
// @Description Delete a category (Admin only)
// @Tags Categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} util.Response "Category deleted successfully"
// @Failure 400 {object} util.Response "Invalid category ID"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 403 {object} util.Response "Admin access required"
// @Router /categories/{id} [delete]
func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid category ID", err)
		return
	}
	productService := services.NewProductService(s.db)

	if err := productService.DeleteCategory(uint(id)); err != nil {
		util.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}

	util.SuccessResponse(c, "Category deleted successfully", nil)
}

// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateProductRequest true "Product data"
// @Success 201 {object} util.Response{data=dto.ProductResponse} "Product created successfully"
// @Failure 400 {object} util.Response "Invalid request data"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 403 {object} util.Response "Admin access required"
// @Router /products [post]
func (s *Server) createProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	productService := services.NewProductService(s.db)

	product, err := productService.CreateProduct(&req)
	if err != nil {
		util.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}

	util.CreatedResponse(c, "Product created successfully", product)
}

// @Summary Get all products
// @Description Retrieve paginated list of active products
// @Tags Products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} util.PaginatedResponse{data=[]dto.ProductResponse} "Products retrieved successfully"
// @Failure 500 {object} util.Response "Internal server error"
// @Router /products [get]
func (s *Server) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	productService := services.NewProductService(s.db)

	products, meta, err := productService.GetProducts(page, limit)
	if err != nil {
		util.InternalServerErrorResponse(c, "Failed to fetch products", err)
		return
	}

	util.PaginatedSuccessResponse(c, "Products retrieved successfully", products, *meta)
}

// @Summary Get a product by ID
// @Description Retrieve detailed information about a specific product
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} util.Response{data=dto.ProductResponse} "Product retrieved successfully"
// @Failure 400 {object} util.Response "Invalid product ID"
// @Failure 404 {object} util.Response "Product not found"
// @Router /products/{id} [get]
func (s *Server) getProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid product ID", err)
		return
	}
	productService := services.NewProductService(s.db)

	product, err := productService.GetProduct(uint(id))
	if err != nil {
		util.NotFoundResponse(c, "Product not found")
		return
	}

	util.SuccessResponse(c, "Product retrieved successfully", product)
}

// @Summary Update a product
// @Description Update an existing product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product update data"
// @Success 200 {object} util.Response{data=dto.ProductResponse} "Product updated successfully"
// @Failure 400 {object} util.Response "Invalid request data"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 403 {object} util.Response "Admin access required"
// @Router /products/{id} [put]
func (s *Server) updateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	productService := services.NewProductService(s.db)

	product, err := productService.UpdateProduct(uint(id), &req)
	if err != nil {
		util.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}

	util.SuccessResponse(c, "Product updated successfully", product)
}

// @Summary Delete a product
// @Description Delete a product (Admin only)
// @Tags Products
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} util.Response "Product deleted successfully"
// @Failure 400 {object} util.Response "Invalid product ID"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 403 {object} util.Response "Admin access required"
// @Router /products/{id} [delete]
func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid product ID", err)
		return
	}
	productService := services.NewProductService(s.db)

	if err := productService.DeleteProduct(uint(id)); err != nil {
		util.InternalServerErrorResponse(c, "Failed to delete product", err)
		return
	}

	util.SuccessResponse(c, "Product deleted successfully", nil)
}

// @Summary Upload product image
// @Description Upload an image for a product (Admin only)
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param image formData file true "Image file"
// @Success 200 {object} util.Response{data=map[string]string} "Image uploaded successfully"
// @Failure 400 {object} util.Response "Invalid request or file"
// @Failure 401 {object} util.Response "Unauthorized"
// @Failure 403 {object} util.Response "Admin access required"
// @Router /products/{id}/images [post]
// Pending
func (s *Server) uploadProductImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		util.BadRequestResponse(c, "No file uploaded", err)
		return
	}
	productService := services.NewProductService(s.db)

	// url, err := s.uploadService.UploadProductImage(uint(id), file)
	// if err != nil {
	// 	util.InternalServerErrorResponse(c, "Failed to upload image", err)
	// 	return
	// }

	if err := productService.AddProductImage(uint(id), "url", file.Filename); err != nil {
		util.InternalServerErrorResponse(c, "Failed to save image record", err)
		return
	}

	util.SuccessResponse(c, "Image uploaded successfully", map[string]string{"url": "url"})
}

// @Summary Search products
// @Description Search products using full-text search with ranking
// @Tags Products
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category_id query int false "Filter by category ID"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Success 200 {object} util.PaginatedResponse{data=[]dto.ProductSearchResult} "Search results"
// @Failure 400 {object} util.Response "Invalid search query"
// @Failure 500 {object} util.Response "Internal server error"
// @Router /search [get]
func (s *Server) searchProducts(c *gin.Context) {
	var req dto.SearchProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.BadRequestResponse(c, "Invalid search parameters", err)
		return
	}
	productService := services.NewProductService(s.db)

	results, meta, err := productService.SearchProducts(&req)
	if err != nil {
		s.logger.Error().Err(err).Msg("Product search failed")
		util.InternalServerErrorResponse(c, "Search failed", errors.New("unable to complete search at this time"))
		return
	}

	util.PaginatedSuccessResponse(c, "OK", results, *meta)
}
