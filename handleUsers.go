package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/samsamisamsam/chirpy/internal/auth"
	"github.com/samsamisamsam/chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	var loginInfo = LoginInfo{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginInfo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding login request", err)
		return
	}
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), loginInfo.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting user info", err)
		return
	}
	err = auth.CheckPasswordHash(loginInfo.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "email/password combination not valid", err)
		return
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making refresh token", err)
		return
	}
	tokenParams := database.StoreRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}
	cfg.dbQueries.StoreRefreshToken(r.Context(), tokenParams)

	userWithoutPassword := UserWithoutPassword{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tokenString,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, userWithoutPassword)
}

type UserWithoutPassword struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type LoginInfo struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var loginInfo LoginInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginInfo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deconding request", err)
		return
	}

	var params database.CreateUserParams
	hashedPassword, err := auth.HashPassword(loginInfo.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	params.Email = loginInfo.Email
	params.HashedPassword = hashedPassword

	user, err := cfg.dbQueries.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating the user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handleDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Endpoint only accessible in developpement", nil)
		return
	}
	cfg.dbQueries.DeleteAllUsers(r.Context())
	cfg.fileserverHits.Store(0)
	respondWithJSON(w, http.StatusOK, nil)
}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {

}
