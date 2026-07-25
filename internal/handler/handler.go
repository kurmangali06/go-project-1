package handler

import (
	"context"
	"encoding/json"
	"errors"
	"go-project-1/internal/service/cart"
	"net/http"
	"strconv"
)

type cartService interface {
	AddItem(ctx context.Context, userID, skuID int64, count uint16) error
	DeleteItem(ctx context.Context, userID, skuID int64) error
	Clear(ctx context.Context, userID int64) error
	GetItems(ctx context.Context, userID int64) ([]cart.CartItem, uint32, error)
}

type item struct {
	SkuID int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
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
	userID, skuID, err := checkID(w, r)
	if err != nil {
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
	userID, skuID, err := checkID(w, r)
	if err != nil {
		return
	}

	err = h.service.DeleteItem(r.Context(), userID, skuID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func checkID(w http.ResponseWriter, r *http.Request) (int64, int64, error) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return 0, 0, err
	}
	skuID, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid sku_id", http.StatusBadRequest)
		return 0, 0, err
	}
	return userID, skuID, nil
}

func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Clear(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListCart(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cartItems, totalPrice, err := h.service.GetItems(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(cartItems) == 0 {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(listResponse{Items: []item{}, TotalPrice: 0})
		if err != nil {
			return
		}
		return
	}
	resp := listResponse{TotalPrice: totalPrice}
	for _, ci := range cartItems {
		resp.Items = append(resp.Items, item{SkuID: ci.SkuID, Name: ci.Name, Count: ci.Count, Price: ci.Price})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
