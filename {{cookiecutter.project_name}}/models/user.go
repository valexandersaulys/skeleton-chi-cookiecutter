package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uuid     string `gorm:"type:CHAR(35);column:uuid;index;not null"`
	Name     string
	Email    string
	Password string
	Posts    []Post `gorm:"foreignKey:author_id"`
}

func (user *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}

// ---------------------------------------- Hooks

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)
	return
}

// Instantiate uuid column in User
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Uuid = uuid.New().String()
	return
}
