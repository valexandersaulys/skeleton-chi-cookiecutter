package middleware

import (
	"encoding/gob"
	"{{cookiecutter.project_name}}/config"
	"{{cookiecutter.project_name}}/models"
	"fmt"
	sessions "github.com/gorilla/sessions"
	// "github.com/quasoft/memstore"
	log "github.com/sirupsen/logrus"
)

var Store sessions.Store

func InitializeSessionStore() {
	Store = sessions.NewCookieStore(
		[]byte(*config.CookieStoreSessionKeyAuthentication),
		[]byte(*config.CookieStoreSessionKeyEncryption),
	)

	// TODO: extend this out to multiple keys -- I think two
	// TODO: have this use an environmental variable

	// -------------------- Register structs here
	gob.Register(&models.User{})
	// ------------------------------------------
}
