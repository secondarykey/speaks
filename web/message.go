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

func messageDeleteHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	if r.Method == "GET" {
		return "", fmt.Errorf("Get Method Error")
	}

	url := r.URL.Path
	msgS := strings.Split(url, "/")
	if len(msgS) != 4 {
		return "", fmt.Errorf("Argument Error")
	}

	msgId := msgS[3]
	user := data["User"].(*db.User)

	err := db.DeleteMessage(msgId, user.Id)
	if err != nil {
		return "", err
	}

	rtn := map[string]string{
		"result":  "0",
		"message": "NO ERROR",
	}

	data["Result"] = rtn

	return "", nil
}
