package utils

import (
	"encoding/json"
	"net/http"
	"ua-service/utils/log"

	"golang.org/x/crypto/bcrypt"
)

func WriteResponse(respPayload any, sc int, w http.ResponseWriter) {
	respJson, err := json.Marshal(respPayload)
	if err != nil {
		log.Error("Error {%v} marshaling response: %v", err, respPayload)
		http.Error(w, "Error converting response to JSON", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sc)
	w.Write(respJson)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func HashPassword(plainPassword string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func ComparePassword(hashedPass, plainPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass))
	return err == nil
}
