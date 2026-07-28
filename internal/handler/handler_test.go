package handler

import (
	"context"
	"encoding/json"
	"go-project-1/internal/service/cart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubService struct {
	addItem    func(ctx context.Context, userID, skuID int64, count uint16) error
	getItems   func(ctx context.Context, userID int64) ([]cart.CartItem, uint32, error)
	deleteItem func(ctx context.Context, userID, skuID int64) error
	clear      func(ctx context.Context, userID int64) error
}

func (s stubService) AddItem(ctx context.Context, userID, skuID int64, count uint16) error {
	return s.addItem(ctx, userID, skuID, count)
}

func (s stubService) GetItems(ctx context.Context, userID int64) ([]cart.CartItem, uint32, error) {
	return s.getItems(ctx, userID)
}

func (s stubService) Clear(ctx context.Context, userID int64) error {
	return s.clear(ctx, userID)
}

func (s stubService) DeleteItem(ctx context.Context, userID, skuID int64) error {
	return s.deleteItem(ctx, userID, skuID)
}

// listRequest собирает запрос к ListCart. PathValue заполняет ServeMux при
// матчинге роута, а httptest.NewRequest мимо мукса не проходит — проставляем руками.
func listRequest(userID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/user/"+userID+"/cart/list", nil)
	req.SetPathValue("user_id", userID)
	return req
}

func TestListCart_Empty(t *testing.T) {
	h := New(stubService{
		getItems: func(ctx context.Context, userID int64) ([]cart.CartItem, uint32, error) {
			return nil, 0, nil
		},
	})

	rec := httptest.NewRecorder()
	h.ListCart(rec, listRequest("1007"))

	if rec.Code != http.StatusNotFound {
		t.Errorf("статус = %d, ожидали %d", rec.Code, http.StatusNotFound)
	}
}

func TestListCart_NotEmpty(t *testing.T) {
	h := New(stubService{
		getItems: func(ctx context.Context, userID int64) ([]cart.CartItem, uint32, error) {
			items := []cart.CartItem{
				{SkuID: 2008, Name: "Клавиатура", Count: 3, Price: 1500},
				{SkuID: 5000, Name: "Мышь", Count: 1, Price: 800},
			}
			return items, 5300, nil
		},
	})

	rec := httptest.NewRecorder()
	h.ListCart(rec, listRequest("1007"))

	if rec.Code != http.StatusOK {
		t.Fatalf("статус = %d, ожидали %d", rec.Code, http.StatusOK)
	}

	var got listResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("разбор тела ответа: %v (тело: %s)", err, rec.Body.String())
	}

	if got.TotalPrice != 5300 {
		t.Errorf("total_price = %d, ожидали 5300", got.TotalPrice)
	}
	if len(got.Items) != 2 {
		t.Fatalf("позиций = %d, ожидали 2", len(got.Items))
	}
	if got.Items[0].SkuID != 2008 || got.Items[0].Name != "Клавиатура" || got.Items[0].Count != 3 {
		t.Errorf("первая позиция = %+v", got.Items[0])
	}
}

func TestListCart_BadUserID(t *testing.T) {
	h := New(stubService{
		getItems: func(ctx context.Context, userID int64) ([]cart.CartItem, uint32, error) {
			t.Fatal("сервис не должен вызываться при кривом user_id")
			return nil, 0, nil
		},
	})

	rec := httptest.NewRecorder()
	h.ListCart(rec, listRequest("abc"))

	if rec.Code != http.StatusBadRequest {
		t.Errorf("статус = %d, ожидали %d", rec.Code, http.StatusBadRequest)
	}
}
