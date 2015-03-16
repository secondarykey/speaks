package web

import (
	"net/http"
	"speakall/db"
	"strings"
)

type messages []db.Message

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "GETしないで><", http.StatusBadRequest)
		return
	}
	url := r.URL.Path
	catS := strings.Split(url, "/")
	if len(catS) > 3 {
		http.Error(w, "存在しません", http.StatusNotFound)
		return
	}
	cat := catS[2]
	id := r.FormValue("lastedId")

	msgs, err := db.SelectMessage(cat, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	setJson(msgs, w)
}
