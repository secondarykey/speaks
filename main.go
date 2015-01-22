package main

import (
	"./db"
	ws "./discussion"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("-secret-dao"))

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
	http.HandleFunc("/", handler)

	log.Println("############### start FileServer")
	http.Handle("/static/", http.FileServer(http.Dir("webroot")))

	http.ListenAndServe(":5555", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	user := session.Values["User"]
	tmplName := "webroot/templates/login.tmpl"
	if user != nil {
		tmplName = "webroot/templates/chat.tmpl"
	}
	tmpl := template.Must(
		template.ParseFiles(tmplName))
	tc := make(map[string]interface{})
	if err := tmpl.Execute(w, tc); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}
