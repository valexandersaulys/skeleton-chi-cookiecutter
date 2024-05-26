package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
	"github.com/mavolin/go-htmx"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"{{cookiecutter.project_name}}/config"
	"{{cookiecutter.project_name}}/views"
)

func CreateApp() *chi.Mux {
	r := chi.NewRouter()

	// TRUE if PROD or DEV -- FALSE if TESTING or LOCAL
	var liveSite bool = *config.Environment != "LOCAL" && *config.Environment != "TESTING"
	var allowedOriginsForCors []string
	var csrfMode csrf.SameSiteMode
	if liveSite {
		csrfMode = csrf.SameSiteStrictMode
		// allowedOriginsForCors := []string{"https://foo.com"} =>  Use this to allow specific origin hosts
		panic("Missing allowedOriginsForCors in app/app.go for live sites")
	} else {
		csrfMode = csrf.SameSiteLaxMode
		allowedOriginsForCors = []string{"https://*", "http://*"}
	}

	// -------  Middleware
	if *config.Environment == "LOCAL" {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqHeaders := fmt.Sprintf("Headers: %+v", r.Header)
				log.Trace(reqHeaders)
				next.ServeHTTP(w, r)
			})
		})
	}
	if *config.Timeout != -1 {
		// check for r.Context().Err() != nil if we should timeout
		r.Use(middleware.Timeout(time.Duration(*config.Timeout) * time.Second))
	}
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5, "gzip"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOriginsForCors,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	if *config.Environment != "TESTING" {
		csrfMiddleware := csrf.Protect(
			*config.CsrfProtectionKey,
			csrf.Secure(liveSite),
			csrf.SameSite(csrfMode))
		r.Use(csrfMiddleware)
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-CSRF-Token", csrf.Token(r))
				next.ServeHTTP(w, r)
			})
		})
	}
	r.Use(htmx.NewMiddleware()) // only needed for setting response headers
	r.Use(middleware.Heartbeat("/ping"))
	if *config.Profiler {
		r.Mount("/debug", middleware.Profiler())
	}
	// --------------------

	if *config.Environment == "LOCAL" {
		r.Mount("/debug", middleware.Profiler())
	}
	views.AddRoutes(r)

	return r
}
