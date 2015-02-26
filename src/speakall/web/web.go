package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	. "speakall/config"
	"speakall/db"
	"strings"
)

func init() {
	http.HandleFunc("/markdown", markdownHandler)
	http.HandleFunc("/me", meHandler)

	http.HandleFunc("/message/", messageHandler)
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

func markdownHandler(w http.ResponseWriter, r *http.Request) {
	templateDir := "templates/"
	tmpls := make([]string, 0)
	tmpls = append(tmpls, templateDir+"markdown.tmpl")
	tmpl := template.Must(template.ParseFiles(tmpls...))
	tc := make(map[string]interface{})
	if err := tmpl.Execute(w, tc); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}

type messages []db.Message

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "GETしないで><", http.StatusBadRequest)
		return
	}
	url := r.URL.Path
	catS := strings.Split(url, "/")
	if len(catS) > 3 {
		http.Error(w, "存在しません", http.StatusNotFound)
		return
	}
	cat := catS[2]
	id := r.FormValue("lastedId")

	msgs, err := db.SelectMessage(cat, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	setJson(msgs, w)
}

func setJson(s interface{}, w http.ResponseWriter) {
	res, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
