package web

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/secondarykey/speaks/db"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		setTemplates(w, nil, "login.tmpl")
		return
	}

	email := r.FormValue("email")
	pswd := r.FormValue("password")

	user, err := db.SelectUser(email, pswd)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = saveLoginUser(r, w, user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := saveLoginUser(r, w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
