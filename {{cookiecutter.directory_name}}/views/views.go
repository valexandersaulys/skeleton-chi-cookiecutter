package views

import (
	"embed"
	"{{cookiecutter.project_name}}/middleware"
	tmpl "{{cookiecutter.project_name}}/templates"
	// "fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

func AddRoutes(r *chi.Mux) *chi.Mux {

	// --- Your routes here
	r.Get("/", indexRoute)
	r.Get("/posts", getPostRoute)
	r.Get("/posts/{postUuid:^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$}",
		getDetailPostRoute)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthRequired)
		r.Get("/posts/new", getNewPostRoute)
		r.Post("/posts/new", postNewPostRoute)
	})
	r.Get("/posts/{postUuid:^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$}/edit",
		getEditPostRoute)
	r.Post("/posts/{postUuid:^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$}/edit",
		postEditPostRoute)
	r.Get("/login", getAuthLogin)
	r.Post("/login", postAuthLogin)
	r.Get("/logout", getAuthLogout)
	// r.Post("/logout", postAuthLogout)
	// --------------------

	// ---------- Static, Missing, MethodNotAllowed
	// r.Get("/static", rootStaticPath(http.FS(StaticFilesFS)))
	HandleStaticFiles(r, "/public")
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_template, err := template.ParseFS(tmpl.MissingFS, "*")
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			log.Error(err)
			return
		}
		_template.Execute(w, &struct{}{})
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_template, err := template.ParseFS(tmpl.MethodNotAllowedFS, "*")
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			log.Error(err)
			return
		}
		_template.Execute(w, &struct{}{})
	})
	// --------------------

	return r
}

// Must be a sub-directory
//
//go:embed public
var staticFilesFS embed.FS

func HandleStaticFiles(r chi.Router, path string) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	subFS, _ := fs.Sub(staticFilesFS, "public") // strip out the public/ from path
	fileServer := http.FileServer(http.FS(subFS))

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, fileServer)
		fs.ServeHTTP(w, r)
	})
}
