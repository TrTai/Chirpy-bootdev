package main

import (
	//"fmt"
	"net/http"
	//"sync/atomic"
)

func main() {

	apiCfg := new(apiConfig)

	mux := http.NewServeMux()
	rootFS := http.FileServer(http.Dir("."))
	rootPrefix := http.StripPrefix("/app/", rootFS)
	assetsFS := http.FileServer(http.Dir("./assets"))
	assetsPrefix := http.StripPrefix("/assets/", assetsFS)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(rootPrefix))
	mux.Handle("/assets/", assetsPrefix)
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /reset", apiCfg.metricsResetHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func healthHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	okBody := ([]byte)("OK")
	rw.Write(okBody)
}
