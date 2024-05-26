package main

import (
	"fmt"
	"github.com/go-chi/docgen"
	"net/http"
	"{{cookiecutter.project_name}}/app"
	"{{cookiecutter.project_name}}/config"
	"{{cookiecutter.project_name}}/middleware"
	"{{cookiecutter.project_name}}/models"
)

func init() {
	config.RunInit() // **Must Run First**
	models.RunInit()

	if *config.AddDummyModels {
		gormDb, err := models.GetDbWithNoContext()
		if err != nil {
			panic(err)
		}
		models.ClearAll(gormDb)
		models.CreateDummyPosts(gormDb)
	}

}

func main() {
	middleware.InitializeSessionStore()

	r := app.CreateApp()

	if *config.PrintRoutes {
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "<My Project Name>",
			Intro:       "My Project Intro",
		}))
		return
	}

	http.ListenAndServe(fmt.Sprintf(":%d", *config.RuntimePort), r) // adding go causes this to exit early
}
