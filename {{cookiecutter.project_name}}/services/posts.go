package services

import (
	"errors"
	"example/skeleton/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	// "github.com/go-playground/validator/v10"
)

func RetrieveAllPublicPosts() *[]models.Post {
	posts := &[]models.Post{}
	models.Db.Where("is_public = ?", true).Find(posts)
	return posts
}

func RetrieveAllPostsForUser(user models.User) *[]models.Post {
	posts := &[]models.Post{}
	models.Db.Where("is_public = ?", true).Or("author_id = ?", user.ID).Find(posts)
	return posts
}

func RetrieveDetailPost(postUuid string) (bool, *models.Post) {
	post := &models.Post{}
	result := models.Db.Where("uuid = ?", postUuid).Preload("Author").First(&post)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound), post
}

// returns bool for whether it was successfully created and map[string]string
// for errors where the key is the problematic field
func CreateNewPost(user *models.User, formData FormData) (bool, map[string]string) {
	title, ok := parseForm(formData, "title")
	if !ok {
		return false, map[string]string{"error": "No name='title' passed"}
	}
	content, ok := parseForm(formData, "content")
	if !ok {
		return false, map[string]string{"error": "No name='content' passed"}
	}
	_, is_public_parsed := parseForm(formData, "is_public")

	post := &models.Post{
		Title:    title,
		Content:  content,
		IsPublic: is_public_parsed,
		AuthorID: user.ID,
	}
	models.Db.Save(post)

	return true, map[string]string{}
}

func UpdatePostByUuid(formData FormData, postUuid string) (bool, *models.Post) {
	found, post := RetrieveDetailPost(postUuid)
	if !found {
		log.Errorf("Could not find post with uuid=%s", postUuid)
		return false, &models.Post{}
	}

	title, ok := formData["title"]
	if ok && len(title) != 0 {
		post.Title = title[0]
	}

	content, ok := formData["content"]
	log.Debug(!ok, len(content), content)
	if ok && len(content) != 0 {
		post.Content = content[0]
	}
	var is_public_parsed bool
	_, exists := formData["is_public"]
	if exists {
		is_public_parsed = true
	} else {
		is_public_parsed = false
	}
	if ok {
		log.Error("Public?  ", is_public_parsed)
		post.IsPublic = is_public_parsed
	}

	models.Db.Save(&post)

	return true, post
}
