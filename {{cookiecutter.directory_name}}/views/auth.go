package views

import (
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"{{cookiecutter.project_name}}/middleware"
	"{{cookiecutter.project_name}}/services"
	tmpl "{{cookiecutter.project_name}}/templates"
)

func getAuthLogin(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.Store.Get(r, "auth")
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	sessionFlashes, err := middleware.Store.Get(r, "flashes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	found, _ := middleware.GetUserFromSession(session)
	if found {
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}

	_template, err := template.ParseFS(tmpl.LoginPageFS, "*")
	if err != nil {
		log.Error(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	errorFlashes := sessionFlashes.Flashes("error")
	log.Error(errorFlashes)
	sessionFlashes.Save(r, w)

	csrfToken := csrf.Token(r)
	log.Debug(csrfToken)

	_template.ExecuteTemplate(w, "all", &struct {
		FlashedError []interface{}
		CsrfToken    string
	}{
		FlashedError: errorFlashes,
		CsrfToken:    csrfToken,
	})
}

func postAuthLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, err := middleware.Store.Get(r, "auth")
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionFlashes, err := middleware.Store.Get(r, "flashes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	success, validationIssues, user := services.AuthenticateUser(ctx, r.Form)
	if ctx.Err() != nil {
		http.Error(w, "Server Timeout", http.StatusInternalServerError)
		return
	} else if !success {
		log.Debug(validationIssues)
		sessionFlashes.AddFlash(validationIssues["error"], "error")
		sessionFlashes.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session.Values["user"] = user // DEBUG: no reference right?
	err = session.Save(r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func getAuthLogout(w http.ResponseWriter, r *http.Request) {
	// not sure why I made this POST only?
	postAuthLogout(w, r)
}

func postAuthLogout(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.Store.Get(r, "auth")
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	sessionFlashes, err := middleware.Store.Get(r, "flashes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, "user") // non-op if "user" doesn't exist
	err = session.Save(r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionFlashes.AddFlash("Successfully Logged Out", "error")
	err = sessionFlashes.Save(r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
