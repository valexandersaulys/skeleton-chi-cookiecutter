package services

import (
	"{{cookiecutter.project_name}}/models"
	"github.com/stretchr/testify/assert"
	// "sort"
	"fmt"
	"slices"
	"testing"
)

func TestRetrieveAndUpdatePosts(t *testing.T) {
	models.RunInit()
	dummyPosts := models.CreateDummyPosts() // author.name == "Vincent"

	privatePost := dummyPosts[0]
	privatePost.IsPublic = false
	models.Db.Save(privatePost)

	publicPosts := RetrieveAllPublicPosts()

	for _, dummyPost := range dummyPosts[1:] {
		if !dummyPost.IsPublic {
			continue
		}
		idx := slices.IndexFunc(*publicPosts, func(p models.Post) bool { return p.ID == dummyPost.ID })
		if idx == -1 {
			assert.False(t, true,
				fmt.Sprintf("Could not find public dummy post: '%s'", dummyPost.Title),
			)
		}
	}
	idx := slices.IndexFunc(*publicPosts, func(p models.Post) bool { return p.ID == privatePost.ID })
	if idx != -1 {
		assert.False(t, true, "Should _not_ be able to find private post in public posts")
	}

	yay, retrievedPost := RetrieveDetailPost(privatePost.Uuid)
	assert.True(t, yay, "Could not successfully retrieve details for _private_ post")
	assert.Equal(t, retrievedPost.Uuid, privatePost.Uuid)
	assert.Equal(t, retrievedPost.ID, privatePost.ID)

	nay, _ := RetrieveDetailPost("3a682ebf-71d7-4d82-98ae-d626539243fe")
	assert.False(t, nay, "_Did_ successfully retrieve post for non-existent post")

	var user models.User
	models.Db.Where("name = ?", "Vincent").First(&user)
	allPostsForUser := RetrieveAllPostsForUser(user)
	assert.Equal(t, 4, len(*allPostsForUser))

	var updatedPost *models.Post
	yay, updatedPost = UpdatePostByUuid(map[string][]string{
		"title":     []string{"My Updated Title"},
		"content":   []string{"bicycle rights ethical raw denim ascot same"},
		"is_public": []string{"1"},
	}, privatePost.Uuid)
	assert.True(t, yay, fmt.Sprintf("Did _not_ successfully update the post with uuid=%s", privatePost.Uuid))
	assert.NotNil(t, updatedPost)
	assert.Equal(t, privatePost.Uuid, updatedPost.Uuid)
	assert.Equal(t, privatePost.ID, updatedPost.ID)
	assert.NotEqual(t, privatePost.Title, updatedPost.Title)
	assert.NotEqual(t, privatePost.Content, updatedPost.Content)
	assert.Equal(t, "My Updated Title", updatedPost.Title)
	assert.Equal(t, "bicycle rights ethical raw denim ascot same", updatedPost.Content)
	assert.True(t, updatedPost.IsPublic)
}

func TestCreateNewPost(t *testing.T) {
	models.RunInit()

	user := &models.User{
		Name:     "Vincent",
		Email:    "vincent@example.com",
		Password: "password",
	}
	models.Db.Create(user)

	yay, validationBits := CreateNewPost(user, map[string][]string{
		"title":     []string{"TestCreateNewPost"},
		"content":   []string{"TestCreateNewPostContent"},
		"is_public": []string{"1"},
	})
	assert.True(t, yay, "Did not successfully create post after title, content, and is_public were active")
	assert.Equal(t, validationBits, map[string]string{})

	yay, validationBits = CreateNewPost(user, map[string][]string{
		"title":     []string{""},
		"content":   []string{"TestCreateNewPostContent"},
		"is_public": []string{"1"},
	})
	assert.False(t, yay, "_Did_ successfully create post despite missing title")
	assert.NotEqual(t, validationBits, map[string]string{})

	yay, validationBits = CreateNewPost(user, map[string][]string{
		"title":     []string{"TestCreateNewPost"},
		"content":   []string{""},
		"is_public": []string{"1"},
	})
	assert.False(t, yay, "_Did_ successfully create post despite missing content")
	assert.NotEqual(t, validationBits, map[string]string{})

	// Note: you _can_ create without is_public present -- r.Form
	// will not have an input[type=checkbox] populated if its not
	// at all checked.

	yay, validationBits = CreateNewPost(user, map[string][]string{})
	assert.False(t, yay)
	assert.NotEqual(t, validationBits, map[string]string{})
}
