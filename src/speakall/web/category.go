package web

import (
	"net/http"
)

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	templateDir := "templates/"
	setTemplates(w, nil, templateDir+"menu.tmpl", templateDir+"category.tmpl")
}
