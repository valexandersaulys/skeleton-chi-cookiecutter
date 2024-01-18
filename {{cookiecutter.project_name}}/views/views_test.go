package views

import (
	"{{cookiecutter.project_name}}/config"
	"{{cookiecutter.project_name}}/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func executeRequest(req *http.Request, httpRouter *chi.Mux, cookies ...string) *httptest.ResponseRecorder {
	for _, cookie := range cookies {
		req.Header.Add("Cookie", cookie)
	}
	rr := httptest.NewRecorder()
	httpRouter.ServeHTTP(rr, req)
	return rr
}

// Return the set-cookie header value for starting a logged in session
func start_session(email string, password string, httpRouter *chi.Mux) string {
	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		convert_to_post(map[string]string{
			"email":    email,
			"password": password,
		}))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp := executeRequest(req, httpRouter)
	cookie := resp.Result().Header.Get("set-cookie")
	if resp.Code != http.StatusSeeOther {
		panic("Could not successfully log in via views.start_session!")
	}
	return cookie
}

func createTestRouter() *chi.Mux {
	r := chi.NewRouter()
	AddRoutes(r)
	return r
}

func convert_to_post(formData map[string]string) *strings.Reader {
	toRet := url.Values{}
	for k, v := range formData {
		toRet.Add(k, v)
	}
	// fmt.Println(strings.NewReader(toRet.Encode()))
	return strings.NewReader(toRet.Encode())
}

func TestMain(m *testing.M) {
	// Write code here to run before tests
	config.RunInit()

	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}

func TestMissing(t *testing.T) {
	middleware.InitializeSessionStore()
	app := createTestRouter()

	req := httptest.NewRequest("GET", "/asdfasdfasdf", nil)
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusNotFound, resp.Code, resp.Body.String())
}

func TestMethodNotAllowed(t *testing.T) {
	middleware.InitializeSessionStore()
	app := createTestRouter()

	req := httptest.NewRequest("OPTIONS", "/", nil)
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code, resp.Body.String())
}

func TestStaticFileServer(t *testing.T) {
	middleware.InitializeSessionStore()
	app := createTestRouter()

	req := httptest.NewRequest("GET", "/public/placeholder.js", nil)
	resp := executeRequest(req, app)
	assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
}
