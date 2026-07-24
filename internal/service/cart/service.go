package cart

import (
	"context"
	"errors"

	"go-project-1/internal/clients/product"
)

var ErrProductNotFound = errors.New("product not found")

type cartRepository interface {
	AddItem(userID, skuID int64, count uint16)
	DeleteItem(userID, skuID int64)
	Clear(userID int64)
	GetItems(userID int64) map[int64]uint16
}

type productService interface {
	GetProduct(ctx context.Context, sku int64) (product.Product, error)
}

type Service struct {
	repo     cartRepository
	products productService
}

func NewService(repo cartRepository, products productService) *Service {
	return &Service{repo: repo, products: products}
}

func (s *Service) AddItem(ctx context.Context, userID, skuID int64, count uint16) error {
	_, err := s.products.GetProduct(ctx, skuID)
	if err != nil {
		if errors.Is(err, product.ErrProductNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	s.repo.AddItem(userID, skuID, count)
	return nil
}

func (s *Service) DeleteItem(userID, skuID int64) {
	s.repo.DeleteItem(userID, skuID)
}

func (s *Service) Clear(userID int64) {
	s.repo.Clear(userID)
}

func (s *Service) GetItems(userID int64) map[int64]uint16 {
	return s.repo.GetItems(userID)
}
