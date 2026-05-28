package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(rw, req)
	})
}

func (cfg *apiConfig) metricsHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	hitCountText := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	respBytes := ([]byte)(hitCountText)
	rw.Write(respBytes)
}

func (cfg *apiConfig) metricsResetHandler(rw http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	bodyText := fmt.Sprint("Hits Reset Successfully")
	respBytes := ([]byte)(bodyText)
	rw.Write(respBytes)

}
