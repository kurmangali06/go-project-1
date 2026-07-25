package cart

import (
	"context"
	"errors"
	"sort"

	"go-project-1/internal/clients/product"
)

var ErrProductNotFound = errors.New("product not found")

type cartRepository interface {
	AddItem(userID, skuID int64, count uint16)
	DeleteItem(userID, skuID int64) error
	Clear(userID int64) error
	GetItems(userID int64) map[int64]uint16
}

type productService interface {
	GetProduct(ctx context.Context, sku int64) (product.Product, error)
}
type CartItem struct {
	SkuID int64
	Name  string
	Count uint16
	Price uint32
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
	err := s.repo.DeleteItem(userID, skuID)
	if err != nil {
		return
	}
}

func (s *Service) Clear(userID int64) {
	err := s.repo.Clear(userID)
	if err != nil {
		return
	}
}

func (s *Service) GetItems(ctx context.Context, userID int64) ([]CartItem, uint32, error) {
	raw := s.repo.GetItems(userID)
	skus := make([]int64, 0, len(raw))
	for sku := range raw {
		skus = append(skus, sku)
	}
	sort.Slice(skus, func(i, j int) bool { return skus[i] < skus[j] })

	items := make([]CartItem, 0, len(skus))
	var totalPrice uint32
	for _, sku := range skus {
		count := raw[sku]

		p, err := s.products.GetProduct(ctx, sku)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, CartItem{SkuID: sku, Name: p.Name, Count: count, Price: p.Price})
		totalPrice += p.Price * uint32(count)
	}
	return items, totalPrice, nil
}
