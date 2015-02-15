package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	. "speakall/config"
	"speakall/db"
)

func init() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/", handler)
}

func Listen(static, port string) {
	startSession(Config.Session.Secret)
	http.Handle("/static/", http.FileServer(http.Dir(static)))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	if url == "/favicon.ico" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	session, _ := store.Get(r, Config.Session.Name)
	user := session.Values["User"]

	templateDir := "templates/"
	layoutName := templateDir + "layout.tmpl"

	tmplName := templateDir + "login.tmpl"
	category := "Dashboard"
	if user != nil {
		tmplName = templateDir + "chat.tmpl"
	}

	tmpl := template.Must(
		template.ParseFiles(layoutName, tmplName))

	tc := make(map[string]interface{})
	tc["User"] = user
	tc["Category"] = category

	if err := tmpl.Execute(w, tc); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pswd := r.FormValue("password")

	user, err := db.SelectUser(email, pswd)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	session, _ := store.Get(r, Config.Session.Name)
	log.Println(session.ID)
	log.Println(session.IsNew)
	session.Values["User"] = user
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, Config.Session.Name)
	session.Values["User"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	path, _ := os.Getwd()
	path += "/test.txt"
	log.Println(path)

	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	out, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
