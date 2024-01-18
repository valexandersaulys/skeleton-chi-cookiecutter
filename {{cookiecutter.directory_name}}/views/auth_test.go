package views

import (
	"{{cookiecutter.project_name}}/middleware"
	"{{cookiecutter.project_name}}/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginUserRoute(t *testing.T) {
	middleware.InitializeSessionStore()
	models.RunInit()
	models.CreateDummyPosts()
	app := createTestRouter()

	// ---- Test /login
	req := httptest.NewRequest("GET", "/login", nil)
	resp := executeRequest(req, app)
	assert.Equal(t, resp.Code, http.StatusOK)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	assert.Nil(t, err, "Should have no errors from parsing response in goquery")
	assert.Equal(t, doc.Find("EGAHASDF").Size(), 0, "Something is very broken in goquery!")
	assert.Equal(t, doc.Find("html").Size(), 1)
	assert.Equal(t, doc.Find("body").Size(), 1)
	assert.Equal(t, doc.Find("input[name='email']").Size(), 1)
	assert.Equal(t, doc.Find("input[name='password']").Size(), 1)

	// ---- Succesfully attempt to Login
	req = httptest.NewRequest(
		http.MethodPost,
		"/login",
		convert_to_post(map[string]string{
			"email":    "vincent@example.com",
			"password": "password",
		}))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp = executeRequest(req, app)
	authCookie := resp.Result().Header.Get("set-cookie")
	assert.Equal(t, http.StatusSeeOther, resp.Result().StatusCode, "Not Getting a 'status see other' response")
	assert.Equal(t, "/posts", resp.Result().Header.Get("Location"), "Did not get redirect to /posts")
	assert.NotEqual(t, "", authCookie)

	// ---- Attempt redirect after logging in
	req = httptest.NewRequest("GET", "/posts", nil)
	resp = executeRequest(req, app, authCookie)
	assert.Equal(t, http.StatusOK, resp.Code)
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	assert.Nil(t, err, "Should have no errors from parsing response in goquery")
	assert.Equal(t, doc.Find("html").Size(), 1)
	assert.Equal(t, doc.Find("h2").Size(), 4, "Could not find the 3 h2 tags")

	// TODO: Attempt log in again to get redirect to /posts

	req = httptest.NewRequest(
		http.MethodPost,
		"/login",
		convert_to_post(map[string]string{
			"email":    "egah@egah.com",
			"password": "wrong-password",
		}))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp = executeRequest(req, app)
	assert.Equal(t, http.StatusSeeOther, resp.Result().StatusCode, "")
	assert.Equal(t, "/login", resp.Result().Header.Get("Location"), "Did not get redirect to /login after bad authentication")
	// TODO: test that this returns an error on the html

	// ---- Test that we can log out
	req = httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp = executeRequest(req, app, authCookie)
	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Equal(t, "/login", resp.Result().Header.Get("Location"), "Did not get redirect to /login after logout")
	req = httptest.NewRequest("GET", "/login", nil)
	resp = executeRequest(req, app, resp.Result().Header.Values("set-cookie")...)
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	assert.Contains(t, resp.Body.String(), "Successfully Logged Out")
}
