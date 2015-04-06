package web

import (
	"database/sql"
	"net/http"
	"speakall/db"
	"strings"
)

func memoListHandler(w http.ResponseWriter, r *http.Request) {
	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	memos, err := db.SelectArchiveMemo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tc := make(map[string]interface{})
	tc["MemoList"] = memos
	tc["User"] = user

	setTemplates(w, tc, "menu.tmpl", "memo/list.tmpl")
}

func memoEditHandler(w http.ResponseWriter, r *http.Request) {
	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	url := r.URL.Path
	pathS := strings.Split(url, "/")
	key := pathS[3]

	if r.Method == "GET" {
		memo, _ := db.SelectMemo(key)

		tc := make(map[string]interface{})
		tc["User"] = user
		tc["Memo"] = memo

		setTemplates(w, tc, "memo/edit.tmpl")
		return
	} else if r.Method == "DELETE" {
		err := db.DeleteMemo(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rtn := map[string]string{
			"result":  "0",
			"message": "NO ERROR",
		}
		setJson(rtn, w)
		return
	}

	name := r.FormValue("Name")
	content := r.FormValue("Content")
	db.UpdateMemo(key, name, content)

	http.Redirect(w, r, "/memo/view/"+key, http.StatusFound)
}

func memoViewHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	url := r.URL.Path
	pathS := strings.Split(url, "/")
	key := pathS[3]
	//search
	memo, err := db.SelectMemo(key)
	if err == sql.ErrNoRows {
		cat, err := db.SelectCategory(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		content := createMemoContent(key)
		name := cat.Name
		_, err = db.InsertMemo(key, name, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//id, _ := result.LastInsertId()
		//memo.Id = int(id)
		memo.Name = name
		memo.Key = key
		memo.Content = content
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tc := make(map[string]interface{})
	tc["Memo"] = memo

	setTemplates(w, tc, "memo/view.tmpl")
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
