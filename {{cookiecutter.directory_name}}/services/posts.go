package services

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"{{cookiecutter.project_name}}/models"
	// "github.com/go-playground/validator/v10"
)

func RetrieveAllPublicPosts(ctx context.Context) *[]models.Post {
	db, err := models.GetDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	posts := &[]models.Post{}
	db.Where("is_public = ?", true).Find(posts)
	return posts
}

func RetrieveAllPostsForUser(ctx context.Context, user models.User) *[]models.Post {
	db, err := models.GetDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	posts := &[]models.Post{}
	db.Where("is_public = ?", true).Or("author_id = ?", user.ID).Find(posts)
	return posts
}

func RetrieveDetailPost(ctx context.Context, postUuid string) (bool, *models.Post) {
	db, err := models.GetDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	post := &models.Post{}
	result := db.Where("uuid = ?", postUuid).Preload("Author").First(&post)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound), post
}

// returns bool for whether it was successfully created and map[string]string
// for errors where the key is the problematic field
func CreateNewPost(ctx context.Context, user *models.User, formData FormData) (bool, map[string]string) {
	db, err := models.GetDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
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
	db.Save(post)

	return true, map[string]string{}
}

func UpdatePostByUuid(ctx context.Context, formData FormData, postUuid string) (bool, *models.Post) {
	db, err := models.GetDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	found, post := RetrieveDetailPost(ctx, postUuid)
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

	db.Save(&post)

	return true, post
}
