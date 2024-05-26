package services

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"{{cookiecutter.project_name}}/models"
)

// DEBUG: to be removed later
func GetDefaultUser(ctx context.Context) *models.User {
	db, err := models.GetDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	user := &models.User{}
	result := db.Where("name = ?", "Vincent").Take(&user) // First adds "ORDER BY id"
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic("Did we not initiate user yet?")
	}
	return user
}

func AuthenticateUser(ctx context.Context, formData FormData) (bool, ValidationErrors, *models.User) {
	db, err := models.GetDb(ctx)
	if err != nil {
		log.Fatal(err)
	}

	email, ok := parseForm(formData, "email")
	if !ok {
		return false, ValidationErrors{"error": "No Email passed"}, &models.User{}
	}
	password, ok := parseForm(formData, "password")
	if !ok {
		return false, map[string]string{"error": "No Password passed"}, &models.User{}
	}

	user := &models.User{}
	result := db.Where("email = ?", email).Take(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, map[string]string{"error": "Could not login. Check your email or password."}, &models.User{}
	}
	if user.VerifyPassword(password) {
		return true, map[string]string{}, user
	}

	return false, map[string]string{"error": "Cannot authenticate password"}, &models.User{}
}
