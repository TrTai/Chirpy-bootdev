package main

import (
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	//"sync/atomic"
	"slices"
	"strings"
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
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.metricsResetHandler)
	mux.HandleFunc("POST /api/validate_chirp", chirpValidateHandler)
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

func chirpValidateHandler(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type ValidText struct {
		ValidBody   bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errmsg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(rw, 500, errmsg)
		return
	}

	if len(params.Body) > 140 {
		errMsg := fmt.Sprintf("Chirp is too long")
		respondWithError(rw, 400, errMsg)
	} else if len(params.Body) <= 140 {
		respBody := ValidText{
			ValidBody:   true,
			CleanedBody: censorChirp(params.Body),
		}
		respondWithJSON(rw, 200, respBody)
	}

}

func censorChirp(post string) string {
	slicePost := strings.Split(post, " ")
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range slicePost {
		if slices.Contains(badWords, strings.ToLower(word)) {
			slicePost[i] = "****"
		}
	}
	return strings.Join(slicePost, " ")

}
