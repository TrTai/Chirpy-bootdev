package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/TrTai/Chirpy-bootdev/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	platform       string
	dbQueries      *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(rw, req)
	})
}

func (cfg *apiConfig) metricsHandler(rw http.ResponseWriter, req *http.Request) {
	// rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(http.StatusOK)
	// hitCountText := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	hitCountText := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	respBytes := ([]byte)(hitCountText)
	rw.Write(respBytes)
}

func (cfg *apiConfig) metricsResetHandler(rw http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(rw, 403, "")
		return
	}
	cfg.fileserverHits.Store(0)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx := req.Context()
	err := cfg.dbQueries.ResetUsers(ctx)
	if err != nil {
		errmsg := fmt.Sprintf("Error resetting user table: %s", err)
		respondWithError(rw, 500, errmsg)
	}
	rw.WriteHeader(http.StatusOK)
	bodyText := fmt.Sprint("Hits and DB Reset Successfully")
	respBytes := ([]byte)(bodyText)
	rw.Write(respBytes)

}

func (cfg *apiConfig) createUser(rw http.ResponseWriter, req *http.Request) {
	type UserJson struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	userReq := UserJson{}
	err := decoder.Decode(&userReq)
	if err != nil {
		errmsg := fmt.Sprintf("Error decoding user request: %s", err)
		respondWithError(rw, 500, errmsg)
	}
	ctx := req.Context()
	userRow, err := cfg.dbQueries.CreateUser(ctx, userReq.Email)
	if err != nil {
		errmsg := fmt.Sprintf("Error creating user: %s", err)
		respondWithError(rw, 500, errmsg)
	}
	respBody := UserJson(userRow)
	respondWithJSON(rw, 201, respBody)
}

func (cfg *apiConfig) NewChirpPost(rw http.ResponseWriter, req *http.Request) {
	type PostJson struct {
		ID        uuid.UUID     `json:"id"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
		Body      string        `json:"body"`
		UserID    uuid.NullUUID `json:"user_id"`
	}
	decoder := json.NewDecoder(req.Body)
	postReq := PostJson{}
	err := decoder.Decode(&postReq)
	if err != nil {
		errmsg := fmt.Sprintf("Error decoding post request: %s", err)
		respondWithError(rw, 500, errmsg)
	}
	postText, err := ValidateChirp(postReq.Body)
	if err != nil {
		errmsg := fmt.Sprintf("Chirp Size Error: %s", err)
		respondWithError(rw, 400, errmsg)
	}
	ctx := req.Context()
	postParams := database.CreatePostParams{
		postText,
		postReq.UserID,
	}
	postRow, err := cfg.dbQueries.CreatePost(ctx, postParams)
	if err != nil {
		errmsg := fmt.Sprintf("Error creating post: %s", err)
		respondWithError(rw, 500, errmsg)
	}
	respBody := PostJson(postRow)
	respondWithJSON(rw, 201, respBody)

}

func (cfg *apiConfig) GetAllChirps(rw http.ResponseWriter, req *http.Request) {

}
