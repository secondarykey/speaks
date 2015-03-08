package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	. "speakall/config"
)

func init() {
	http.HandleFunc("/markdown", markdownHandler)
	http.HandleFunc("/me", meHandler)

	http.HandleFunc("/category", categoryHandler)

	http.HandleFunc("/message/", messageHandler)
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

func chatHandler(w http.ResponseWriter, r *http.Request) {

	session := getSession(r)
	user := session.Values["User"]
	category := "Dashboard"

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tc := make(map[string]interface{})
	tc["User"] = user
	tc["Category"] = category

	templateDir := "templates/"
	setTemplates(w, tc, templateDir+"menu.tmpl", templateDir+"chat.tmpl")
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	templateDir := "templates/"
	setTemplates(w, nil, templateDir+"menu.tmpl", templateDir+"me.tmpl")
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

func markdownHandler(w http.ResponseWriter, r *http.Request) {
	templateDir := "templates/"
	setTemplates(w, nil, templateDir+"menu.tmpl", templateDir+"markdown.tmpl")
}

func setJson(s interface{}, w http.ResponseWriter) {
	res, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
