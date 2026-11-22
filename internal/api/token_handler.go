package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/hero-pilot/late-aligator/internal/store"
	"github.com/hero-pilot/late-aligator/internal/tokens"
	"github.com/hero-pilot/late-aligator/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type CreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (th *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req CreateTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		th.logger.Printf("ERROR: CreateToken Request:%v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payloas"})
		return
	}

	user, err := th.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		th.logger.Printf("ERROR: GetUSerByUsername:%v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server Error"})
		return
	}

	matches, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		th.logger.Printf("ERROR: PasswordHash.Matches: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if !matches {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid credentials"})
		return
	}
	token, err := th.tokenStore.CreateNewToken(user.ID, time.Hour*24, tokens.ScopeAuth)
	if err !=nil {
		th.logger.Printf("ERROR: CreateNewToken: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"auth_token":token })
}
