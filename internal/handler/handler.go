package handler

import (
	"context"
	"encoding/json"
	"errors"
	"go-project-1/internal/service/cart"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type cartService interface {
	AddItem(ctx context.Context, userID, skuID int64, count uint16) error
	DeleteItem(userID, skuID int64)
	Clear(userID int64)
	GetItems(userID int64) map[int64]uint16
}

type item struct {
	SkuID int64  `json:"sku_id"`
	Count uint16 `json:"count"`
}
type listResponse struct {
	Items      []item `json:"items"`
	TotalPrice uint32 `json:"total_price"`
}
type Handler struct {
	service cartService
}

func New(service cartService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, skuID, done := checkID(w, r)
	if done {
		return
	}
	var req struct {
		Count uint16 `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Count == 0 {
		http.Error(w, "count must be greater than 0", http.StatusBadRequest)
		return
	}

	if err := h.service.AddItem(r.Context(), userID, skuID, req.Count); err != nil {
		if errors.Is(err, cart.ErrProductNotFound) {
			http.Error(w, "product not found", http.StatusPreconditionFailed)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, skuID, done := checkID(w, r)
	if done {
		return
	}

	h.service.DeleteItem(userID, skuID)
	w.WriteHeader(http.StatusOK)
}

func checkID(w http.ResponseWriter, r *http.Request) (int64, int64, bool) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, 0, true
	}
	skuID, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, 0, true
	}
	return userID, skuID, false
}

func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.service.Clear(userID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListCart(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	items := h.service.GetItems(userID)
	if len(items) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	skus := make([]int64, 0, len(items))
	for sku := range items {
		skus = append(skus, sku)
	}
	sort.Slice(skus, func(i, j int) bool { return skus[i] < skus[j] })
	resp := listResponse{}
	for _, sku := range skus {
		count := items[sku]
		resp.Items = append(resp.Items, item{SkuID: sku, Count: count})
		resp.TotalPrice += uint32(count)
	}

	log.Printf("items: %v", items)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
