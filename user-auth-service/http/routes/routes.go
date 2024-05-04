package routes

import (
	"net/http"
	"ua-service/controller"
)

var RegisterAuthRoutes = func(router *http.ServeMux, authController *controller.AuthController) {
	router.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", router))

	router.HandleFunc("POST /register", authController.Register)
	router.HandleFunc("POST /login", authController.Login)
	router.HandleFunc("POST /validate", authController.ValidateToken)
}
