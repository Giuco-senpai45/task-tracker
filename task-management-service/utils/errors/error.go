package errors

import (
	"errors"
	"tm-service/utils/log"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials  = errors.New("invalid login credentials")
	ErrInRequestMarshaling = errors.New("invalid/bad request paramenters")
	ErrDuplicateEmail      = errors.New("email already exists")
	ErrMalformedToken      = errors.New("malformed jwt token")
	ErrNonExistingService  = errors.New("service does not exist")
)

func Error(e error) {
	log.Error(e.Error())
	// panic(e)
}

func DBErrorCheck(db *gorm.DB) {
	if err := db.Error; err != nil {
		Error(err)
	}
}

func ErrorCheck(e error) {
	if e != nil {
		Error(e)
	}
}
