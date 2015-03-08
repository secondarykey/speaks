package web

import (
	"net/http"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {

	session := getSession(r)
	user := session.Values["User"]
	category := "Dashboard"

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tc := make(map[string]interface{})
	tc["User"] = user
	tc["Category"] = category

	templateDir := "templates/"
	setTemplates(w, tc, templateDir+"menu.tmpl", templateDir+"chat.tmpl")
}
