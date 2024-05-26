package models

import (
	_ "ariga.io/atlas-provider-gorm/gormschema"
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"{{cookiecutter.project_name}}/config"
)

func GetDb(ctx context.Context) (*gorm.DB, error) {
	gormDb, err := GetDbWithNoContext()
	if err != nil {
		return nil, err
	}

	return gormDb.WithContext(ctx), nil
}

func GetDbWithNoContext() (*gorm.DB, error) {
	var gormDb *gorm.DB

	var sqlDb gorm.Dialector
	switch *config.Environment {
	case "TESTING":
		sqlDb = sqlite.Open("file::memory:?cache=shared&_journal_mode=WAL")
	case "LOCAL":
		sqlDb = sqlite.Open("/tmp/{{cookiecutter.application_name}}.db?_journal_mode=WAL")
	case "DEV":
		sqlDb = sqlite.Open("/data/db/{{cookiecutter.application_name}}.db?_journal_mode=WAL")
	case "PROD":
		sqlDb = sqlite.Open("/data/db/{{cookiecutter.application_name}}.db?_journal_mode=WAL")
	default:
		sqlDb = sqlite.Open("/tmp/{{cookiecutter.application_name}}.db")
	}

	var err error
	gormDb, err = gorm.Open(sqlDb, &gorm.Config{
		Logger: CustomGormLogger(true)(),
	})
	if err != nil {
		return nil, err
	}

	return gormDb, nil
}

func RunInit() {
	gormDb, _ := GetDbWithNoContext()

	underlyingDb, err := gormDb.DB()
	if err != nil {
		panic(err)
	}

	if err := underlyingDb.Ping(); err != nil {
		panic(err)
	}

	underlyingDb.SetMaxIdleConns(*config.MaxIdleDbConnections)
	underlyingDb.SetMaxOpenConns(*config.MaxOpenDbConnections)

	if *config.RunMigrations || *config.Environment == "TESTING" {
		gormDb.AutoMigrate(&User{}, &Post{})
	}
}

func ClearAll(db *gorm.DB) {
	db.Where("id > ?", 0).Delete(&User{})
	db.Where("id > ?", 0).Delete(&Post{})
}

func CreateDummyPosts(db *gorm.DB) []*Post {
	user := &User{
		Name:     "Vincent",
		Email:    "vincent@example.com",
		Password: "password",
	}
	db.Create(user)

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
	db.Create(posts)
	return posts
}
