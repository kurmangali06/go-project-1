package cart

import (
	"context"
	"errors"
	"fmt"
	"go-project-1/internal/clients/product"
	"sort"

	"golang.org/x/sync/errgroup"
)

var ErrProductNotFound = errors.New("product not found")

type cartRepository interface {
	AddItem(userID, skuID int64, count uint16) error
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
	return s.repo.AddItem(userID, skuID, count)
}

func (s *Service) DeleteItem(ctx context.Context, userID, skuID int64) error {
	return s.repo.DeleteItem(userID, skuID)
}

func (s *Service) Clear(ctx context.Context, userID int64) error {
	return s.repo.Clear(userID)
}

func (s *Service) GetItems(ctx context.Context, userID int64) ([]CartItem, uint32, error) {
	raw := s.repo.GetItems(userID)
	skus := make([]int64, 0, len(raw))
	for sku := range raw {
		skus = append(skus, sku)
	}
	sort.Slice(skus, func(i, j int) bool { return skus[i] < skus[j] })

	items := make([]CartItem, len(skus))

	// errgroup сам управляет горутинами и первой ошибкой
	g, gCtx := errgroup.WithContext(ctx)

	for i, sku := range skus {
		i, sku := i, sku // захват переменных (Go < 1.22)
		g.Go(func() error {
			p, err := s.products.GetProduct(gCtx, sku)
			if err != nil {
				return fmt.Errorf("get product %d: %w", sku, err)
			}
			items[i] = CartItem{
				SkuID: sku,
				Name:  p.Name,
				Count: raw[sku],
				Price: p.Price,
			}
			return nil
		})
	}

	// Ждём все горутины, получаем первую ошибку
	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	var totalPrice uint32
	for _, item := range items {
		totalPrice += item.Price * uint32(item.Count)
	}

	return items, totalPrice, nil
}
