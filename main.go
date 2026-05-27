package main

import (
	"net/http"
)

func main() {

	serveMux := http.NewServeMux()
	rootFS := http.FileServer(http.Dir("."))
	rootPrefix := http.StripPrefix("/app/", rootFS)
	assetsFS := http.FileServer(http.Dir("./assets"))
	assetsPrefix := http.StripPrefix("/assets/", assetsFS)
	serveMux.Handle("/app/", rootPrefix)
	serveMux.Handle("/assets/", assetsPrefix)
	serveMux.HandleFunc("/healthz", healthHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()
}

func healthHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	okBody := ([]byte)("OK")
	rw.Write(okBody)
}
