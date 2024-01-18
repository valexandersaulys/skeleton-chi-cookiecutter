package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostCreationAndRetrieval(t *testing.T) {
	user := &User{
		Name:     "Vincent",
		Email:    "vincent@example.com",
		Password: "password",
	}
	Db.Create(user)

	post := &Post{
		Title:   "My Title",
		Content: "shorts echo park. Kogi mustache pabst tumeric. Etsy photo booth",
		Author:  *user,
	}
	Db.Create(post)
	assert.Len(t, post.Uuid, 36, "Uuid is not a length of 36 as expected")
	assert.Regexp(t,
		"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
		post.Uuid,
		"post.Uuid is _not_ a uuid")
	assert.Equal(t, post.Title, "My Title")
	assert.Equal(t, post.Content, "shorts echo park. Kogi mustache pabst tumeric. Etsy photo booth")
	assert.Equal(t, post.Author.ID, user.ID)
	assert.False(t, post.IsPublic)
	post.IsPublic = true
	Db.Save(&post)

	aPost := &Post{}
	Db.Where("uuid = ?", post.Uuid).First(&aPost)
	assert.Equal(t, aPost.Uuid, post.Uuid)
	assert.Equal(t, aPost.ID, post.ID)
	assert.NotEqual(t, aPost.Author.ID, user.ID, "ID should not be retrieved unless we specify eager")
	assert.True(t, aPost.IsPublic,
		"IsPublic should have changed after updating and then committing to database")

	// Preload("<colname>") can be a non-existent "<colname>"
	aPost = &Post{}
	Db.Preload("Author").Where("uuid = ?", post.Uuid).First(&aPost)
	assert.Equal(t, aPost.Author.ID, user.ID, "User ID should be populated if we specify eager loading")
}
