package main

import (
	ws "./discussion"
	"net/http"
)

func main() {

	server := ws.NewServer()
	go server.Listen("/ws/")

	http.HandleFunc("/", handler)
	http.Handle("/static/", http.FileServer(http.Dir("webroot")))

	http.ListenAndServe(":5555", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
}
