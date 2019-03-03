package web

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/secondarykey/speaks/db"
)

type messages []db.Message

func messageHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	if r.Method == "GET" {
		return "", fmt.Errorf("GET Message")
	}

	// /api/message/{category}

	url := r.URL.Path
	catS := strings.Split(url, "/")
	if len(catS) != 4 {
		return "", fmt.Errorf("Error Request URL[%s]", url)
	}

	u := data["User"].(*db.User)

	cat := catS[3]
	id := r.FormValue("lastedId")
	project := u.CurrentProject.Key

	log.Println(project)

	msgs, err := db.SelectMessages(project, cat, id)
	if err != nil {
		return "", err
	}

	data["MessageList"] = msgs
	return "", nil
}

func messageDeleteHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getLoginUser(r)
	if err != nil {
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
	err = db.DeleteMessage(msgId, user.Id)
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
