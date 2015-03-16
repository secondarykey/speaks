package web

import (
	"log"
	"net/http"
	"speakall/db"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		templateDir := "templates/"
		setTemplates(w, nil, templateDir+"login.tmpl")
		return
	}

	email := r.FormValue("email")
	pswd := r.FormValue("password")

	user, err := db.SelectUser(email, pswd)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = saveLoginUser(r, w, user)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	session.Values["User"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
