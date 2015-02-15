package web

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"speakall/db"
)

var store *sessions.CookieStore

func init() {
	gob.Register(&db.User{})
}

func startSession(secret string) {
	store = sessions.NewCookieStore([]byte(secret))
}
