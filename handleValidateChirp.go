package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpJSON struct {
		Body string `json:"body"`
	}
	type cleanedChirp struct {
		CleanedChirp string `json:"cleaned_body"`
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
	respondWithJSON(w, 200, cleanedChirp{CleanedChirp: cleanChirp(chirpBody.Body)})
}

func cleanChirp(chirp string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	wordList := strings.Split(chirp, " ")
	for i, word := range wordList {
		loweredWord := strings.ToLower(word)
		if _, found := badWords[loweredWord]; found {
			wordList[i] = "****"
		}
	}
	cleanedChirp := strings.Join(wordList, " ")
	return cleanedChirp
}
