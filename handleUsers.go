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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	requestRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting authorization", err)
		return
	}
	refreshTokenTable, err := cfg.dbQueries.GetUserWithToken(r.Context(), requestRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no refresh token", err)
		return
	}
	if refreshTokenTable.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token revoked", err)
		return
	}
	if refreshTokenTable.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "token expired", err)
		return
	}
	accessToken, err := auth.MakeJWT(refreshTokenTable.UserID, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making JWT", err)
		return
	}
	accessTokenResponse := AccessTokenResponse{Token: accessToken}
	respondWithJSON(w, http.StatusOK, accessTokenResponse)
}

type AccessTokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	authTokenBearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting bearer", err)
		return
	}
	err = cfg.dbQueries.RevokeToken(r.Context(), authTokenBearer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking token", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}

type UpdateUserRequest struct {
	Email    string
	Password string
}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	accesToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid token", err)
		return
	}
	userID, err := auth.ValidateJWT(accesToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	updateUserRequest := UpdateUserRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updateUserRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding user update request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(updateUserRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	updateParams := database.UpdateUserParams{
		Email:          updateUserRequest.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	}
	updatedUserInfo, err := cfg.dbQueries.UpdateUser(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user info", err)
		return
	}
	respondWithJSON(w, http.StatusOK, UpdateUserResponse{
		Email:     updatedUserInfo.Email,
		CreatedAt: updatedUserInfo.CreatedAt,
		UpdatedAt: updatedUserInfo.UpdatedAt,
		ID:        updatedUserInfo.ID,
		IsChirpyRed: updatedUserInfo.,
	})
}

type UpdateParams struct {
	Email          string
	HashedPassword string
	UserID         uuid.UUID
}

type UpdateUserResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type UpgradeUserRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handleUpgradeUser(w http.ResponseWriter, r *http.Request) {
	upgradeUserRequest := UpgradeUserRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&upgradeUserRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request", err)
		return
	}
	if upgradeUserRequest.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	err = cfg.dbQueries.UpgradeUser(r.Context(), upgradeUserRequest.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
