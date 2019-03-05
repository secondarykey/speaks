package web

import (
	"net/http"
)

func chatHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {
	//TODO const
	data["Category"] = "Dashboard"
	return "chat.tmpl", nil
}
