package web

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"net/http"
	. "speakall/config"
	"speakall/db"
)

var store *sessions.CookieStore

func init() {
	gob.Register(&db.User{})
}

func startSession(secret string) {
	store = sessions.NewCookieStore([]byte(secret))
}

func getSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, Config.Session.Name)
	return session
}

func getLoginUser(r *http.Request) interface{} {
	session := getSession(r)
	return session.Values["User"]
}

func saveLoginUser(r *http.Request, w http.ResponseWriter, u interface{}) error {
	session := getSession(r)
	session.Values["User"] = u
	return session.Save(r, w)
}
