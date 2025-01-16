package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samsamisamsam/chirpy/internal/auth"
	"github.com/samsamisamsam/chirpy/internal/database"
)

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("id")
	id, err := uuid.Parse(pathValue)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		respondWithError(w, 404, "Error converting id parameter to uuid", err)
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Error getting chirp from the database", err)
		return
	}
	respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting all chirps from the database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	id, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	type chirpRequest struct {
		Body string `json:"body"`
	}

	chirpReq := chirpRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	const chirpMaxLength = 140
	if len(chirpReq.Body) > chirpMaxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedChirpBody := cleanChirp(chirpReq.Body)
	params := database.CreateChirpParams{
		Body:   cleanedChirpBody,
		UserID: id,
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating the chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, database.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    id,
	})
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	Token     string    `json:"token"`
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

type chirpDB struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
