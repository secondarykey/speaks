package web

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/secondarykey/speaks/db"
)

func memoListHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	//user, err := getLoginUser(r)
	memos, err := db.SelectArchiveMemo()
	if err != nil {
		return "", err
	}

	data["MemoList"] = memos
	return "memo/list.tmpl", nil
}

func memoEditHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	url := r.URL.Path
	pathS := strings.Split(url, "/")
	key := pathS[3]

	if r.Method == "GET" {
		memo, _ := db.SelectMemo(key)
		data["Memo"] = memo
		return "memo/edit.tmpl", nil

	} else if r.Method == "DELETE" {

		//API化
		err := db.DeleteMemo(key)
		if err != nil {
			return "", err
		}
		rtn := map[string]string{
			"result":  "0",
			"message": "NO ERROR",
		}
		setJson(rtn, w)
		return "", NewNoWrite("Memo Delete")
	}

	name := r.FormValue("Name")
	content := r.FormValue("Content")
	db.UpdateMemo(key, name, content)

	return "", NewRedirect("/memo/view/" + key)

}

func memoViewHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	var err error
	url := r.URL.Path
	pathS := strings.Split(url, "/")
	key := pathS[3]

	//search
	memo, err := db.SelectMemo(key)
	if err == sql.ErrNoRows {
		cat, err := db.SelectCategory(key)
		if err != nil {
			return "", err
		}
		content := createMemoContent(key)
		name := cat.Name
		_, err = db.InsertMemo(key, name, content)
		if err != nil {
			return "", err
		}
		//id, _ := result.LastInsertId()
		//memo.Id = int(id)
		memo.Name = name
		memo.Key = key
		memo.Content = content
	} else if err != nil {
		return "", err
	}

	data["Memo"] = memo

	return "memo/view.tmpl", nil
}

func createMemoContent(key string) string {

	msgs, err := db.SelectAllMessage(key)
	if err != nil {
		return err.Error()
	}

	content := ""
	for _, elm := range msgs {
		content += elm.UserName + ":" + elm.Created
		content += "\n"
		content += elm.Content
		content += "\n\n"
	}
	return content
}
