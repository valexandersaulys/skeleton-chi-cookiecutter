package middleware

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"{{cookiecutter.project_name}}/config"
	"{{cookiecutter.project_name}}/models"
	"{{cookiecutter.project_name}}/services"
)

func TestAuthRequiredMiddleware(t *testing.T) {
	config.RunInit()
	InitializeSessionStore() // from middleware
	models.RunInit()

	db, err := models.GetDbWithNoContext()
	if err != nil {
		panic(err)
	}
	models.CreateDummyPosts(db)
	*config.Environment = "LOCAL" // change to force AuthRequired to not fail

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { /*Add optional tests here*/ })
	handlerToTest := AuthRequired(nextHandler)

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	handlerToTest.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusForbidden, resp.Code)

	session, err := Store.Get(req, "auth")
	assert.Nil(t, err, "Error initializing the session")
	user := services.GetDefaultUser(context.TODO())
	assert.NotNil(t, user)
	session.Values["user"] = user
	err = session.Save(req, resp)
	assert.Nil(t, err, "Error saving the session")

	resp = httptest.NewRecorder()
	handlerToTest.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	*config.Environment = "TESTING"
}
