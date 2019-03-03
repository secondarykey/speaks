package web

import (
	"net/http"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {

	user, err := getLoginUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	category := "Dashboard"
	tc := make(map[string]interface{})

	tc["User"] = user
	tc["Category"] = category

	setTemplates(w, tc, "chat.tmpl")
}
