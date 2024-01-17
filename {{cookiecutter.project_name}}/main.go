package main

import (
	"example/skeleton/app"
	"example/skeleton/config"
	"example/skeleton/middleware"
	"example/skeleton/models"
	"fmt"
	"github.com/go-chi/docgen"
	"net/http"
)

func init() {
	config.RunInit() // **Must Run First**
	models.RunInit()

	if *config.AddDummyModels {
		models.ClearAll()
		models.CreateDummyPosts()
	}

}

func main() {
	middleware.InitializeSessionStore()

	r := app.CreateApp()
	if models.Db == nil {
		panic("Database Not initiated! Panicing and exiting")
	}

	if *config.PrintRoutes {
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "<My Project Name>",
			Intro:       "My Project Intro",
		}))
		return
	}

	http.ListenAndServe(fmt.Sprintf(":%d", *config.RuntimePort), r) // adding go causes this to exit early
}
