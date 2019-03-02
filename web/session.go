package web

import (
	"encoding/gob"
	"net/http"

	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/db"

	"github.com/gorilla/sessions"
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
	v := session.Values[Config.Session.Name]
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
	session.Values[Config.Session.Name] = u
	return session.Save(r, w)
}
