package middleware

import (
	"log"
	"net/http"
	"time"
)

type WrapperHandler struct {
	Wrap http.Handler
}

func (h WrapperHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	rec := &statusRecorder{ResponseWriter: writer, status: http.StatusOK}

	h.Wrap.ServeHTTP(rec, request)

	log.Printf("%s %s -> %d (%s)", request.Method, request.URL.Path, rec.status, time.Since(start))
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
