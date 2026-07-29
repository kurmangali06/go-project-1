package kata

import (
	"encoding/json"
	"net/http"
)

// УРОВЕНЬ 5 — HTTP-хендлер через httptest.NewRecorder.
//
// PriceHandler отдаёт цену товара: GET /price?sku=2008
//
// Контракт:
//   - параметр sku отсутствует или пустой → 400
//   - товара нет в прайсе                 → 404
//   - товар есть → 200, Content-Type: application/json,
//     тело {"sku":"2008","price":1500}
//
// Задача: httptest.NewRequest + httptest.NewRecorder, проверить статус,
// заголовок и тело. Один из статусов на практике окажется не тем.
type PriceHandler struct {
	Prices map[string]uint32
}

type priceResponse struct {
	SKU   string `json:"sku"`
	Price uint32 `json:"price"`
}

func (h PriceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sku := r.URL.Query().Get("sku")
	if sku == "" {
		http.Error(w, "sku required", http.StatusBadRequest)
		return
	}

	price, ok := h.Prices[sku]
	if !ok {
		w.Write([]byte(`{"error":"unknown sku"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(priceResponse{SKU: sku, Price: price})
}
