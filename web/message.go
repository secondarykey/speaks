package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	if len(msgS) != 5 {
		return "", fmt.Errorf("Argument Error")
	}

	msgId := msgS[4]
	intVal, err := strconv.Atoi(msgId)
	if err != nil {
		return "", err
	}

	err = db.DeleteMessage(intVal)
	if err != nil {
		return "", err
	}

	//data["Result"] = "0"
	//data["Message"] = "No Error"

	return "", nil
}
