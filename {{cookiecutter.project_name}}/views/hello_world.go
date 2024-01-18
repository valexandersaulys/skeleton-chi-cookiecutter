package views

import (
	"{{cookiecutter.project_name}}/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func indexRoute(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "my-big-session")
	session.Values["foo"] = "bar"
	session.Values[42] = 45
	err := session.Save(r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Hello World!<br/><br/><a href=\"/login\">Login Here</a>"))
}
