package web

import (
	"database/sql"
	"net/http"
	"speakall/db"
	"strings"
)

func memoListHandler(w http.ResponseWriter, r *http.Request) {
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

		templateDir := "templates/"
		setTemplates(w, tc, templateDir+"memo/edit.tmpl")
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
		content := createMemoContent(key)
		_, err = db.InsertMemo(key, "", content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//id, _ := result.LastInsertId()
		//memo.Id = int(id)
		memo.Key = key
		memo.Content = content
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tc := make(map[string]interface{})
	tc["Memo"] = memo

	templateDir := "templates/"
	setTemplates(w, tc, templateDir+"memo/view.tmpl")

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
