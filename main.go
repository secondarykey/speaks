package main

import (
	"./db"
	ws "./discussion"
	"html/template"
	"log"
	"net/http"
)

func main() {

	log.Println("############### start DBServer")
	db.Listen("db/data/SpeakAll.db")

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
	tmpl := template.Must(
		template.ParseFiles("webroot/templates/chat.tmpl"))
	tc := make(map[string]interface{})
	if err := tmpl.Execute(w, tc); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}
