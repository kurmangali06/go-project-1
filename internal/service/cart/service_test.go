package cart

import (
	"context"
	"errors"
	"go-project-1/internal/clients/product"
	"go-project-1/internal/repository/memory"
	"sync/atomic"
	"testing"
)

// GetItems ходит в Product Service параллельно (errgroup), поэтому счётчик
// вызовов обязан быть потокобезопасным, а ресивер — указателем.
type stubProducts struct {
	calls      atomic.Int64
	getProduct func(ctx context.Context, sku int64) (product.Product, error)
}

func (s *stubProducts) GetProduct(ctx context.Context, sku int64) (product.Product, error) {
	s.calls.Add(1)
	return s.getProduct(ctx, sku)
}

// каталог, который отдаёт цены как настоящий Product Service
func catalogStub() *stubProducts {
	prices := map[int64]product.Product{
		2008: {Name: "Клавиатура", Price: 1500},
		5000: {Name: "Мышь", Price: 800},
	}
	return &stubProducts{
		getProduct: func(ctx context.Context, sku int64) (product.Product, error) {
			p, ok := prices[sku]
			if !ok {
				return product.Product{}, product.ErrProductNotFound
			}
			return p, nil
		},
	}
}

func TestAddItem_UnknownSKU(t *testing.T) {
	repo := memory.NewRepository()
	svc := NewService(repo, catalogStub())

	err := svc.AddItem(context.Background(), 1007, 9999, 1)

	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("ошибка = %v, ожидали ErrProductNotFound", err)
	}
	if items := repo.GetItems(1007); len(items) != 0 {
		t.Errorf("несуществующий товар попал в корзину: %v", items)
	}
}

func TestAddItem_OK(t *testing.T) {
	repo := memory.NewRepository()
	svc := NewService(repo, catalogStub())

	if err := svc.AddItem(context.Background(), 1007, 2008, 3); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if got := repo.GetItems(1007)[2008]; got != 3 {
		t.Errorf("count = %d, ожидали 3", got)
	}
}

// Ошибка Product Service, не связанная с «товар не найден», не должна
// подменяться на ErrProductNotFound — иначе хендлер отдаст 412 вместо 500.
func TestAddItem_ProductServiceDown(t *testing.T) {
	repo := memory.NewRepository()
	boom := errors.New("connection refused")
	svc := NewService(repo, &stubProducts{
		getProduct: func(ctx context.Context, sku int64) (product.Product, error) {
			return product.Product{}, boom
		},
	})

	err := svc.AddItem(context.Background(), 1007, 2008, 1)

	if !errors.Is(err, boom) {
		t.Fatalf("ошибка = %v, ожидали проброс %v", err, boom)
	}
	if errors.Is(err, ErrProductNotFound) {
		t.Error("сетевая ошибка подменена на ErrProductNotFound → хендлер отдаст 412 вместо 500")
	}
}

func TestGetItems(t *testing.T) {
	repo := memory.NewRepository()
	// кладём в обратном порядке — проверим, что сервис отсортирует
	if err := repo.AddItem(1007, 5000, 1); err != nil {
		t.Fatalf("repo.AddItem: %v", err)
	}
	if err := repo.AddItem(1007, 2008, 3); err != nil {
		t.Fatalf("repo.AddItem: %v", err)
	}

	products := catalogStub()
	svc := NewService(repo, products)

	items, totalPrice, err := svc.GetItems(context.Background(), 1007)
	if err != nil {
		t.Fatalf("GetItems: %v", err)
	}

	if totalPrice != 5300 {
		t.Errorf("total_price = %d, ожидали 5300 (1500*3 + 800*1)", totalPrice)
	}
	if len(items) != 2 {
		t.Fatalf("позиций = %d, ожидали 2", len(items))
	}
	if items[0].SkuID != 2008 || items[1].SkuID != 5000 {
		t.Errorf("порядок = [%d %d], ожидали по возрастанию [2008 5000]", items[0].SkuID, items[1].SkuID)
	}
	if items[0].Name != "Клавиатура" || items[0].Count != 3 || items[0].Price != 1500 {
		t.Errorf("первая позиция = %+v", items[0])
	}
	if got := products.calls.Load(); got != 2 {
		t.Errorf("походов в Product Service = %d, ожидали 2 (по одному на sku)", got)
	}
}

func TestGetItems_EmptyCart(t *testing.T) {
	repo := memory.NewRepository()
	products := catalogStub()
	svc := NewService(repo, products)

	items, totalPrice, err := svc.GetItems(context.Background(), 1007)
	if err != nil {
		t.Fatalf("GetItems: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("позиций = %d, ожидали 0", len(items))
	}
	if totalPrice != 0 {
		t.Errorf("total_price = %d, ожидали 0", totalPrice)
	}
	if got := products.calls.Load(); got != 0 {
		t.Errorf("походов в Product Service = %d, на пустой корзине ходить некуда", got)
	}
}

// Если Product Service отвалился хотя бы по одному sku — отдаём ошибку целиком,
// а не половину корзины с нулевыми ценами.
func TestGetItems_ProductError(t *testing.T) {
	repo := memory.NewRepository()
	if err := repo.AddItem(1007, 2008, 3); err != nil {
		t.Fatalf("repo.AddItem: %v", err)
	}
	if err := repo.AddItem(1007, 5000, 1); err != nil {
		t.Fatalf("repo.AddItem: %v", err)
	}

	boom := errors.New("connection refused")
	svc := NewService(repo, &stubProducts{
		getProduct: func(ctx context.Context, sku int64) (product.Product, error) {
			if sku == 5000 {
				return product.Product{}, boom
			}
			return product.Product{Name: "Клавиатура", Price: 1500}, nil
		},
	})

	items, totalPrice, err := svc.GetItems(context.Background(), 1007)

	if !errors.Is(err, boom) {
		t.Fatalf("ошибка = %v, ожидали проброс %v", err, boom)
	}
	if items != nil || totalPrice != 0 {
		t.Errorf("при ошибке вернули данные: items=%v total=%d", items, totalPrice)
	}
}

func TestDeleteItemAndClear(t *testing.T) {
	repo := memory.NewRepository()
	svc := NewService(repo, catalogStub())
	ctx := context.Background()

	if err := svc.AddItem(ctx, 1007, 2008, 3); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := svc.AddItem(ctx, 1007, 5000, 1); err != nil {
		t.Fatalf("AddItem: %v", err)
	}

	if err := svc.DeleteItem(ctx, 1007, 2008); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}
	if _, ok := repo.GetItems(1007)[2008]; ok {
		t.Error("sku 2008 остался после DeleteItem")
	}

	if err := svc.Clear(ctx, 1007); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if items := repo.GetItems(1007); len(items) != 0 {
		t.Errorf("после Clear корзина не пуста: %v", items)
	}
}
