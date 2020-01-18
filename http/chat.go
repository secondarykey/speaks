package http

import (
	"net/http"
)

//表示用の
func chatHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	u, err := getLoginUser(r)
	if err != nil {
	}
	if u.CurrentCategory == "" {
		u.CurrentCategory = "DashBoard"
	}

	data["Category"] = u.CurrentCategory

	return "chat.tmpl", nil
}
