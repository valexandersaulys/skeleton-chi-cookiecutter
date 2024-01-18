package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"{{cookiecutter.project_name}}/config"
)

var Db *gorm.DB

func RunInit() {
	var sqlDb gorm.Dialector
	switch *config.Environment {
	case "TESTING":
		sqlDb = sqlite.Open("file::memory:?cache=shared")
	case "LOCAL":
		sqlDb = sqlite.Open("/tmp/{{cookiecutter.application_name}}i.db")
	case "DEV":
		sqlDb = sqlite.Open("/tmp/{{cookiecutter.application_name}}i.db")
	case "PROD":
		sqlDb = sqlite.Open("/tmp/{{cookiecutter.application_name}}i.db")
	default:
		sqlDb = sqlite.Open("/tmp/{{cookiecutter.application_name}}i.db")
	}

	var err error
	Db, err = gorm.Open(sqlDb, &gorm.Config{
		Logger: CustomGormLogger(true)(),
	})
	if err != nil {
		panic(err)
	}
	// Silence is golden: no error? We're good

	if *config.RunMigrations {
		Db.AutoMigrate(&User{}, &Post{})
	}
}

func ClearAll() {
	Db.Where("id > ?", 0).Delete(&User{})
	Db.Where("id > ?", 0).Delete(&Post{})
}

func CreateDummyPosts() []*Post {
	user := &User{
		Name:     "Vincent",
		Email:    "vincent@example.com",
		Password: "password",
	}
	Db.Create(user)

	posts := []*Post{
		&Post{
			Title:    "My First Title",
			Content:  "shorts echo park. Kogi mustache pabst tumeric. Etsy photo booth",
			Author:   *user,
			IsPublic: true,
		},
		&Post{
			Title:    "My Second Title",
			Content:  "shorts echo park. Kogi mustache pabst tumeric. Etsy photo booth",
			Author:   *user,
			IsPublic: true,
		},
		&Post{
			Title:    "My Third Title",
			Content:  "shorts echo park. Kogi mustache pabst tumeric. Etsy photo booth",
			Author:   *user,
			IsPublic: true,
		},
		&Post{
			Title:    "My Private Fourth Title",
			Content:  "shorts echo park. Kogi mustache pabst tumeric. Etsy photo booth",
			Author:   *user,
			IsPublic: false,
		},
	}
	Db.Create(posts)
	return posts
}
