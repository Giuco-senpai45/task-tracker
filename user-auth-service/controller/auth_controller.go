package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"ua-service/http/requests"
	"ua-service/jwt"
	"ua-service/service"
	"ua-service/utils"
	"ua-service/utils/log"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var createRequest requests.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		log.Error("Err decoding body : %v", err.Error())
		http.Error(w, "Error decoding the body", http.StatusInternalServerError)
		return
	}

	res, err := c.service.Register(createRequest)
	if err != nil {
		log.Error("Error registering user : %v", err.Error())
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(res, http.StatusOK, w)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest requests.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		log.Error("Err decoding body : %v", err.Error())
		http.Error(w, "Error decoding the body", http.StatusInternalServerError)
		return
	}

	loginToken, err := c.service.Login(loginRequest)
	if err != nil {
		log.Error("Error logging in user : %v", err.Error())
		http.Error(w, "Error logging in user", http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(loginToken, http.StatusOK, w)
}

func (c *AuthController) ValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		log.Error("No authorization header")
		http.Error(w, "No authorization header", http.StatusUnauthorized)
		return
	}

	token := strings.Split(authHeader, "Bearer ")
	if len(token) != 2 {
		log.Error("Invalid token")
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	wt, err := jwt.ValidateToken(token[1])
	if err != nil {
		log.Error("Incorrect token: %v", err.Error())
		http.Error(w, "Incorrect token", http.StatusUnauthorized)
		return
	}

	if !wt.Valid {
		log.Error("Invalid token")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userId, err := jwt.GetUserIdFromToken(token[1])
	if err != nil {
		log.Error("Error getting user id from token: %v", err.Error())
		http.Error(w, "Error getting user id from token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("X-User-Id", strconv.Itoa(userId))
	w.WriteHeader(http.StatusOK)
}
