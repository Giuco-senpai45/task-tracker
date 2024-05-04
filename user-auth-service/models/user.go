package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func FindUserByEmail(email string) (User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func GetAllUsers() ([]User, error) {
	var users []User
	result := db.Find(&users)
	return users, result.Error
}

func GetUserById(id string) (User, error) {
	var user User
	result := db.First(&user, id)
	return user, result.Error
}

func DeleteUserById(id string) *gorm.DB {
	db.Where("ID = ?", id).Delete(&User{})
	return db
}

func (u *User) AddUser() *gorm.DB {
	db.Create(u)
	return db
}
