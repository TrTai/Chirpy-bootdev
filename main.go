package main

import (
	"database/sql"
	//"encoding/json"
	//"fmt"
	"os"

	"log"
	"net/http"

	//"sync/atomic"
	//"slices"
	//"strings"

	"github.com/TrTai/Chirpy-bootdev/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening db: %s", err)
	}

	dbQueries := database.New(db)

	apiCfg := new(apiConfig)
	apiCfg.dbQueries = dbQueries
	apiCfg.platform = platform

	mux := http.NewServeMux()

	AssignMuxHandles(mux, apiCfg)
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

func AssignMuxHandles(mux *http.ServeMux, apiCfg *apiConfig) {
	mux.HandleFunc("POST /api/chirps", apiCfg.NewChirpPost)
	assetsFS := http.FileServer(http.Dir("./assets"))
	assetsPrefix := http.StripPrefix("/assets/", assetsFS)
	rootFS := http.FileServer(http.Dir("."))
	rootPrefix := http.StripPrefix("/app/", rootFS)
	mux.Handle("/assets/", assetsPrefix)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(rootPrefix))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.metricsResetHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
}
