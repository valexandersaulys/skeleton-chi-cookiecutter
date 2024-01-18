package middleware

import (
	"{{cookiecutter.project_name}}/config"
	"{{cookiecutter.project_name}}/models"
	"{{cookiecutter.project_name}}/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRequiredMiddleware(t *testing.T) {
	config.RunInit()
	InitializeSessionStore() // from middleware
	models.RunInit()
	models.CreateDummyPosts()
	*config.Environment = "LOCAL" // change to force AuthRequired to not fail

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { /*Add optional tests here*/ })
	handlerToTest := AuthRequired(nextHandler)

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	handlerToTest.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusForbidden, resp.Code)

	session, err := Store.Get(req, "auth")
	assert.Nil(t, err, "Error initializing the session")
	user := services.GetDefaultUser()
	assert.NotNil(t, user)
	session.Values["user"] = user
	err = session.Save(req, resp)
	assert.Nil(t, err, "Error saving the session")

	resp = httptest.NewRecorder()
	handlerToTest.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	*config.Environment = "TESTING"
}
