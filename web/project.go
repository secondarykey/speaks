package web

import (
	"net/http"

	"github.com/secondarykey/speaks/db"
)

func projectHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	if r.Method == "POST" {
		//Project Register
		name := r.FormValue("dispName")
		desc := r.FormValue("desc")

		err := db.InsertProject(name, desc)
		if err != nil {
			return "", err
		}
		return "/admin/", NewRedirect("/admin/project/")
	}

	plist, err := db.SelectProjects()
	if err != nil {
		return "", err
	}
	data["ProjectList"] = plist

	return "admin/project/list.tmpl", nil
}

func projectMemberHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	//Key値の取得

	if r.Method == "POST" {
		//Member update
	}

	return "admin/project/member.tmpl", nil
}
