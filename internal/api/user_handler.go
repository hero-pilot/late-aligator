package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/hero-pilot/late-aligator/internal/store"
	"github.com/hero-pilot/late-aligator/internal/utils"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 4 {
		return errors.New("username must be at least 4 characters")
	}
	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 chars")
	}
	if !usernameRegex.MatchString(req.Username) {
		return errors.New("username may only contain letters, numbers, and underscores")
	}
	if req.Email == "" {
		return errors.New("Email is required")
	}
	if len(req.Email) > 254 {
		return errors.New("email is too long")
	}
	if !emailRegex.MatchString(req.Email) {
		return errors.New("Invalid email format")
	}
	if req.Password == "" {
		return errors.New("Password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if len(req.Password) > 72 {
		return errors.New("password must be at most 72 characters")
	}

	return nil

}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: Decoding register request: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
	}
	err = uh.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user := &store.User{
		Username: req.Username,
		Email: req.Email,
	}
	if req.Bio != ""{
		user.Bio = req.Bio
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil{
		uh.logger.Printf("ERROR: Hashing Password: %v",err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Inernal server error"})
		return 
	}
	err = uh.userStore.CreateUser(user)
	if err !=nil{
		uh.logger.Printf("ERROR: Registering User: %v",err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Inernal server error"})
		return 
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"user": user})

}
