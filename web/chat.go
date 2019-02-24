package web

import (
	"net/http"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {

	user := getLoginUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	category := "Dashboard"
	tc := make(map[string]interface{})
	tc["User"] = user
	tc["Category"] = category

	setTemplates(w, tc, "menu.tmpl", "chat.tmpl")
}
