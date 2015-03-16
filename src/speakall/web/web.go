package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	. "speakall/config"
)

func init() {

	http.HandleFunc("/memo", memoHandler)

	http.HandleFunc("/message/", messageHandler)

	http.HandleFunc("/me", meHandler)
	http.HandleFunc("/me/upload", iconRegisterHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/user/register/", userRegisterHandler)

	http.HandleFunc("/category", categoryHandler)
	http.HandleFunc("/category/list", categoryListHandler)
	http.HandleFunc("/category/view/", categoryViewHandler)
	//http.HandleFunc("/category/edit/", categoryHandler)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/store/", storeHandler)

	http.HandleFunc("/", chatHandler)
}

func Listen(static, port string) {
	startSession(Config.Session.Secret)
	http.Handle("/static/", http.FileServer(http.Dir(static)))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func setTemplates(w http.ResponseWriter, param interface{}, templateFiles ...string) {

	templateDir := "templates/"
	tmpls := make([]string, 0)
	tmpls = append(tmpls, templateDir+"layout.tmpl")

	for _, elm := range templateFiles {
		tmpls = append(tmpls, elm)
	}

	tmpl := template.Must(template.ParseFiles(tmpls...))
	if err := tmpl.Execute(w, param); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}

func setJson(s interface{}, w http.ResponseWriter) {
	res, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
