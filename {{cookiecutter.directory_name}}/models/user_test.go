package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreationAndRetrieval(t *testing.T) {
	db, err := GetDbWithNoContext()
	if err != nil {
		panic(err)
	}

	user := &User{
		Name:     "Vincent",
		Email:    "vincent@example.com",
		Password: "password",
	}
	db.Create(user)

	assert.Len(t, user.Uuid, 36, "Uuid is not a length of 36 as expected")
	assert.Regexp(t,
		"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
		user.Uuid,
		"user.Uuid is _not_ a uuid")
	assert.Equal(t, user.Name, "Vincent")
	assert.Equal(t, user.Email, "vincent@example.com")
	assert.NotEqual(t, user.Password, "password", "Password is not being hashed")
	assert.False(t, user.VerifyPassword(user.Password),
		"Hashed password should not pass VerifyPassword")
	assert.False(t, user.VerifyPassword("egah"),
		"User.VerifyPassword should not work on incorrect passwords")
	assert.True(t, user.VerifyPassword("password"), "correct password should verify")
}
