package views

import (
	"errors"
	"example/skeleton/middleware"
	"example/skeleton/models"
	"example/skeleton/services"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestPostsListRoute(t *testing.T) {
	middleware.InitializeSessionStore()
	models.RunInit()
	posts := models.CreateDummyPosts()
	app := createTestRouter()

	req := httptest.NewRequest("GET", "/posts", nil)
	resp := executeRequest(req, app)
	assert.Equal(t, resp.Code, http.StatusOK)
	assert.Contains(t, resp.Body.String(), posts[0].Title)
	assert.Contains(t, resp.Body.String(), posts[1].Title)
	assert.Contains(t, resp.Body.String(), posts[2].Title)

	tmp := fmt.Sprintf("/posts/%s", posts[2].Uuid)
	assert.Contains(t, resp.Body.String(), tmp, "Should contain the UUID of post[2]")
}

func TestNewPostRoutes(t *testing.T) {
	middleware.InitializeSessionStore()
	models.RunInit()
	models.CreateDummyPosts()
	app := createTestRouter()
	cookie := start_session("vincent@saulys.me", "password", app)
	assert.NotEqual(t, "", cookie)

	reqN := httptest.NewRequest("GET", "/posts/new", nil)
	respN := executeRequest(reqN, app, cookie)
	assert.Equal(t, http.StatusOK, respN.Code)
	doc, err := goquery.NewDocumentFromReader(respN.Body)
	assert.Nil(t, err, "Should have no errors from parsing response in goquery")
	assert.Equal(t, doc.Find("input[name='title']").Size(), 1)
	assert.Equal(t, doc.Find("textarea[name='content']").Size(), 1)
	assert.Equal(t, doc.Find("input[name='is_public']").Size(), 1)

	req := httptest.NewRequest(
		http.MethodPost,
		"/posts/new",
		convert_to_post(map[string]string{
			"title":     "TestNewPostRoute",
			"content":   "TestNewPostRouteContent",
			"is_public": "1",
		}))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp := executeRequest(req, app, cookie).Result()
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode, "Not Getting a 'status see other' response")
	assert.Equal(t, "/posts", resp.Header.Get("Location"), "Did not get redirect to /posts")

	post := &models.Post{}
	models.Db.Where("title = ?", "TestNewPostRoute").First(&post)
	assert.Equal(t, post.Title, "TestNewPostRoute", "New Post could not be found in Database?")
	assert.Equal(t, post.Content, "TestNewPostRouteContent")
	assert.True(t, post.IsPublic, "Post is not public by default as we expected")

	// ---- attempt to get detail of _public_ post
	detailPostUrl := fmt.Sprintf("/posts/%s", post.Uuid)
	req = httptest.NewRequest("GET", detailPostUrl, nil)
	respDetailGet := executeRequest(req, app)
	assert.Equal(t, respDetailGet.Code, http.StatusOK)
	doc, err = goquery.NewDocumentFromReader(respDetailGet.Body)
	assert.Nil(t, err, "Should have no errors from parsing in goquery")
	assert.Contains(t, doc.Text(), post.Title)
	assert.Contains(t, doc.Text(), post.Content)

	// ---- Confirm Validation Errors if we send empty body
	req = httptest.NewRequest(
		http.MethodPost,
		"/posts/new",
		convert_to_post(map[string]string{}))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	respV := executeRequest(req, app, cookie)
	assert.Equal(t, http.StatusSeeOther, respV.Code)
	assert.Equal(t, "/posts/new", respV.Header().Get("Location"), "Did not get redirect to /posts/new after bad input")
	reqN = httptest.NewRequest("GET", "/posts/new", nil)
	reqN.Header.Add("Cookie", cookie) // combine cookies
	reqN.Header.Add("Cookie", respV.Result().Header.Get("set-cookie"))
	respN = executeRequest(reqN, app)
	assert.Equal(t, http.StatusOK, respN.Code)
	doc, err = goquery.NewDocumentFromReader(respN.Body)
	assert.Nil(t, err, "Should have no errors from parsing in goquery")
	assert.Contains(t, doc.Text(), "No name='title' passed", "Could not find \"No name='title' passed\" error in HTML response")

	// TODO: attempt to get _private_ post?
}

func TestUpdatePostRoutes(t *testing.T) {
	middleware.InitializeSessionStore()
	models.RunInit()
	app := createTestRouter()
	cookie := start_session("vincent@saulys.me", "password", app)
	assert.NotEqual(t, "", cookie)

	// Retrieve default user
	user := &models.User{}
	result := models.Db.Where("name = ?", "Vincent").Take(&user) // First adds "ORDER BY id"
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic("Did we not initiate user yet?")
	}

	// Get first private post for User
	allPostsForUser := services.RetrieveAllPostsForUser(*user)
	assert.GreaterOrEqual(t, len(*allPostsForUser), 1, "We don't have at least one post")
	idx := slices.IndexFunc(*allPostsForUser, func(p models.Post) bool { return !p.IsPublic })
	assert.NotEqual(t, -1, idx, "Cannot find a private post for user")
	dummyPost := (*allPostsForUser)[idx]
	assert.NotEqual(t, "", dummyPost.Uuid)

	// ---- GET update post via routes
	req := httptest.NewRequest("GET", fmt.Sprintf("/posts/%s/edit", dummyPost.Uuid), nil)
	resp := executeRequest(req, app, cookie)
	assert.Equal(t, http.StatusOK, resp.Code)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	raw_string, _ := doc.Html()
	assert.Nil(t, err, "Should have no errors from parsing response in goquery")
	assert.Equal(t, 1, doc.Find("input[name='title']").Size(), raw_string)
	assert.Equal(t, 1, doc.Find("textarea[name='content']").Size(), raw_string)
	assert.Equal(t, 1, doc.Find("input[name='is_public']").Size(), raw_string)
	anotherCookie := resp.Result().Header.Get("set-cookie") // Needed for checking errors

	// ---- POST update post via routes
	req = httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/posts/%s/edit", dummyPost.Uuid),
		convert_to_post(map[string]string{
			"title":     "TestUpdatePostRoute",
			"content":   "TestUpdatePostRouteContent",
			"is_public": "1",
		}))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp = executeRequest(req, app, cookie, anotherCookie)
	cookie3 := resp.Result().Header.Get("set-cookie")
	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Equal(t, "/posts", resp.Result().Header.Get("Location"), "Did not get redirect to /posts")
	req = httptest.NewRequest("GET", "/posts", nil)
	resp = executeRequest(req, app, cookie, anotherCookie, cookie3)
	assert.Contains(t, resp.Body.String(), "Successfully Updated!")
	assert.Contains(t, resp.Body.String(), "TestUpdatePostRoute")

	// TODO: check that errors come up in views

	// ---- Check that this post is updated as expected
	req = httptest.NewRequest("GET", fmt.Sprintf("/posts/%s", dummyPost.Uuid), nil)
	resp = executeRequest(req, app)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "TestUpdatePostRoute")
	assert.Contains(t, resp.Body.String(), "TestUpdatePostRouteContent")
	// IMPLICIT: IsPublic is tested because we view this without cookies
}
