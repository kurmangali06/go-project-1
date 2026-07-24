package handler

import (
	"log"
	"net/http"
	"strconv"
)

type cartService interface {
	AddItem(userID, skuID int64, count uint16) error
	DeleteItem(userID, skuID int64)
	Clear(userID int64)
	GetItems(userID int64) map[int64]uint16
}
type Handler struct {
	service cartService
}

func New(service cartService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	skuID, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// count из тела читаем на 4.3, пока заглушка:
	var count uint16 = 1 // TODO(4.3): распарсить из body

	if err := h.service.AddItem(userID, skuID, count); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	skuID, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.service.DeleteItem(userID, skuID)
	w.WriteHeader(http.StatusOK)
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
	log.Printf("items: %v", items)
	w.WriteHeader(http.StatusOK)
}
