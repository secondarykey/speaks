package http

import (
	"net/http"
)

//表示用の
func chatHandler(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, error) {

	//今いるカテゴリにできるか？
	q := r.URL.Query()

	cat := "Dashboard"
	catS := q["cat"]
	if catS != nil || len(catS) > 0 {
		cat = catS[0]
	}

	data["Category"] = cat

	return "chat.tmpl", nil
}
