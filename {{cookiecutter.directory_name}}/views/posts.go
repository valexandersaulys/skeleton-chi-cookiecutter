package views

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"{{cookiecutter.project_name}}/middleware"
	"{{cookiecutter.project_name}}/models"
	"{{cookiecutter.project_name}}/services"
	tmpl "{{cookiecutter.project_name}}/templates"
)

func getPostRoute(w http.ResponseWriter, r *http.Request) {
	sessionFlashes, err := middleware.Store.Get(r, "flashes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionAuth, err := middleware.Store.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var allPosts *[]models.Post
	found, user := middleware.GetUserFromSession(sessionAuth)
	if found {
		allPosts = services.RetrieveAllPostsForUser(r.Context(), *user)
	} else {
		allPosts = services.RetrieveAllPublicPosts(r.Context())
	}
	log.Debug(allPosts)

	_template, err := template.ParseFS(tmpl.ListPostTemplateFS, "*")
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	infoFlashes := sessionFlashes.Flashes("info")
	log.Debug(infoFlashes)

	sessionFlashes.Save(r, w)
	// must pass in pointer to struct here via &StructType
	_template.Execute(w, &struct {
		Posts       *[]models.Post
		FlashedInfo []interface{}
		User        *models.User
	}{
		Posts:       allPosts,
		FlashedInfo: infoFlashes,
		User:        user,
	})
}

func getDetailPostRoute(w http.ResponseWriter, r *http.Request) {
	var postUuid string = chi.URLParam(r, "postUuid")
	found, post := services.RetrieveDetailPost(r.Context(), postUuid)
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	_template, err := template.ParseFS(tmpl.DetailPostTemplateFS, "*")
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	_template.Execute(w, &struct {
		Post *models.Post
	}{Post: post})
}

func getNewPostRoute(w http.ResponseWriter, r *http.Request) {
	sessionFlashes, err := middleware.Store.Get(r, "flashes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_template, err := template.ParseFS(tmpl.NewPostTemplateFS, "*")
	if err != nil {
		log.Error(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	errorFlashes := sessionFlashes.Flashes("error")
	sessionFlashes.Save(r, w)
	log.Debug(errorFlashes)

	csrfToken := csrf.Token(r)
	log.Debug(csrfToken)

	_template.Execute(w, &struct {
		FlashedError []interface{}
		CsrfToken    string
	}{
		FlashedError: errorFlashes,
		CsrfToken:    csrfToken,
	})
}

func postNewPostRoute(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Trace(r.Body)
	log.Trace(r.Form)

	session, err := middleware.Store.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionFlashes, err := middleware.Store.Get(r, "flashes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	success, user := middleware.GetUserFromSession(session)
	if !success {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	success, validationIssues := services.CreateNewPost(r.Context(), user, r.Form)
	if !success {
		sessionFlashes.AddFlash(validationIssues["error"], "error")
		err := sessionFlashes.Save(r, w)
		if err != nil {
			log.Error(err)
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/posts/new", http.StatusSeeOther)
		return
	}

	if success {
		log.Debug("Successfully saved!")
		sessionFlashes.AddFlash("Successfully saved!", "info")
		sessionFlashes.Save(r, w)
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func getEditPostRoute(w http.ResponseWriter, r *http.Request) {
	var postUuid string = chi.URLParam(r, "postUuid")
	found, post := services.RetrieveDetailPost(r.Context(), postUuid)
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	_template, err := template.ParseFS(tmpl.EditPostTemplateFS, "*")
	if err != nil {
		log.Error(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	_template.Execute(w, &struct {
		Post models.Post
	}{
		Post: *post,
	})
}

func postEditPostRoute(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	session, err := middleware.Store.Get(r, "flashes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var postUuid string = chi.URLParam(r, "postUuid")
	success, _ := services.UpdatePostByUuid(r.Context(), r.Form, postUuid)
	if !success {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if success {
		log.Debug("Successfully Updated!")
		session.AddFlash("Successfully Updated!", "info")
		session.Save(r, w)
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
