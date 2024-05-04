package responses

import "ua-service/models"

type UserResponse struct {
	User models.User `json:"user"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
