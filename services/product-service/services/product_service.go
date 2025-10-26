package services

import (
	"errors"
	"product-service/dto"
	"product-service/models"

	"gorm.io/gorm"
)

// ProductService handles product business logic
type ProductService struct {
	db *gorm.DB
}

// NewProductService creates a new product service
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
	}

	if err := s.db.Create(&product).Error; err != nil {
		return nil, err
	}

	return s.modelToResponse(&product), nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(id uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return s.modelToResponse(&product), nil
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts() ([]dto.ProductResponse, error) {
	var products []models.Product
	if err := s.db.Find(&products).Error; err != nil {
		return nil, err
	}

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, *s.modelToResponse(&product))
	}

	return responses, nil
}

// GetProductsByCategory retrieves products by category
func (s *ProductService) GetProductsByCategory(category string) ([]dto.ProductResponse, error) {
	var products []models.Product
	if err := s.db.Where("category = ?", category).Find(&products).Error; err != nil {
		return nil, err
	}

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, *s.modelToResponse(&product))
	}

	return responses, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(id uint, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Category = req.Category

	if err := s.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return s.modelToResponse(&product), nil
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(id uint) error {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	if err := s.db.Delete(&product).Error; err != nil {
		return err
	}

	return nil
}

// modelToResponse converts a Product model to ProductResponse DTO
func (s *ProductService) modelToResponse(product *models.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
