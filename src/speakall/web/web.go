package web

import (
	"html/template"
	"log"
	"net/http"
	. "speakall/config"
)

func init() {
	http.HandleFunc("/me", meHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/store/", storeHandler)

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
	category := "Dashboard"

	log.Println(user)

	templateDir := "templates/"
	tmpls := make([]string, 0)
	tmpls = append(tmpls, templateDir+"layout.tmpl")

	if user == nil {
		tmpls = append(tmpls, templateDir+"login.tmpl")
	} else {
		tmpls = append(tmpls, templateDir+"menu.tmpl")
		tmpls = append(tmpls, templateDir+"chat.tmpl")
	}

	tmpl := template.Must(template.ParseFiles(tmpls...))

	tc := make(map[string]interface{})
	tc["User"] = user
	tc["Category"] = category

	//template exec
	if err := tmpl.Execute(w, tc); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}

func meHandler(w http.ResponseWriter, r *http.Request) {

	templateDir := "templates/"
	tmpls := make([]string, 0)
	tmpls = append(tmpls, templateDir+"layout.tmpl")
	tmpls = append(tmpls, templateDir+"menu.tmpl")
	tmpls = append(tmpls, templateDir+"me.tmpl")
	tmpl := template.Must(template.ParseFiles(tmpls...))
	tc := make(map[string]interface{})

	//tc["User"] = user
	//tc["Category"] = category

	if err := tmpl.Execute(w, tc); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}
