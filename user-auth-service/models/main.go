package models

import "gorm.io/gorm"

var db *gorm.DB

func New(dbPool *gorm.DB) Models {
	db = dbPool
	return Models{}
}

type Models struct {
	User User
}
