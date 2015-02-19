package web

import (
	"html/template"
	"net/http"
	. "speakall/config"
)

func init() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/", handler)
}

func Listen(static, port string) {
	startSession(Config.Session.Secret)
	http.Handle("/static/", http.FileServer(http.Dir(static)))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	if url == "/favicon.ico" {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	session := getSession(r)
	user := session.Values["User"]

	templateDir := "templates/"
	layoutName := templateDir + "layout.tmpl"

	tmplName := templateDir + "login.tmpl"
	category := "Dashboard"
	if user != nil {
		tmplName = templateDir + "chat.tmpl"
	}

	tmpl := template.Must(
		template.ParseFiles(layoutName, tmplName))

	tc := make(map[string]interface{})
	tc["User"] = user
	tc["Category"] = category

	//template exec
	if err := tmpl.Execute(w, tc); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}
