package main

import (
	"./db"
	ws "./discussion"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte("session-secret-dao"))

func init() {
	gob.Register(&db.User{})
}

func main() {
	var err error
	log.Println("############### start DBServer")
	err = db.Listen("db/data/SpeakAll.db")
	if err != nil {
		panic(err)
	}

	log.Println("############### start WSServer")
	server := ws.NewServer()
	go server.Listen("/ws/")

	log.Println("############### start HTTPServer")
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", handler)

	log.Println("############### start FileServer")
	http.Handle("/static/", http.FileServer(http.Dir("webroot")))

	http.ListenAndServe(":5555", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	if url == "/favicon.ico" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	session, _ := store.Get(r, "session-name")
	user := session.Values["User"]
	tmplName := "webroot/templates/login.tmpl"
	category := "Dashboard"
	if user != nil {
		tmplName = "webroot/templates/chat.tmpl"
	}

	tmpl := template.Must(
		template.ParseFiles(tmplName))

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
	session, _ := store.Get(r, "session-name")
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
	session, _ := store.Get(r, "session-name")
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
	bit, err := io.Copy(out, file)

	fmt.Println(bit)
	fmt.Println(err)

	/*
		file, err := os.Create(path)
		//file, err := os.OpenFile(path, os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fw, err := writer.CreateFormFile("uploadFile", filepath.Base(path))
		if err != nil {
			log.Println(err)
		}

		_, err = io.Copy(fw, file)
		if err != nil {
			log.Println("copy")
			log.Println(err)
		}

		err = writer.Close()
		if err != nil {
			log.Println("close")
			log.Println(err)
		}
	*/
	http.Redirect(w, r, "/", http.StatusFound)

}
