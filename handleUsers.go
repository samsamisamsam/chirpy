package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/samsamisamsam/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Email struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var email Email
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deconding request", err)
		return
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), email.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating the user", err)
		return
	}
	userJSON := userToJSON(user)
	respondWithJSON(w, http.StatusCreated, userJSON)
}

func (cfg *apiConfig) deleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Endpoint only accessible in developpement", nil)
		return
	}
	cfg.dbQueries.DeleteAllUsers(r.Context())
	respondWithJSON(w, http.StatusOK, nil)
}

func userToJSON(user database.User) User {
	convertedUser := User{}
	convertedUser.ID = user.ID
	convertedUser.CreatedAt = user.CreatedAt
	convertedUser.UpdatedAt = user.UpdatedAt
	convertedUser.Email = user.Email
	return convertedUser
}
