package middleware

import (
	"encoding/gob"
	sessions "github.com/gorilla/sessions"
	"{{cookiecutter.project_name}}/config"
	"{{cookiecutter.project_name}}/models"
	// "github.com/quasoft/memstore"
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
