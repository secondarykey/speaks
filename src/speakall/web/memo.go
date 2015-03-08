package web

import (
	"net/http"
)

func markdownHandler(w http.ResponseWriter, r *http.Request) {
	templateDir := "templates/"
	setTemplates(w, nil, templateDir+"menu.tmpl", templateDir+"markdown.tmpl")
}
