// Мок Product Service для локальной отладки: настоящий route256.pavl.uk:8080
// недоступен. Отвечает на те же запросы, что ждёт internal/clients/product.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type product struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

// Каталог из шапки cart.http. Всё, чего тут нет, отдаётся как 404,
// чтобы AddItem вернул 412 Precondition Failed.
var catalog = map[uint32]product{
	2008: {Name: "Клавиатура", Price: 1500},
	5000: {Name: "Мышь", Price: 800},
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
		Sku   uint32 `json:"sku"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "empty token", http.StatusUnauthorized)
		return
	}

	p, ok := catalog[req.Sku]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func main() {
	// Порт вынесен в ENV: 8080 на машине разработчика часто уже занят.
	addr := os.Getenv("MOCK_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /get_product", getProduct)

	log.Printf("product mock listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
