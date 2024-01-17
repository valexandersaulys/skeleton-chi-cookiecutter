package app

import (
	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	"example/skeleton/config"
	"example/skeleton/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Write code here to run before tests
	config.RunInit()

	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}

func executeRequest(req *http.Request, httpRouter *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	httpRouter.ServeHTTP(rr, req)
	return rr
}

func TestHeartbeatPing(t *testing.T) {
	assert.Nil(t, middleware.Store, "Store should be Nil before initializing session")
	middleware.InitializeSessionStore()
	assert.NotNil(t, middleware.Store, "Store should _not_ be Nil _after_ initializing session")
	app := CreateApp()
	req, _ := http.NewRequest("GET", "/ping", nil)
	resp := executeRequest(req, app)
	assert.Equal(t, resp.Code, http.StatusOK)
}
