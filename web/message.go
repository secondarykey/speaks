package web

import (
	"net/http"
	"strings"

	"github.com/secondarykey/speaks/db"
)

type messages []db.Message

func messageHandler(w http.ResponseWriter, r *http.Request) {
	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

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

func messageDeleteHandler(w http.ResponseWriter, r *http.Request) {
	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == "GET" {
		http.Error(w, "GETしないで><", http.StatusBadRequest)
		return
	}

	url := r.URL.Path
	msgS := strings.Split(url, "/")
	if len(msgS) > 4 {
		http.Error(w, "存在しません", http.StatusNotFound)
		return
	}
	msgId := msgS[3]
	err := db.DeleteMessage(msgId, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rtn := map[string]string{
		"result":  "0",
		"message": "NO ERROR",
	}

	setJson(rtn, w)
}
