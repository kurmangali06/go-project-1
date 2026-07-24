package main

import (
	"go-project-1/internal/clients/product"
	"go-project-1/internal/handler"
	"go-project-1/internal/repository/memory"
	"go-project-1/internal/service/cart"
	"log"
	"net/http"
	"os"
)

func main() {
	repo := memory.NewRepository()
	client := product.New("http://route256.pavl.uk:8080", os.Getenv("PRODUCT_TOKEN"))
	service := cart.NewService(repo, client)
	h := handler.New(service)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", h.AddItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", h.DeleteItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart", h.Clear)
	mux.HandleFunc("GET /user/{user_id}/cart/list", h.ListCart)
	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	log.Println("listening on :8082")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
