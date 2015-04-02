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
	gob.Register(&db.RoleMap{})
	gob.Register(&db.User{})
}

func startSession() {
	secret := Config.Session.Secret
	store = sessions.NewCookieStore([]byte(secret))
}

func getSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, Config.Session.Name)
}

func getLoginUser(r *http.Request) *db.User {
	session, err := getSession(r)
	if err != nil {
		return nil
	}
	v := session.Values["User"]
	if v == nil {
		return nil
	}
	return v.(*db.User)
}

func saveLoginUser(r *http.Request, w http.ResponseWriter, u interface{}) error {
	session, err := getSession(r)
	if err != nil {
		return err
	}
	session.Values["User"] = u
	return session.Save(r, w)
}
