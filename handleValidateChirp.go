package main

import (
	"encoding/json"
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpJSON struct {
		Body string `json:"body"`
	}
	type validChirp struct {
		Valid bool `json:"valid"`
	}
	chirpBody := chirpJSON{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirpBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	const chirpMaxLength = 140
	if len(chirpBody.Body) > chirpMaxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	respondWithJSON(w, 200, validChirp{Valid: true})
}
