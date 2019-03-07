package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/db"
)

func memoListHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	u := data["User"].(*db.User)
	memos, err := db.SelectProjectMemo(u.CurrentProject.Key)
	if err != nil {
		return "", err
	}
	data["MemoList"] = memos
	return "memo/list.tmpl", nil
}

func memoUpdateHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	url := r.URL.Path
	pathS := strings.Split(url, "/")

	// /api/memo/edit/{key}
	if len(pathS) != 5 {
		return "", fmt.Errorf("URL Error." + url)
	}

	key := pathS[4]

	name := r.FormValue("Name")
	content := r.FormValue("Content")

	u := data["User"].(*db.User)
	err := db.UpdateMemo(key, u.CurrentProject.Key, name, content)
	if err != nil {
		return "", err
	}
	return "", nil
}

func memoDeleteHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {
	url := r.URL.Path
	pathS := strings.Split(url, "/")
	// /api/memo/delete/{key}
	if len(pathS) != 5 {
		return "", fmt.Errorf("URL Error." + url)
	}
	key := pathS[4]
	u := data["User"].(*db.User)
	err := db.DeleteMemo(key, u.CurrentProject.Key)
	if err != nil {
		return "", err
	}
	return "", nil
}

func memoHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	url := r.URL.Path
	pathS := strings.Split(url, "/")

	if len(pathS) != 4 {
		return "", fmt.Errorf("URL Error." + url)
	}
	key := pathS[3]

	u := data["User"].(*db.User)
	memo, _ := db.SelectMemo(key, u.CurrentProject.Key)
	data["Memo"] = memo
	return "memo/edit.tmpl", nil
}

func memoViewHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	var err error
	url := r.URL.Path

	pathS := strings.Split(url, "/")
	key := pathS[2]

	u := data["User"].(*db.User)
	memo, err := db.SelectMemo(key, u.CurrentProject.Key)
	if err != nil {
		return "", err
	}
	data["Memo"] = memo

	templateDir := Config.Base.Root + "/" + Config.Web.Template
	tmplFile := templateDir + "/memo/full.tmpl"

	tmpl := template.Must(template.ParseFiles(tmplFile))

	if err := tmpl.Execute(w, data); err != nil {
		return "", err
	}
	return "", NewNoWrite("Memo Full Mode")
}
