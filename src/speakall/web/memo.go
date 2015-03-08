package web

import (
	"net/http"
)

func memoHandler(w http.ResponseWriter, r *http.Request) {
	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tc := make(map[string]interface{})
	tc["User"] = user

	templateDir := "templates/"
	setTemplates(w, tc, templateDir+"menu.tmpl", templateDir+"memo.tmpl")
}
