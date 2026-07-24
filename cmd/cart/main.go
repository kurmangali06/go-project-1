package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	log.Println("listening on :8082")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
