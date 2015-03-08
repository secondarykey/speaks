package web

import (
	"net/http"
)

func meHandler(w http.ResponseWriter, r *http.Request) {
	templateDir := "templates/"
	setTemplates(w, nil, templateDir+"menu.tmpl", templateDir+"me.tmpl")
}
