package middleware

import (
	// "github.com/go-chi/chi/v5/middleware"
	// "{{cookiecutter.project_name}}/config"
	"fmt"
	sessions "github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"net/http"
	"{{cookiecutter.project_name}}/models"
)

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "auth")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		exists, _ := GetUserFromSession(session)
		if !exists {
			http.Error(w, fmt.Sprint(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserFromSession(session *sessions.Session) (bool, *models.User) {
	var user *models.User

	tmp := session.Values["user"]
	if tmp == nil {
		return false, user
	}
	user = tmp.(*models.User)
	log.Debug(fmt.Sprintf("Found User: %+v", user))

	return true, user
}
