package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

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
		cleanBody, err := ValidateChirp(params.Body)
		if err != nil {
			return
		}
		respBody := ValidText{
			ValidBody:   true,
			CleanedBody: cleanBody,
		}
		respondWithJSON(rw, 200, respBody)
	}

}

func ValidateChirp(post string) (string, error) {
	if len(post) > 140 || len(post) < 1 {
		return "", errors.New("Invalid String Length")
	}
	slicePost := strings.Split(post, " ")
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range slicePost {
		if slices.Contains(badWords, strings.ToLower(word)) {
			slicePost[i] = "****"
		}
	}
	return strings.Join(slicePost, " "), nil

}
