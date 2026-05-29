package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(rw http.ResponseWriter, code int, msg string) {
	type errorResp struct {
		ErrorBody string `json:"error"`
	}
	rw.WriteHeader(code)
	if msg != "" {
		log.Print(msg)
		respJson, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.Write(respJson)
	}

}

func respondWithJSON(rw http.ResponseWriter, code int, payload interface{}) {
	rw.WriteHeader(code)
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(dat)
}
